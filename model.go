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

var db *sql.DB

func (m *model) Init() {
	base, err := sql.Open("sqlite3", "proxbase.db")
	fatalError("Невозможно загрузить базу данных!", err)
	//defer base.Close()

	db = base

	///DEBUG
	//_, err = db.Exec("DROP TABLE IF EXISTS users")
	//fatalError("Невозможно удалить таблицу users", err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`ip` TEXT NOT NULL, `name` TEXT, `traffic` INTEGER NOT NULL)")
	fatalError("Невозможно создать таблицу users: %v", err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `restricts` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`site` TEXT NOT NULL)")
	fatalError("Невозможно создать таблицу restricts: %v", err)

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

	//	log.Print(users)

}

func (m *model) Close() {
	db.Close()
}

func (m *model) UserExists(ip string) bool {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE ip=?", ip).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Fatal(err)
	}
	return true
}

func (m *model) GetUsers() []string {
	rUsers, err := db.Query("SELECT ip FROM users")
	fatalError("Ошибка при использовании SELECT", err)

	var list []string

	for rUsers.Next() {
		var ip string
		rUsers.Scan(&ip)
		list = append(list, ip)
	}
	return list
}

func (m *model) UpdateUser(ip string, name string, traffic int64) bool {
	if m.UserExists(ip) {
		r, err := db.Exec("UPDATE users SET name=?, traffic=? WHERE ip=?", name, traffic, ip)
		fatalError("Невозможно обновить базу данных", err)
		num, _ := r.RowsAffected()
		if num > 0 {
			return true
		}
	} else {
		r, err := db.Exec("INSERT INTO users(ip, name, traffic) values(?,?,?)", ip, name, traffic)
		fatalError("Невозможно обновить базу данных", err)
		num, _ := r.RowsAffected()
		if num > 0 {
			return true
		}
	}
	return false
}

func fatalError(text string, err error) {
	if err != nil {
		log.Fatal(text + " : " + err.Error())
	}
}
