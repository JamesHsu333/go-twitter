package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/server"
	"github.com/JamesHsu333/go-twitter/pkg/database/postgres"
	"github.com/JamesHsu333/go-twitter/pkg/database/redis"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/JamesHsu333/go-twitter/pkg/version"
	"github.com/common-nighthawk/go-figure"
)

const (
	banner = `
____________________________________O/_______
                                    O\
`
)

func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	figure.NewColorFigure(cfg.Server.Name, "pepper", "cyan", true).Print()
	fmt.Printf("%s%s", cfg.Server.Description, banner)

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Info(version.PrintVersion())
	appLogger.Infof("LogLevel: %s, Mode: %s, SSL: %v", cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	traceProvider, err := tracer.NewJaeger(cfg)
	if err != nil {
		appLogger.Fatal("Cannot create tracer", err)
	} else {
		appLogger.Info("Jaeger connected")
	}

	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			appLogger.Error("Cannot shutdown tracer", err)
		}
	}()
	s := server.NewServer(cfg, psqlDB, redisClient, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
