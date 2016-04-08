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

var pServ *prx.ProxyServ

func main() {
	// Просим Go использовать все имеющиеся в системе процессоры.
	runtime.GOMAXPROCS(runtime.NumCPU())

	//logFile, _ := os.Create("logs.txt")
	//serv := http.Server{}

	base := new(model)
	base.Init()
	defer base.Close()

	logger := log.New(os.Stdout, "log: ", log.Ltime)
	logger.Printf("Сервер запущен: %v", time.Now())

	pServ = prx.NewServ(logger, *dbg)

	go updateBase(1, base)

	logger.Fatal(http.ListenAndServe(":"+*port, pServ))
}

func updateBase(mins int, m *model) {
	timer := time.NewTicker(time.Minute * time.Duration(mins))
	for {
		<-timer.C
		for k := range pServ.Users {
			name, tr := pServ.GetUser(k)
			m.UpdateUser(k, name, tr)
			lst := m.GetUsers()
			log.Print("##### USERS #######")
			log.Print(lst)
		}
	}
}
