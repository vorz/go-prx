package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type model struct{}

//Site - Структура используется при формировании выборки из бд stats
type Site struct {
	SiteName string `json:"sitename"`
	Traffic  int64  `json:"traffic"`
}

//User - Структура используется при формировании выборки из бд users
type User struct {
	IP      string
	Name    string
	Traffic int64
}

var db *sql.DB

//Большинство ошибок с бд можно считать фатальными, при которой
//невозможно продолжение работы (log.Fatal завершает программу),
//поэтому обработку ошибок завернем в функцию
func fatalError(text string, err error) {
	if err != nil {
		log.Fatal(text + " : " + err.Error())
	}
}

func (m *model) Init() {
	var err error
	db, err = sql.Open("sqlite3", "proxbase.db")
	fatalError("Невозможно загрузить базу данных!", err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`ip` TEXT NOT NULL, `name` TEXT, `traffic` INTEGER NOT NULL)")
	fatalError("Невозможно создать таблицу users: %v", err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `stats` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`site` TEXT NOT NULL, `user_id` TEXT NOT NULL, `bytes` INTEGER NOT NULL)")
	fatalError("Невозможно создать таблицу restricts: %v", err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `restricts` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`site` TEXT NOT NULL)")
	fatalError("Невозможно создать таблицу restricts: %v", err)
}

func (m *model) Close() {
	db.Close()
}

func (m *model) AddRestricts(sites ...string) {
	for _, str := range sites {
		var id int
		err := db.QueryRow("SELECT id FROM restricts WHERE site=?", str).Scan(&id)
		switch {
		case err == sql.ErrNoRows:
			_, e := db.Exec("INSERT INTO restricts(site) values(?)", str)
			fatalError("Невозможно обновить базу данных", e)
		case err != nil:
			fatalError("Ошибка доступа к базе данных", err)
		default:
			continue
		}
	}
}

func (m *model) GetRestricts() []string {
	r, err := db.Query("SELECT site FROM restricts")
	fatalError("Ошибка при использовании SELECT", err)

	var list []string

	for r.Next() {
		var site string
		r.Scan(&site)
		list = append(list, site)
	}
	return list
}

func (m *model) GetUserId(ip string) int {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE ip=?", ip).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return -1
	case err != nil:
		log.Fatal(err)
	}
	return id
}

func (m *model) GetUsers() []User {
	rUsers, err := db.Query("SELECT ip, name, traffic FROM users ORDER BY traffic DESC")
	fatalError("Ошибка при использовании SELECT", err)

	var users []User

	for rUsers.Next() {
		var u User
		rUsers.Scan(&u.IP, &u.Name, &u.Traffic)
		users = append(users, u)
	}
	return users
}

func (m *model) GetTraffic(id int) int64 {
	var traffic int64
	err := db.QueryRow("SELECT traffic FROM users WHERE id=?", id).Scan(&traffic)
	switch {
	case err == sql.ErrNoRows:
		return -1
	case err != nil:
		log.Fatal(err)
	}
	return traffic
}

func (m *model) UpdateUser(ip string, name string, traffic int64) bool {
	if m.GetUserId(ip) >= 0 {
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

func (m *model) UpdateStat(ip string, site string, bytes int64) {
	if userid := m.GetUserId(ip); userid >= 0 {
		var id int
		var b int64
		err := db.QueryRow("SELECT id, bytes FROM stats WHERE site=? AND user_id=?", site, userid).Scan(&id, &b)
		switch {
		case err == sql.ErrNoRows:
			_, e := db.Exec("INSERT INTO stats(site, user_id, bytes) values(?,?,?)", site, userid, bytes)
			fatalError("Невозможно обновить базу данных", e)
		case err != nil:
			fatalError("Ошибка доступа к базе данных", err)
		default:
			b += bytes
			_, e := db.Exec("UPDATE stats SET bytes=? WHERE id=?", b, id)
			fatalError("Невозможно обновить базу данных", e)
		}
	}
}

//Получить список сайтов и трафик по id пользователя
func (m *model) GetSitesStats(id int) []Site {
	rSites, err := db.Query("SELECT site, bytes FROM stats WHERE user_id=? ORDER BY bytes DESC", id)
	fatalError("Ошибка при использовании SELECT", err)

	var sites []Site

	for rSites.Next() {
		var s Site
		rSites.Scan(&s.SiteName, &s.Traffic)
		sites = append(sites, s)
	}

	return sites
}
