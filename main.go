package main

import (
	"context"
	"flag"
	"fmt"
	"gen/log"
	"gen/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		configFile = flag.String("config", "app.ini", "path to config file")
	)
	flag.Parse()

	defer func() {
		if err := log.Logger.Sync(); err != nil {
			fmt.Printf("Failed to close log: %s\n", err)
		}
	}()

	s, err := server.New(server.Config{
		ConfigFile: *configFile,
	})

	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
		os.Exit(-1)
	}

	ctx := context.Background()
	go listenToSystemSignals(ctx, s)

	if err := s.Run(); err != nil {
		fmt.Printf("Failed to run server: %s\n", err)
		os.Exit(-1)
	}
}

func listenToSystemSignals(ctx context.Context, s *server.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case sig := <-signalChan:
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			if err := s.Shutdown(ctx, fmt.Sprintf("System signal: %s", sig)); err != nil {
				fmt.Printf("Timed out waiting for server to shutdown\n")
			}
			cancel()
			return
		}
	}
}
