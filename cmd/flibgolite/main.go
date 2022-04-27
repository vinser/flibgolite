package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"

	"github.com/kardianos/service"
	"github.com/vinser/flibgolite/pkg/config"
	"github.com/vinser/flibgolite/pkg/database"
	"github.com/vinser/flibgolite/pkg/genres"
	"github.com/vinser/flibgolite/pkg/rlog"
	"github.com/vinser/flibgolite/pkg/stock"

	"golang.org/x/text/language"
)

type Handler struct {
	CFG    *config.Config
	DB     *database.DB
	GT     *genres.GenresTree
	LANG   *language.Tag
	Exit   chan struct{}
	S_Exit chan struct{}
	O_Exit chan struct{}
	S_LOG  *rlog.Log
	O_LOG  *rlog.Log
	Server *http.Server
}

func (h *Handler) reindex() {
	stockHandler := &stock.Handler{
		CFG: h.CFG,
		LOG: h.S_LOG,
		DB:  h.DB,
		GT:  h.GT,
	}
	stockHandler.InitStockFolders()
	stockHandler.Reindex()
}

var serviceActions []string

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
	  actions are: %q 
  -reindex              empty book stock index and then scan book stock directory to add books to index (database)
  -config               create default config file in ./config folder for customization and exit
  -help                 display this help and exit
  -version              output version information and exit

Examples:
  ./flibgolite                      Run FLibGoLite in console mode
  ./flibgolite -service install     Install FLibGoLite as a system service

Documentation at: <https://github.com/vinser/flibgolite>

`,
		runtime.GOOS, runtime.GOARCH, serviceActions)
}

func displayVersion() {}

func main() {
	serviceActions = append(service.ControlAction[:], "status")
	serviceFlag := flag.String("service", "", `control FLibGoLite system service`)
	reindexFlag := flag.Bool("reindex", false, `empty book stock database and then scan book stock directory to add books to database`)
	configFlag := flag.Bool("config", false, `create default config file in ./config folder for customization and exit`)
	helpFlag := flag.Bool("help", false, `display extended command help and exit`)
	versionFlag := flag.Bool("version", false, `output version information and exit`)
	flag.Parse()
	if flag.NFlag() > 1 {
		fmt.Println(`Error: More than one OPTION used`)
		displayHelp()
		return
	}

	if *helpFlag {
		displayHelp()
		return
	}

	if *versionFlag {
		displayVersion()
		return
	}

	cfg := config.LoadConfig()
	if *configFlag {
		fmt.Println(`Default config file "./config/flibgolite.yml" was created for customization`)
		return
	}
	cfg.LoadLocales()
	langTag := language.Make(cfg.Language.DEFAULT)

	stockLog := rlog.NewLog(cfg.Logs.SCAN, cfg.Logs.DEBUG)
	defer stockLog.File.Close()
	opdsLog := rlog.NewLog(cfg.Logs.OPDS, cfg.Logs.DEBUG)
	defer opdsLog.File.Close()

	db := database.NewDB(cfg.Database.DSN)
	defer db.Close()
	if !db.IsReady() {
		db.InitDB()
		f := "Book stock was inited. Tables were created in empty database"
		stockLog.I.Println(f)
	}

	genresTree := genres.NewGenresTree(cfg.Genres.TREE_FILE)

	h := &Handler{
		CFG:   cfg,
		DB:    db,
		GT:    genresTree,
		LANG:  &langTag,
		S_LOG: stockLog,
		O_LOG: opdsLog,
	}

	if *reindexFlag {
		h.reindex()
		return
	}

	if s := h.ServiceControl(*serviceFlag); s != nil {
		err := s.Run()
		if err != nil {
			logger.Error(err)
		}
	}
}
