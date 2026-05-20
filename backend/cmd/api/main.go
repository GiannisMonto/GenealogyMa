package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/internal/infrastructure/cache"
	"github.com/genealogy-ma/platform/internal/infrastructure/config"
	"github.com/genealogy-ma/platform/internal/infrastructure/persistence"
	"github.com/genealogy-ma/platform/internal/interfaces/http/controller"
	"github.com/genealogy-ma/platform/internal/interfaces/http/middleware"
	"github.com/genealogy-ma/platform/pkg/auth"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 族谱数字化管理平台 API
// @version 1.0
// @description 族谱数字化管理平台后端API接口文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db, err := persistence.NewDatabase(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	// 执行数据库迁移
	migrator := persistence.NewMigrator(db.DB)
	if err := migrator.Run(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// 初始化Redis
	redisCache, err := cache.NewRedisCache(&cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect Redis: %v", err)
	} else {
		defer redisCache.Close()
		log.Println("Redis connected successfully")
	}

	// 初始化JWT服务
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.ExpireHours,
		cfg.JWT.RefreshExpireHours,
	)
	// 设置全局JWT服务实例
	auth.SetGlobalJWTService(jwtService)

	// 创建Gin引擎
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ===== 初始化仓储和服务 =====

	// 令牌黑名单服务
	var tokenBlacklist *cache.TokenBlacklist
	if redisCache != nil {
		tokenBlacklist = cache.NewTokenBlacklist(redisCache.Client())
	}

	// 人物领域
	personRepo := persistence.NewPersonRepository(db.DB)
	personService := service.NewPersonService(personRepo)
	personController := controller.NewPersonController(personService)

	// 用户领域
	userRepo := persistence.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo, jwtService)
	authController := controller.NewAuthController(userService, tokenBlacklist)
	rbacController := controller.NewRBACController(userService)

	// 墓园领域
	cemeteryRepo := persistence.NewCemeteryRepository(db.DB)
	cemeteryGraveRepo := persistence.NewGraveRepository(db.DB)
	cemeteryService := service.NewCemeteryService(cemeteryRepo, cemeteryGraveRepo)
	cemeteryController := controller.NewCemeteryController(cemeteryService)

	// 文化领域
	documentRepo := persistence.NewDocumentRepository(db.DB)
	storyRepo := persistence.NewStoryRepository(db.DB)
	teachingsRepo := persistence.NewFamilyTeachingsRepository(db.DB)
	cultureService := service.NewCultureService(documentRepo, storyRepo, teachingsRepo)
	cultureController := controller.NewCultureController(cultureService)

	// 社区领域
	userProfileRepo := persistence.NewUserProfileRepository(db.DB)
	postRepo := persistence.NewPostRepository(db.DB)
	commentRepo := persistence.NewCommentRepository(db.DB)
	messageRepo := persistence.NewMessageRepository(db.DB)
	communityService := service.NewCommunityService(userProfileRepo, postRepo, commentRepo, messageRepo)
	communityController := controller.NewCommunityController(communityService)

	// ===== 路由设置 =====
	authMiddleware := middleware.JWTAuth(jwtService)

	// API路由组
	apiV1 := r.Group("/api/v1")
	{
		// 健康检查
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"time":   time.Now().Format(time.RFC3339),
			})
		})

		// 注册认证控制器路由（包含登录、注册、用户管理等）
		authController.RegisterRoutes(apiV1, authMiddleware)

		// 注册人物控制器路由
		personController.RegisterRoutes(apiV1, authMiddleware)

		// 注册角色权限管理控制器路由
		rbacController.RegisterRoutes(apiV1, authMiddleware)

		// 注册墓园控制器路由
		cemeteryController.RegisterRoutes(apiV1, authMiddleware)

		// 注册文化控制器路由
		cultureController.RegisterRoutes(apiV1, authMiddleware)

		// 注册社区控制器路由
		communityController.RegisterRoutes(apiV1, authMiddleware)
	}

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 5秒超时关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
