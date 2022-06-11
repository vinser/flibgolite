package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/kardianos/service"
)

var svcLog service.Logger

type program struct{}

var doShutdown chan struct{}

func (p *program) Start(s service.Service) error {
	doShutdown = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) run() {
	svcLog.Infof("Service is running %v.", service.Platform())
	run()
}

func (p *program) Stop(s service.Service) error {
	svcLog.Infof("%s service is stopping!", s)
	close(doShutdown)
	return nil
}

func initService() service.Service {
	serviceCfg := &service.Config{}
	serviceCfg.Name = "FLibGoLite"
	serviceCfg.DisplayName = "FLibGoLite Service"
	serviceCfg.Description = "FLibGoLite serves new books acquisitions scan and opds http server"
	switch runtime.GOOS {
	case "linux":
		serviceCfg.Dependencies = []string{"Requires=network.target", "After=network-online.target syslog.target"}
		if suid := os.Getenv("SUDO_USER"); suid != "" {
			serviceCfg.UserName = suid
		}
		options := make(service.KeyValue)
		options["Restart"] = "on-success"
		options["SuccessExitStatus"] = "1 2 8 SIGKILL"
		serviceCfg.Option = options
	case "windows":
	case "darwin":
	default:
		log.Fatalf("FLibGoLite on %s is not available yet\n", runtime.GOOS)
	}

	p := &program{}

	s, err := service.New(p, serviceCfg)
	if err != nil {
		log.Fatalln("Failed to instantinate service:", err)
	}

	errs := make(chan error, 5)
	svcLog, err = s.Logger(errs)
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

	return s
}

func controlService(action string) {
	s := initService()
	switch action {
	case "install":
		err := s.Install()
		if err != nil {
			log.Fatalln("Error installing the service:", err)
		}
		fmt.Println("Service installed!")
	case "uninstall":
		err := s.Uninstall()
		if err != nil {
			log.Fatalln("Error uninstalling the service:", err)
		}
		fmt.Println("Service uninstalled!")
	case "start":
		err := s.Start()
		if err != nil {
			log.Fatalln("Error starting the service:", err)
		}
		fmt.Println("Service started!")
	case "stop":
		err := s.Stop()
		if err != nil {
			log.Fatalln("Error stopping the service:", err)
		}
		fmt.Println("Service stopped!")
	case "restart":
		err := s.Restart()
		if err != nil {
			log.Fatalln("Error restarting the service:", err)
		}
		fmt.Println("Service restarted!")
	case "status":
		status, err := s.Status()
		if err != nil {
			log.Fatalln("Error getting the status of the service:", err)
		}
		fmt.Print("Status of the service: ")
		switch status {
		case service.StatusRunning:
			fmt.Println("running")
		case service.StatusStopped:
			fmt.Println("stopped")
		case service.StatusUnknown:
			fmt.Println("unknown")
		}
	default:
		log.Fatalf("Unknown action '%s'\n", action)
	}
	os.Exit(0)
}

func runService(s service.Service) {
	err := s.Run()
	if err != nil {
		svcLog.Error(err)
	}
}
