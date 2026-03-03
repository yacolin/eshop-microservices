package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"eshop-microservices/internal/order-service/api/handlers"
	"eshop-microservices/internal/order-service/api/routes"
	"eshop-microservices/internal/order-service/clients"
	"eshop-microservices/internal/order-service/domain/models"
	"eshop-microservices/internal/order-service/domain/repositories"
	ordermq "eshop-microservices/internal/order-service/mq"
	"eshop-microservices/internal/order-service/service"
	"eshop-microservices/pkg/config"
	"eshop-microservices/pkg/database"
	"eshop-microservices/pkg/mq"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// App order-service 应用入口，负责配置加载、依赖装配与服务启动
type App struct {
	cfg             *config.Config
	db              *gorm.DB
	mqClient        *mq.Client
	engine          *gin.Engine
	inventoryClient *clients.InventoryClient
	mqConsumer      *ordermq.Consumer
}

// New 加载配置并创建 App，configPath 为空时从环境变量 CONFIG_PATH 读取，再默认 configs/order-service.yaml
func New(configPath string) (*App, error) {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "configs/order-service.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &App{cfg: cfg}, nil
}

// Run 装配依赖、启动 HTTP 服务并阻塞直到收到退出信号
func (a *App) Run() error {
	if err := a.wire(); err != nil {
		return err
	}

	// 启动 HTTP 服务
	addr := fmt.Sprintf(":%d", a.cfg.Server.Port)
	go func() {
		log.Printf("order-service HTTP listen %s", addr)
		if err := a.engine.Run(addr); err != nil {
			log.Printf("http server: %v", err)
		}
	}()

	// 启动 MQ 消费者
	if a.mqConsumer != nil {
		if err := a.mqConsumer.Start(); err != nil {
			log.Printf("mq consumer: %v", err)
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	// 优雅关闭
	if a.mqConsumer != nil {
		a.mqConsumer.Stop()
	}
	if a.mqClient != nil {
		a.mqClient.Close()
	}
	if a.inventoryClient != nil {
		a.inventoryClient.Close()
	}

	return nil
}

// wire 初始化 DB、Redis、MQ，注册路由
func (a *App) wire() error {
	gin.SetMode(a.cfg.Server.Mode)

	db, err := database.NewMySQL(a.cfg.MySQL.DSN(), logger.Info)
	if err != nil {
		return fmt.Errorf("mysql: %w", err)
	}
	a.db = db
	if err := a.db.AutoMigrate(&models.Order{}, &models.OrderItem{}); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	rdb, err := database.NewRedis(a.cfg.Redis.Addr(), a.cfg.Redis.Password, a.cfg.Redis.DB)
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	_ = rdb // 预留：订单缓存等

	// 初始化 MQ 客户端
	if a.cfg.RabbitMQ.URL != "" {
		a.mqClient, err = mq.NewClient(a.cfg.RabbitMQ.URL, a.cfg.RabbitMQ.Exchange)
		if err != nil {
			log.Printf("rabbitmq (optional): %v", err)
		}
	}

	// 初始化库存服务 gRPC 客户端
	// 从环境变量或配置获取库存服务地址，默认 localhost:50051
	inventoryAddr := os.Getenv("INVENTORY_SERVICE_ADDR")
	if inventoryAddr == "" {
		inventoryAddr = "localhost:50051"
	}

	a.inventoryClient, err = clients.NewInventoryClient(inventoryAddr)
	if err != nil {
		log.Printf("inventory client (optional): %v", err)
		a.inventoryClient = nil // 如果连接失败，继续运行但不使用库存服务
	}

	orderRepo := repositories.NewOrderRepository(a.db)
	orderSvc := service.NewOrderService(orderRepo, a.inventoryClient)
	var pub *ordermq.Publisher
	if a.mqClient != nil {
		pub = ordermq.NewPublisher(a.mqClient)
	}
	orderHandler := handlers.NewOrderHandler(orderSvc, pub)

	// 创建 MQ 消费者
	if a.mqClient != nil {
		conn := a.mqClient.GetConnection()
		if conn != nil {
			a.mqConsumer, err = ordermq.NewConsumer(conn, orderSvc, a.cfg.RabbitMQ.Exchange)
			if err != nil {
				log.Printf("mq consumer init (optional): %v", err)
			}
		}
	}

	a.engine = gin.New()
	routes.Setup(a.engine, orderHandler)
	return nil
}
