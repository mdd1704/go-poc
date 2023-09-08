package main

import (
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	joonix "github.com/joonix/log"
	"github.com/opentracing/opentracing-go"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"

	"go-poc/external"
	"go-poc/service"
	inventoryHandler "go-poc/service/inventory/handler"
	inventoryAdapter "go-poc/service/inventory/repository/adapter"
	inventoryPort "go-poc/service/inventory/repository/port"
	inventoryUsecase "go-poc/service/inventory/usecase"
	salesChannelHandler "go-poc/service/saleschannel/handler"
	salesChannelAdapter "go-poc/service/saleschannel/repository/adapter"
	salesChannelPort "go-poc/service/saleschannel/repository/port"
	salesChannelUsecase "go-poc/service/saleschannel/usecase"
	"go-poc/utils"
	"go-poc/utils/activity"
	"go-poc/utils/log"
)

const (
	salesChannelService = "saleschannel"
	inventoryService    = "inventory"
)

func main() {
	godotenv.Load(".env")
	configureLogging()
	ctx := activity.NewContext("init_app")

	salesChannelJaeger, closer, err := external.NewJaeger(salesChannelService)
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "jaeger open telemetry error"))
		panic(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(salesChannelJaeger)

	redisDB, err := external.NewRedis()
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "redis connection error"))
		panic(err)
	}

	memcacheDB, err := external.NewMemcache()
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "memcache connection error"))
		panic(err)
	}

	// Register sales channel service
	var salesChannelDB *sql.DB
	var salesChannelMain salesChannelPort.MainRepository
	switch os.Getenv("SALES_CHANNEL_MAIN") {
	case "mysql":
		salesChannelDB, err = external.NewMySQL(salesChannelService)
		if err != nil {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "mysql connection error"))
			panic(err)
		}

		salesChannelMain = salesChannelAdapter.NewMySQL(salesChannelDB)
	case "postgres":
		salesChannelDB, err = external.NewPostgres(salesChannelService)
		if err != nil {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "mysql connection error"))
			panic(err)
		}

		salesChannelMain = salesChannelAdapter.NewPostgres(salesChannelDB)
	}

	var salesChannelCache salesChannelPort.CacheRepository
	switch os.Getenv("SALES_CHANNEL_CACHE") {
	case "redis":
		salesChannelCache = salesChannelAdapter.NewRedis(redisDB)
	case "memcache":
		salesChannelCache = salesChannelAdapter.NewMemcache(memcacheDB)
	}

	salesChannelUsecase := salesChannelUsecase.NewChannel(salesChannelMain, salesChannelCache)
	salesChannelHandler := salesChannelHandler.NewChannel(salesChannelUsecase)

	// Register inventory service
	var inventoryDB *sql.DB
	var inventoryMain inventoryPort.MainRepository
	switch os.Getenv("INVENTORY_MAIN") {
	case "mysql":
		inventoryDB, err = external.NewMySQL(inventoryService)
		if err != nil {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "mysql connection error"))
			panic(err)
		}

		inventoryMain = inventoryAdapter.NewMySQL(inventoryDB)
	case "postgres":
		inventoryDB, err = external.NewPostgres(inventoryService)
		if err != nil {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "postgres connection error"))
			panic(err)
		}

		inventoryMain = inventoryAdapter.NewPostgres(inventoryDB)
	}

	var inventoryCache inventoryPort.CacheRepository
	switch os.Getenv("INVENTORY_CACHE") {
	case "redis":
		inventoryCache = inventoryAdapter.NewRedis(redisDB)
	case "memcache":
		inventoryCache = inventoryAdapter.NewMemcache(memcacheDB)
	}

	inventoryUsecase := inventoryUsecase.NewLocation(inventoryMain, inventoryCache)
	inventoryHandler := inventoryHandler.NewLocation(inventoryUsecase)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		default:
			return
		}
	} else {
		// Set application mode
		mode := os.Getenv("APP_MODE")
		gin.SetMode(mode)

		corsConfig := cors.New(cors.Config{
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"*"},
			AllowOrigins:     []string{"*"},
			AllowCredentials: true,
		})

		// Define application
		app := gin.Default()
		app.Use(
			corsConfig,
			gin.Recovery(),
			gin.Logger(),
		)

		// Init route
		service.InitRoute(
			ctx,
			app,
			salesChannelHandler,
			inventoryHandler,
		)

		// Start HTTP server
		srv := &http.Server{
			Addr:    ":" + os.Getenv("SERVER_PORT"),
			Handler: app,
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "running error"))
		}

		if err := srv.Shutdown(ctx); err != nil {
			log.WithContext(ctx).Error("cannot shutdown http server")
		}
	}

	// OS notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.WithContext(ctx).Info("service stopped")
}

func configureLogging() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.AddHook(utils.LogrusSourceContextHook{})

	if gin.Mode() != "release" {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	} else {
		logrus.SetFormatter(&joonix.FluentdFormatter{})
	}
}
