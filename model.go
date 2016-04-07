package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type model struct{}

type user struct {
	ip      string
	name    string
	traffic int64
}

var users []user
var restricts []string

func (m *model) Init() {
	db, err := sql.Open("sqlite3", "proxbase.db")
	if err != nil {
		log.Fatal("Невозможно загрузить базу данных!")
	}
	defer db.Close()

	///DEBUG
	//_, err = db.Exec("DROP TABLE IF EXISTS users")
	//fatalError("Невозможно удалить таблицу users", err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`ip` TEXT NOT NULL, `name` TEXT, `traffic` INTEGER NOT NULL)")
	if err != nil {
		log.Fatalf("Невозможно создать таблицу users: %v", err.Error())
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `restricts` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`site` TEXT NOT NULL)")
	if err != nil {
		log.Fatalf("Невозможно создать таблицу restricts: %v", err.Error())
	}

	//TEMP
	ins, err := db.Prepare("INSERT INTO users(ip, name, traffic) values(?,?,?)")
	if err != nil {
		log.Fatalf("Невозможно создать запрос: %v", err.Error())
	}
	_, err = ins.Exec("192.168.57.5", "unknown", 0)
	fatalError("Cannt insert vals", err)
	//if err != nil {
	//	log.Fatalf("Cannt insert vals coz %v", err.Error())
	//}

	rUsers, err := db.Query("SELECT ip, name, traffic FROM users")
	fatalError("Ошибка при использовании SELECT", err)

	for rUsers.Next() {
		var ip, name string
		var traffic int64
		err = rUsers.Scan(&ip, &name, &traffic)
		fatalError("", err)
		users = append(users, user{ip, name, traffic})
	}

	log.Print(users)

}

func fatalError(text string, err error) {
	if err != nil {
		log.Fatal(text + " : " + err.Error())
	}
}
