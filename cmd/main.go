package main

import (
	"os"
	config "stream-session-api/internal/conf"
	"stream-session-api/internal/service/worker"
	"stream-session-api/pkg"

	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file (once)
	err := godotenv.Load(".env")
	if err != nil {
		pkg.LogFatal("failed to read .env")
		os.Exit(1)
	}

	pkg.LogInfo(os.Getenv("APPLICATION_NAME") + " " + os.Getenv("APPLICATION_VERSION") + " is running... ")

	// Get config
	if err := config.Get(); err != nil {
		pkg.LogFatal("get config fail!")
		os.Exit(1)
	}
}

func main() {

	// Init gRPC server
	if err := worker.InitGrpcServer(); err != nil {
		pkg.LogFatal("init gRPC server fail!")
		os.Exit(2)
	}

	go worker.GrpcServer()
	go worker.PeriodicStreamSessionCheck()
	select {}
}
