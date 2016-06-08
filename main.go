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
var base *model

func main() {
	// Просим Go использовать все имеющиеся в системе процессоры.
	runtime.GOMAXPROCS(runtime.NumCPU())

	//logFile, _ := os.Create("logs.txt")
	//serv := http.Server{}

	base = new(model)
	base.Init()
	defer base.Close()

	logger := log.New(os.Stdout, "log: ", log.Ltime)
	logger.Printf("Сервер запущен: %v", time.Now())

	pServ = prx.NewServ(logger, *dbg)
	pServ.Users = initUsers(base)
	//pServ.Grab = base

	//base.AddRestricts("yandex.ru", "reddit.com", "lenta.ru")
	//list := base.GetRestricts()
	//pServ.Restricts = list

	//Функция обновления базы данных раз в несколько минут
	go updateBase(1, base)

	//Инициализация роутера
	router := routerInit() //from controller.go
	go http.ListenAndServe(":800", router)

	logger.Fatal(http.ListenAndServe(":"+*port, pServ))
}

func updateBase(mins int, m *model) {
	timer := time.NewTicker(time.Minute * time.Duration(mins))
	currentMonth := time.Now().Month()
	for {
		select {
		case stat := <-pServ.Stats:
			m.UpdateStat(stat.IP, stat.Name, stat.Site, stat.Bytes, stat.Date)
		case <-timer.C:
			if time.Now().Month() != currentMonth {
				for n := range pServ.Users {
					pServ.Users[n].Traffic = 0
				}
			}
		}
		// <-timer.C
		// for k := range pServ.Users {
		// 	name, tr := pServ.GetUser(k)
		// 	m.UpdateUser(k, name, tr)
		// }
	}
}

func initUsers(m *model) map[string]*prx.User {
	users := m.GetUsers()
	prxUsers := make(map[string]*prx.User)
	for _, v := range users {
		u := new(prx.User)
		u.Name = v.Name
		u.Limit = v.Limit
		u.Traffic = v.Traffic
		prxUsers[v.IP] = u
	}
	return prxUsers
}
