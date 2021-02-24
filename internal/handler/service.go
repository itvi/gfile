package handler

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/kardianos/service"
)

var logger service.Logger

type program struct {
	port, dir string
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	logger.Infof("Arguments: %v", os.Args)
	logger.Infof("I'm running %v.", service.Platform())

	cfg, db := Config(p.dir)
	defer db.Close()

	server := &http.Server{
		Addr:    p.port,
		Handler: cfg.Route(),
	}

	log.Printf("Starting port [%s], From [%s]", server.Addr, p.dir)
	log.Fatal(server.ListenAndServe())
}

func (p *program) Stop(s service.Service) error {
	logger.Info("I'm Stopping!")
	return nil
}

// Service setup
/* -------usage-----------------------------------
   1. go build
   2. gfile -s install -p :9000 -d g:/test
   3. gfile -s start

   如果要改变端口或目录，则需uninstall后再重新install。
--------------------------------------------------*/
func Service() {
	portFlag := flag.String("p", ":8000", "TCP address port")
	dirFlag := flag.String("d", ".", "Monitor directory") // such as: c:/test,d:/test ...
	svcFlag := flag.String("s", "", "Control the system service.")
	flag.Parse()

	// define service config
	svcConfig := &service.Config{
		Name:        "GFile",
		DisplayName: "GFile service",
		Description: "Download center.",
	}

	// add another arguments
	if *svcFlag == "install" {
		var prgArgs []string
		if len(*portFlag) != 0 {
			prgArgs = append(prgArgs, "-p="+*portFlag)
		}
		if len(*dirFlag) != 0 {
			prgArgs = append(prgArgs, "-d="+*dirFlag)
		}
		svcConfig.Arguments = prgArgs
	}

	prg := &program{
		port: *portFlag,
		dir:  *dirFlag,
	}

	// create service
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// setup logger
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	// handle service controls
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	// run the service
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
