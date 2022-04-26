package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/kardianos/service"
	"github.com/vinser/flibgolite/pkg/opds"
	"github.com/vinser/flibgolite/pkg/stock"
	"golang.org/x/text/message"
)

var logger service.Logger

func (h *Handler) ServiceControl(controlAction string) service.Service {
	serviceCfg := &service.Config{}
	serviceCfg.Name = "FLibGoLiteService"
	serviceCfg.DisplayName = "FLibGoLite Service for Linux"
	serviceCfg.Description = "FLibGoLite service controls new aqusitions scan and opds server"
	switch runtime.GOOS {
	case "linux":
		serviceCfg.Dependencies = []string{"Requires=network.target", "After=network-online.target syslog.target"}
		options := make(service.KeyValue)
		options["Restart"] = "on-success"
		options["SuccessExitStatus"] = "1 2 8 SIGKILL"
		serviceCfg.Option = options
	case "windows":
	default:
		log.Printf("FLibGoLite on %s is not available yet", runtime.GOOS)
		return nil
	}

	s, err := service.New(h, serviceCfg)
	if err != nil {
		log.Fatal(err)
	}

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(controlAction) != 0 {
		err := service.Control(s, controlAction)
		if err != nil {
			log.Printf(`Valid actions for "-service" option are: %q\n`, service.ControlAction)
			log.Fatal(err)
		}
		return nil
	}
	h.Exit = make(chan struct{})
	h.S_Exit = make(chan struct{})
	h.O_Exit = make(chan struct{})

	return s
}

func (h *Handler) Start(s service.Service) error {
	go h.run()
	return nil
}

func (h *Handler) run() {
	stockHandler := &stock.Handler{
		CFG: h.CFG,
		LOG: h.S_LOG,
		DB:  h.DB,
		GT:  h.GT,
	}
	stockHandler.InitStockFolders()
	go func() {
		defer func() { h.S_Exit <- struct{}{} }()
		f := "New aquisitions scanning started...\n"
		h.S_LOG.I.Printf(f)
		logger.Info(f)
		for {
			stockHandler.ScanDir(h.CFG.Library.NEW_ACQUISITIONS)
			time.Sleep(time.Duration(h.CFG.Database.POLL_DELAY) * time.Second)
			select {
			case <-h.S_Exit:
				return
			default:
				continue
			}
		}
	}()

	opdsHandler := &opds.Handler{
		CFG: h.CFG,
		LOG: h.O_LOG,
		DB:  h.DB,
		GT:  h.GT,
		P:   message.NewPrinter(*h.LANG),
	}
	server := &http.Server{
		Addr:    fmt.Sprint(":", h.CFG.OPDS.PORT),
		Handler: opdsHandler,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	f := "Server at %s is listening...\n"
	h.O_LOG.I.Printf(f, fmt.Sprint("http://localhost:", h.CFG.OPDS.PORT))
	logger.Infof(f, fmt.Sprint("http://localhost:", h.CFG.OPDS.PORT))
	h.Server = server
}

func (h *Handler) Stop(s service.Service) error {
	f := "Shutdown started...\n"
	h.O_LOG.I.Printf(f)
	logger.Info(f)

	// Stop scanning for new aquisitions and wait for completion
	h.S_Exit <- struct{}{}
	<-h.S_Exit
	f = "New aquisitions scanning was stoped correctly\n"
	h.S_LOG.I.Printf(f)
	logger.Info(f)

	// Shutdown OPDS server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.Server.Shutdown(ctx); err != nil {
		f := "Shutdown error: %v\n"
		h.O_LOG.E.Printf(f, err)
		logger.Errorf(f, err)
	}
	f = "Server at %s was shut down correctly\n"
	h.O_LOG.I.Printf(f, fmt.Sprint("http://localhost:", h.CFG.OPDS.PORT))
	logger.Infof(f, fmt.Sprint("http://localhost:", h.CFG.OPDS.PORT))

	close(h.Exit)
	return nil
}
