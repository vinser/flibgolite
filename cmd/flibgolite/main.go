package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/kardianos/service"
	"github.com/vinser/flibgolite/pkg/config"
	"github.com/vinser/flibgolite/pkg/database"
	"github.com/vinser/flibgolite/pkg/genres"
	"github.com/vinser/flibgolite/pkg/opds"
	"github.com/vinser/flibgolite/pkg/stock"
	"golang.org/x/text/message"
)

var version, buildTime, target, goversion string

func main() {
	serviceFlag := flag.String("service", "", `control FLibGoLite system service`)
	reindexFlag := flag.Bool("reindex", false, `empty book stock database and then scan book stock directory to add books to database`)
	configFlag := flag.Bool("config", false, `create default config file in ./config folder for customization and exit`)
	helpFlag := flag.Bool("help", false, `display extended command help and exit`)
	versionFlag := flag.Bool("version", false, `output version information and exit`)
	flag.Parse()
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
	fmt.Printf(
		`	
FLibGoLite is multiplatform lightweight OPDS server with SQLite database book search index
This program was build for %s-%s

Usage: flibgolite [OPTION]

With no OPTION program will run in console mode (Ctrl+C to exit)
Caution: Only one OPTION can be used at a time

OPTION should be one of:
  -service [action]     control FLibGoLite system service
	  where action is one of: install, start, stop, restart, uninstall, status 
  -reindex              empty book stock index and then scan book stock directory to add books to index (database)
  -config               create default config file in ./config folder for customization and exit
  -help                 display this help and exit
  -version              output version information and exit

Examples:
  ./flibgolite                      Run FLibGoLite in console mode
  ./flibgolite -service install     Install FLibGoLite as a system service

Documentation at: <https://github.com/vinser/flibgolite>

`,
		runtime.GOOS, runtime.GOARCH)
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
	config.LoadConfig()
	fmt.Println(`Default config file "./config/config.yml" was created for customization`)
	os.Exit(0)
}

func reindexStock() {
	cfg := config.LoadConfig()
	svc := initService()
	svcStatus, err := svc.Status()
	if err == nil && svcStatus == service.StatusRunning {
		svc.Stop()
	}
	stockLog, _ := cfg.InitLogs(false)
	defer stockLog.Close()

	db := database.NewDB(cfg.Database.DSN)
	defer db.Close()
	if !db.IsReady() {
		db.InitDB()
		f := "Book stock was inited. Tables were created in empty database"
		stockLog.S.Println(f)
	}

	genresTree := genres.NewGenresTree(cfg.Genres.TREE_FILE)

	stockHandler := &stock.Handler{
		CFG: cfg,
		LOG: stockLog,
		DB:  db,
		GT:  genresTree,
	}
	stockHandler.InitStockFolders()
	stockHandler.Reindex()
	if err == nil && svcStatus == service.StatusRunning {
		svc.Start()
	}
	os.Exit(0)
}

func run() {
	cfg := config.LoadConfig()

	cfg.Locales.LoadLocales()

	stockLog, opdsLog := cfg.InitLogs(true)
	defer stockLog.Close()
	defer opdsLog.Close()

	db := database.NewDB(cfg.Database.DSN)
	defer db.Close()
	if !db.IsReady() {
		db.InitDB()
		stockLog.S.Println("Book stock was inited. Tables were created in empty database")
	}

	genresTree := genres.NewGenresTree(cfg.Genres.TREE_FILE)

	stockHandler := &stock.Handler{
		CFG: cfg,
		LOG: stockLog,
		DB:  db,
		GT:  genresTree,
	}
	stockHandler.InitStockFolders()
	stockHandler.SY.Stop = make(chan struct{})
	defer close(stockHandler.SY.Stop)
	go func() {
		defer func() { stockHandler.SY.Stop <- struct{}{} }()
		dir := cfg.Library.STOCK_DIR
		if len(cfg.Library.NEW_DIR) > 0 {
			dir = cfg.Library.NEW_DIR
		}
		for {
			stockHandler.ScanDir(dir)
			time.Sleep(time.Duration(cfg.Database.POLL_DELAY) * time.Second)
			select {
			case <-stockHandler.SY.Stop:
				return
			default:
				continue
			}
		}
	}()
	stockHandler.LOG.S.Printf("New acquisitions scanning started...\n")

	opdsHandler := &opds.Handler{
		CFG: cfg,
		LOG: opdsLog,
		DB:  db,
		GT:  genresTree,
		MP:  make(map[string]*message.Printer, len(cfg.Locales.Languages)),
	}

	for k, v := range cfg.Locales.Languages {
		opdsHandler.MP[k] = message.NewPrinter(v.Tag)
	}
	server := &http.Server{
		Addr:    fmt.Sprint(":", cfg.OPDS.PORT),
		Handler: opdsHandler,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	opdsHandler.LOG.S.Printf("Server started listening at %s \n", fmt.Sprint("http://localhost:", cfg.OPDS.PORT))

	// <<<<<<<<<<<<<<<<<- Wait for shutdown
	<-doShutdown

	opdsHandler.LOG.S.Printf("Shutdown started...\n")

	// Stop scanning for new acquisitions and wait for completion
	stockHandler.SY.Stop <- struct{}{}
	<-stockHandler.SY.Stop
	stockHandler.LOG.S.Printf("New acquisitions scanning was stoped correctly\n")

	// Shutdown OPDS server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		opdsHandler.LOG.E.Printf("Shutdown error: %v\n", err)
	}
	opdsHandler.LOG.S.Printf("Server at %s was shut down correctly\n", fmt.Sprint("http://localhost:", cfg.OPDS.PORT))

}
