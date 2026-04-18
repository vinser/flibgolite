package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kardianos/service"
	"github.com/vinser/flibgolite/internal/app"
	"github.com/vinser/flibgolite/internal/core/model"
	"github.com/vinser/flibgolite/internal/index"
)

var version, buildTime, target, goversion string
var rootDir string

func main() {
	serviceFlag := flag.String("service", "", `control FLibGoLite system service`)
	reindexFlag := flag.Bool("reindex", false, `empty book stock database and then scan book stock directory to add books to database`)
	configFlag := flag.Bool("config", false, `create default config file in ./config folder for customization and exit`)
	helpFlag := flag.Bool("help", false, `display extended command help and exit`)
	versionFlag := flag.Bool("version", false, `output version information and exit`)
	flag.Parse()
	if flag.Arg(0) != "" {
		rootDir = flag.Arg(0)
	} else if os.Getenv("FLIBGOLITE_ROOT") != "" {
		rootDir = os.Getenv("FLIBGOLITE_ROOT")
	} else {
		x, _ := os.Executable()
		rootDir = filepath.Dir(x)
	}

	switch {
	case flag.NFlag() > 1:
		fmt.Println(`Error: More than one OPTION used`)
		displayHelp()
		os.Exit(1)
	case *helpFlag:
		displayHelp()
	case *versionFlag:
		displayVersion()
	case *configFlag:
		defaultConfig()
	case *reindexFlag:
		reindexStock()
	case *serviceFlag != "":
		controlService(*serviceFlag)
	default:
		runService(initService())
	}
}

func displayHelp() {
	exeName, _ := os.Executable()
	fmt.Printf(
		`	
FLibGoLite is multiplatform lightweight OPDS server with SQLite database book search index
This program was built for %s-%s

Usage: %s [OPTION] [data directory]

With no OPTION program will run in console mode (Ctrl+C to exit)
Caution: Only one OPTION can be used at a time

OPTION should be one of:
  -service [action]     control FLibGoLite system service
	  where action is one of: install, start, stop, restart, uninstall, status 
  -reindex              empty book stock index and then scan book stock folder to add books to index (database)
  -config               create default config file in ./config folder for customization
  -help                 display this help
  -version              output version information

data directory is optional (current directory by default)
  
Examples:
  ./flibgolite                      Run FLibGoLite in console mode with app data in current directory
  ./flibgolite -service install     Install FLibGoLite as a system service

Documentation: <https://vinser.github.io/flibgolite-docs>
Sources: <https://github.com/vinser/flibgolite>
`,
		runtime.GOOS, runtime.GOARCH, filepath.Base(exeName))
	os.Exit(0)
}

func displayVersion() {
	fmt.Printf("FLibGoLite OPDS server\n")
	fmt.Printf("Version: %s (%s)\n", version, target)
	fmt.Printf("Build time: %s\n", buildTime)
	fmt.Printf("Golang version: %s\n", goversion)
	os.Exit(0)
}

func defaultConfig() {
	appInstance := app.New(nil)
	appInstance.InitConfig(rootDir)
	fmt.Printf("Default config file %s/config/config.yml was created for customization\n", rootDir)
	os.Exit(0)
}

func reindexStock() {
	svc := initService()
	runningService := false
	svcStatus, err := svc.Status()
	if err == nil && svcStatus == service.StatusRunning {
		svc.Stop()
		runningService = true
	}

	appInstance := app.New(nil)
	cfg := appInstance.InitConfig(rootDir)
	os.Remove(cfg.Database.DSN)

	if runningService {
		svc.Start()
		os.Exit(0)
	}

	appInstance.InitLocales(cfg)

	stockLog, _ := appInstance.InitLogs(cfg, false)
	defer stockLog.Close()

	start := time.Now()
	stockLog.S.Println(">>> Book stock reindex started  >>>>>>>>>>>>>>>>>>>>>>>>>>>")

	db := appInstance.InitDatabase(cfg)
	defer db.Close()
	if !db.IsReady() {
		stockLog.S.Println("Book stock was inited. Tables were created in empty database")
	}

	genresTree := appInstance.InitGenres(cfg)

	bookQueue := make(chan model.Book, cfg.Database.BOOK_QUEUE_SIZE)
	defer close(bookQueue)
	fileQueue := make(chan index.File, cfg.Database.FILE_QUEUE_SIZE)
	defer close(fileQueue)
	stockHandler := appInstance.InitIndexerOnce(cfg, db, genresTree, stockLog)
	stockHandler.StopDB = make(chan struct{})
	defer close(stockHandler.StopDB)
	stockHandler.StopScan = make(chan struct{})
	defer close(stockHandler.StopScan)

	defer func() { stockHandler.StopScan <- struct{}{} }()
	stockHandler.StopScan <- struct{}{}

	stockHandler.StopDB <- struct{}{}
	<-stockHandler.StopDB

	stockLog.S.Println("<<< Book stock reindex finished <<<<<<<<<<<<<<<<<<<<<<<<<<<")
	stockLog.S.Println("Time elapsed: ", time.Since(start))

	os.Exit(0)
}

func run() {
	appInstance := app.New(nil)
	cfg := appInstance.InitConfig(rootDir)

	appInstance.InitLocales(cfg)

	stockLog, opdsLog := appInstance.InitLogs(cfg, true)
	defer stockLog.Close()
	defer opdsLog.Close()

	db := appInstance.InitDatabase(cfg)
	defer db.Close()
	if !db.IsReady() {
		stockLog.S.Println("Book stock was inited. Tables were created in empty database")
	}

	genresTree := appInstance.InitGenres(cfg)

	// Starting OPDS
	opdsHandler, server := appInstance.InitOPDS(cfg, db, genresTree, opdsLog)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			opdsHandler.LOG.E.Println(err)
			close(doShutdown)
		}
	}()
	opdsHandler.LOG.S.Printf("Server started listening at http://localhost:%d \n", cfg.OPDS.PORT)

	// Starting book stock
	bookQueue := make(chan model.Book, cfg.Database.BOOK_QUEUE_SIZE)
	defer close(bookQueue)
	fileQueue := make(chan index.File, cfg.Database.FILE_QUEUE_SIZE)
	defer close(fileQueue)
	stockHandler := appInstance.InitIndexer(cfg, db, genresTree, stockLog)
	stockHandler.StopDB = make(chan struct{})
	defer close(stockHandler.StopDB)
	stockHandler.StopScan = make(chan struct{})
	defer close(stockHandler.StopScan)

	stockHandler.LOG.S.Printf("Book cache warming started...\n")

	go func() {
		defer func() { stockHandler.StopScan <- struct{}{} }()
		dir := cfg.Library.STOCK_DIR
		if len(cfg.Library.NEW_DIR) > 0 {
			dir = cfg.Library.NEW_DIR
		}
		for {
			stockHandler.ScanDir(dir)
			time.Sleep(time.Duration(cfg.Database.POLL_DELAY) * time.Second)
			select {
			case <-stockHandler.StopScan:
				return
			default:
				continue
			}
		}
	}()
	stockHandler.LOG.S.Printf("New acquisitions scanning started...\n")

	// <<<<<<<<<<<<<<<<<- Wait for shutdown
	<-doShutdown

	opdsHandler.LOG.S.Printf("Shutdown started...\n")

	// Stop scanning for new acquisitions and wait for completion
	stockHandler.StopScan <- struct{}{}
	<-stockHandler.StopScan
	stockHandler.LOG.S.Printf("New acquisitions scanning was stoped correctly\n")

	// Stop addind new acquisitions to index and wait for completion
	stockHandler.StopDB <- struct{}{}
	<-stockHandler.StopDB
	stockHandler.LOG.S.Printf("New acquisitions adding was stoped correctly\n")

	// Shutdown OPDS server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		opdsHandler.LOG.E.Printf("Shutdown error: %v\n", err)
	}
	opdsHandler.LOG.S.Printf("Server at %s was shut down correctly\n", fmt.Sprint("http://localhost:", cfg.OPDS.PORT))

}
