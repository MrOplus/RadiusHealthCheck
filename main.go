package main

import (
	"fmt"
	"github.com/kooroshh/RadiusHealthCheck/cli"
	"github.com/kooroshh/RadiusHealthCheck/config"
	"github.com/kooroshh/RadiusHealthCheck/healthcheck"
	"os"
	"os/signal"
	"syscall"
)
func main() {
	ch :=make(chan os.Signal ,1)
	signal.Notify(ch,syscall.SIGTERM,syscall.SIGKILL,syscall.SIGINT)
	opts := cli.Parse()
	conf := config.Parse(opts.Config)
	fmt.Printf("Total Servers %d\n", len(conf.Servers))

	for i := 0 ; i < len(conf.Servers) ; i++ {
		go healthcheck.StartHealthCheck(&conf.Servers[i],conf)
	}
	_ = <-ch
	fmt.Printf("Exiting ...")
}
