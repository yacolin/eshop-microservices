package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"eshop-microservices/internal/inventory-service/api/grpc"
	"eshop-microservices/internal/inventory-service/api/handlers"
	"eshop-microservices/internal/inventory-service/api/routes"
	"eshop-microservices/internal/inventory-service/domain/models"
	"eshop-microservices/internal/inventory-service/domain/repositories"
	inventorymq "eshop-microservices/internal/inventory-service/mq"
	"eshop-microservices/internal/inventory-service/service"
	"eshop-microservices/pkg/config"
	"eshop-microservices/pkg/database"
	"eshop-microservices/pkg/mq"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// App inventory-service 应用入口，负责配置加载、依赖装配与服务启动
type App struct {
	cfg        *config.Config
	db         *gorm.DB
	mqClient   *mq.Client
	engine     *gin.Engine
	grpcServer *grpc.Server
	mqConsumer *inventorymq.Consumer
}

// New 加载配置并创建 App，configPath 为空时从环境变量 CONFIG_PATH 读取，再默认 configs/inventory-service.yaml
func New(configPath string) (*App, error) {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "configs/inventory-service.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &App{cfg: cfg}, nil
}

// Run 装配依赖、启动 HTTP 和 gRPC 服务并阻塞直到收到退出信号
func (a *App) Run() error {
	if err := a.wire(); err != nil {
		return err
	}
	if a.mqClient != nil {
		defer a.mqClient.Close()
	}

	// 启动 HTTP 服务
	addr := fmt.Sprintf(":%d", a.cfg.Server.Port)
	go func() {
		log.Printf("inventory-service HTTP listen %s", addr)
		if err := a.engine.Run(addr); err != nil {
			log.Printf("http server: %v", err)
		}
	}()

	// 启动 gRPC 服务
	if a.grpcServer != nil {
		go func() {
			log.Printf("inventory-service gRPC listen :%s", "50051")
			if err := a.grpcServer.Start(); err != nil {
				log.Printf("grpc server: %v", err)
			}
		}()
	}

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

	// 优雅关闭 MQ 消费者
	if a.mqConsumer != nil {
		a.mqConsumer.Stop()
	}

	// 优雅关闭 gRPC 服务
	if a.grpcServer != nil {
		a.grpcServer.Stop()
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
	// 自动迁移数据库表
	if err := a.db.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.Inventory{},
		&models.Comment{},
	); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	rdb, err := database.NewRedis(a.cfg.Redis.Addr(), a.cfg.Redis.Password, a.cfg.Redis.DB)
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	_ = rdb // 预留：库存缓存等

	if a.cfg.RabbitMQ.URL != "" {
		a.mqClient, err = mq.NewClient(a.cfg.RabbitMQ.URL, a.cfg.RabbitMQ.Exchange)
		if err != nil {
			log.Printf("rabbitmq (optional): %v", err)
		}
	}

	inventoryRepo := repositories.NewInventoryRepository(a.db)
	inventorySvc := service.NewInventoryService(inventoryRepo)
	var pub *inventorymq.Publisher
	if a.mqClient != nil {
		pub = inventorymq.NewPublisher(a.mqClient)
	}
	productHandler := handlers.NewProductHandler(inventorySvc, pub)
	inventoryHandler := handlers.NewInventoryHandler(inventorySvc, pub)
	categoryHandler := handlers.NewCategoryHandler(inventorySvc, pub)

	a.engine = gin.New()
	routes.Setup(a.engine, productHandler, inventoryHandler, categoryHandler)

	// 创建 gRPC 服务器
	a.grpcServer = grpc.NewServer(inventoryRepo, "50051")

	// 创建 MQ 消费者
	if a.mqClient != nil {
		conn := a.mqClient.GetConnection()
		if conn != nil {
			a.mqConsumer, err = inventorymq.NewConsumer(conn, inventorySvc, a.cfg.RabbitMQ.Exchange)
			if err != nil {
				log.Printf("mq consumer init (optional): %v", err)
			}
		}
	}

	return nil
}
