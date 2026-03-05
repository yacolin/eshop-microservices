package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"eshop-microservices/internal/user-service/api/handlers"
	"eshop-microservices/internal/user-service/api/routes"
	"eshop-microservices/internal/user-service/domain/models"
	"eshop-microservices/internal/user-service/domain/repositories"
	usermq "eshop-microservices/internal/user-service/mq"
	"eshop-microservices/internal/user-service/service"
	"eshop-microservices/pkg/config"
	"eshop-microservices/pkg/database"
	"eshop-microservices/pkg/mq"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	cfg      *config.Config
	db       *gorm.DB
	mqClient *mq.Client
	engine   *gin.Engine
}

func New(configPath string) (*App, error) {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "configs/user-service.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &App{cfg: cfg}, nil
}

func (a *App) Run() error {
	if err := a.wire(); err != nil {
		return err
	}
	if a.mqClient != nil {
		defer a.mqClient.Close()
	}

	addr := fmt.Sprintf(":%d", a.cfg.Server.Port)
	go func() {
		log.Printf("user-service listen %s", addr)
		if err := a.engine.Run(addr); err != nil {
			log.Printf("http server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
	return nil
}

func (a *App) wire() error {
	gin.SetMode(a.cfg.Server.Mode)

	db, err := database.NewMySQL(a.cfg.MySQL.DSN(), logger.Info)
	if err != nil {
		return fmt.Errorf("mysql: %w", err)
	}
	a.db = db
	if err := a.db.AutoMigrate(
		&models.User{},
		&models.UserInfo{},
		&models.UserIdentity{},
		&models.AuthToken{},
		&models.LoginHistory{},
		&models.Permission{},
		&models.RolePermission{},
	); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	rdb, err := database.NewRedis(a.cfg.Redis.Addr(), a.cfg.Redis.Password, a.cfg.Redis.DB)
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	_ = rdb

	if a.cfg.RabbitMQ.URL != "" {
		a.mqClient, err = mq.NewClient(a.cfg.RabbitMQ.URL, a.cfg.RabbitMQ.Exchange)
		if err != nil {
			log.Printf("rabbitmq (optional): %v", err)
		}
	}

	// 初始化 repositories
	userRepo := repositories.NewUserRepository(a.db)
	identityRepo := repositories.NewUserIdentityRepository(a.db)
	tokenRepo := repositories.NewAuthTokenRepository(a.db)
	loginHistoryRepo := repositories.NewLoginHistoryRepository(a.db)
	permissionRepo := repositories.NewPermissionRepository(a.db)

	// 初始化 token service
	tokenSvc := service.NewTokenService(a.cfg.JWT.Secret, tokenRepo)

	// 初始化 auth service
	authSvc := service.NewAuthService(a.db, userRepo, identityRepo, tokenRepo, loginHistoryRepo, tokenSvc)

	// 初始化 user service
	userSvc := service.NewUserService(userRepo)
	userSvc.SetJWTSecret(a.cfg.JWT.Secret)

	// 初始化 permission service
	permissionSvc := service.NewPermissionService(permissionRepo, userRepo)

	var pub *usermq.Publisher
	if a.mqClient != nil {
		pub = usermq.NewPublisher(a.mqClient)
	}

	// 初始化 handlers
	userHandler := handlers.NewUserHandler(userSvc, pub)
	authHandler := handlers.NewAuthHandler(authSvc, tokenSvc, userSvc)
	permissionHandler := handlers.NewPermissionHandler(permissionSvc, userSvc)

	a.engine = gin.New()
	routes.Setup(a.engine, userHandler, authHandler, permissionHandler)
	return nil
}
