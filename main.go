package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/vorz/go-prx/prx"
)

var (
	port = flag.String("port", "8080", "Port for proxy")
	dbg  = flag.Bool("d", false, "Debug information")
)

func main() {
	// Просим Go использовать все имеющиеся в системе процессоры.
	runtime.GOMAXPROCS(runtime.NumCPU())

	//logFile, _ := os.Create("logs.txt")
	//serv := http.Server{}

	base := new(model)
	base.Init()

	logger := log.New(os.Stdout, "log: ", log.Ltime)
	logger.Printf("Сервер запущен: %v", time.Now())

	pServ := &prx.ProxyServ{Log: logger}

	logger.Fatal(http.ListenAndServe(":"+*port, pServ))
}
