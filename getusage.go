package main

import (
	"flag"
	"log"
	"encoding/xml"
	"net/http"
	"strconv"
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var logger service.Logger

type program struct {
	exit chan struct{}
}

type CpuInfo struct {
	Hostname         string
	CPU_Percent      string
	Total_Ram        string
	Available_Ram    string
	Used_Ram_Percent string
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func index(w http.ResponseWriter, r *http.Request) {

	v, _ := mem.VirtualMemory()
	h, _ := host.Info()
	c2, _ := cpu.Percent(0, true)
	cpup := strconv.FormatFloat(c2[0], 'f', 4, 64)
	tram := strconv.FormatUint(v.Total, 10)
	aram := strconv.FormatUint(v.Available, 10)
	pram := strconv.FormatFloat(v.UsedPercent, 'f', 2, 64)
	CpuInfo := &CpuInfo{h.Hostname, cpup, tram, aram, pram}
	buf, _ := xml.Marshal(CpuInfo)
	w.Header().Set("Content-Type", "application/xml")
	w.Write(buf)
}

func (p *program) run() error {
	http.HandleFunc("/getusage", index)
	http.ListenAndServe(":8083", nil)
	return nil
}
func (p *program) Stop(s service.Service) error {
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "XML server",
		DisplayName: "Http://localhost:8083/getusage",
		Description: "web server to return CPU usage",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
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

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
