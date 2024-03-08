package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	appBackendVersion = "0.1.0"
)

type syncController struct {
	stopChan  chan struct{}
	fatalChan chan error
	wg        sync.WaitGroup
}

func main() {
	var args struct {
		ConfigPath   *string `arg:"positional" default:"config.toml" help:"configuration file path"`
		Debug        bool    `arg:"-d,--debug" help:"enable debug logging"`
		PrintVersion bool    `arg:"-V,--version" help:"print program version"`
	}
	arg.MustParse(&args)

	if args.PrintVersion {
		fmt.Println("App backend version", appBackendVersion)
		os.Exit(0)
	}

	stopChan := make(chan struct{})
	go exitHandler(stopChan)

	log.SetTimeFormat(time.DateTime)

	if args.Debug {
		if args.Debug {
			log.SetLevel(log.DebugLevel)
		}
	}

	cfg, err := getConfig(*args.ConfigPath)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	sc := syncController{
		stopChan:  stopChan,
		fatalChan: make(chan error),
		wg:        sync.WaitGroup{},
	}

	dbPool, err := startDbPool(cfg)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}

	sc.wg.Add(1)
	go serve(cfg, &sc, dbPool)

	select {
	case <-sc.stopChan:
		dbPool.Close()
		sc.wg.Wait()
		break
	case err := <-sc.fatalChan:
		log.Errorf("Exiting due to fatal runtime error: %v", err)
		close(sc.stopChan)
		dbPool.Close()
		sc.wg.Wait()
		log.Fatalf("Runtime: %v", err)
	}

	os.Exit(0)
}

func exitHandler(stopChan chan struct{}) {
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, os.Interrupt, syscall.SIGTERM)

	log.Debug("received stop signal", "signal", <-stopSig)
	close(stopChan)

	// Force exit on second signal from stopSig
	log.Fatal("force exit:", "signal", <-stopSig)
}
