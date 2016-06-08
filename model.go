package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const defaultLimit = 1000000000

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
	Limit   int64
}

//Log - структура для передачи "сырых" данных статистики
type Log struct {
	SiteName string
	IP       string
	Name     string
	Traffic  int64
	Date     string
}

//Stats - структура содержит общую статистику базы данных (сервера)
type Stats struct {
	UsersNum int
	Traffic  int64
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
		"`ip` TEXT NOT NULL, `name` TEXT NOT NULL, `lim` INTEGER NOT NULL)")
	fatalError("Невозможно создать таблицу users: %v", err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `sites` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`name` TEXT NOT NULL)")
	fatalError("Невозможно создать таблицу sites: %v", err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `stats` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`user_id` INTEGER NOT NULL, site_id INTEGER NOT NULL, `bytes` INTEGER NOT NULL, `date` INTEGER NOT NULL, " +
		"FOREIGN KEY(user_id) REFERENCES users(id), FOREIGN KEY(site_id) REFERENCES sites(id))")
	fatalError("Невозможно создать таблицу stats: %v", err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `restricts` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`site_id` INTEGER NOT NULL, FOREIGN KEY(site_id) REFERENCES sites(id))")
	fatalError("Невозможно создать таблицу restricts: %v", err)
}

func (m *model) Close() {
	db.Close()
}

func (m *model) AddRestricts(sites ...string) {
	for _, str := range sites {
		var id int
		err := db.QueryRow("SELECT id FROM restricts WHERE site_id=?", str).Scan(&id)
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
	defer r.Close()

	var list []string

	for r.Next() {
		var site string
		r.Scan(&site)
		list = append(list, site)
	}
	return list
}

func (m *model) GetOverall() Stats {
	var stat Stats
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stat.UsersNum)
	switch {
	case err == sql.ErrNoRows:
		stat.UsersNum = 0
	case err != nil:
		log.Fatal(err)
	}
	err = db.QueryRow("SELECT SUM(bytes) FROM stats WHERE strftime('%m', stats.date,'unixepoch')=strftime('%m', 'now')").Scan(&stat.Traffic)
	switch {
	case err == sql.ErrNoRows:
		stat.Traffic = 0
	case err != nil:
		log.Fatal(err)
	}

	return stat
}

func (m *model) GetUserID(ip string) int {
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

func (m *model) AddUser(ip string, name string, limit int64) bool {
	r, err := db.Exec("INSERT INTO users(ip, name, lim) values(?,?,?)", ip, name, limit)
	fatalError("Невозможно обновить базу данных users", err)
	num, _ := r.RowsAffected()
	if num > 0 {
		return true
	}
	return false
}

func (m *model) GetSiteID(name string) int {
	var id int
	err := db.QueryRow("SELECT id FROM sites WHERE name=?", name).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return -1
	case err != nil:
		log.Fatal(err)
	}
	return id
}

func (m *model) AddSite(name string) bool {
	r, err := db.Exec("INSERT INTO sites(name) values(?)", name)
	fatalError("Невозможно обновить базу данных sites", err)
	num, _ := r.RowsAffected()
	if num > 0 {
		return true
	}
	return false
}

func (m *model) GetUsers() []*User {
	rUsers, err := db.Query("SELECT id FROM users")
	fatalError("Ошибка при использовании SELECT", err)
	defer rUsers.Close()

	var users []*User

	for rUsers.Next() {
		var id int
		rUsers.Scan(&id)
		u := m.GetUserInfo(id)
		users = append(users, u)
	}
	return users
}

func (m *model) GetUserInfo(id int) *User {
	var usr User
	err := db.QueryRow("SELECT ip, name, lim FROM users WHERE id=?", id).Scan(&usr.IP, &usr.Name, &usr.Limit)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		log.Fatal(err)
	}
	//Трафик за текущий месяц
	err = db.QueryRow("SELECT SUM(bytes) FROM stats WHERE strftime('%m', stats.date,'unixepoch')=strftime('%m', 'now') AND user_id=?", id).Scan(&usr.Traffic)
	switch {
	case err == sql.ErrNoRows:
		usr.Traffic = 0
	case err != nil:
		log.Fatal(err)
	}
	return &usr
}

// func (m *model) GetTraffic(id int) int64 {
// 	var traffic int64
// 	err := db.QueryRow("SELECT traffic FROM users WHERE id=?", id).Scan(&traffic)
// 	switch {
// 	case err == sql.ErrNoRows:
// 		return -1
// 	case err != nil:
// 		log.Fatal(err)
// 	}
// 	return traffic
// }

// func (m *model) UpdateUser(ip string, name string, traffic int64) bool {
// 	if m.GetUserID(ip) >= 0 {
// 		r, err := db.Exec("UPDATE users SET name=?, traffic=? WHERE ip=?", name, traffic, ip)
// 		fatalError("Невозможно обновить базу данных", err)
// 		num, _ := r.RowsAffected()
// 		if num > 0 {
// 			return true
// 		}
// 	} else {
// 		r, err := db.Exec("INSERT INTO users(ip, name, traffic) values(?,?,?)", ip, name, traffic)
// 		fatalError("Невозможно обновить базу данных", err)
// 		num, _ := r.RowsAffected()
// 		if num > 0 {
// 			return true
// 		}
// 	}
// 	return false
// }

func (m *model) UpdateStat(ip string, dns string, site string, bytes int64, date int64) {

	userid := m.GetUserID(ip)
	if userid < 0 {
		if !m.AddUser(ip, dns, defaultLimit) {
			return
		}
		userid = m.GetUserID(ip)
	}

	siteid := m.GetSiteID(site)
	if siteid < 0 {
		if !m.AddSite(site) {
			return
		}
		siteid = m.GetSiteID(site)
	}

	_, e := db.Exec("INSERT INTO stats(user_id, site_id, bytes, date) values(?,?,?,?)", userid, siteid, bytes, date)
	fatalError("Невозможно обновить базу данных stats", e)

}

//Получить список сайтов и трафик по id пользователя
func (m *model) GetSitesStats(id int) []Site {
	var rSites *sql.Rows
	var err Error
	if id < 0 {
		rSites, err = db.Query("SELECT sites.name, SUM(stats.bytes) AS bt FROM stats, sites "+
			"WHERE sites.id = stats.site_id AND strftime('%m', stats.date,'unixepoch')=strftime('%m', 'now') GROUP BY stats.site_id ORDER BY bt DESC", id)
		fatalError("Ошибка при использовании SELECT", err)
	} else {
		rSites, err = db.Query("SELECT sites.name, SUM(stats.bytes) AS bt FROM stats, sites "+
			"WHERE sites.id = stats.site_id AND stats.user_id=? AND strftime('%m', stats.date,'unixepoch')=strftime('%m', 'now') GROUP BY stats.site_id ORDER BY bt DESC", id)
		fatalError("Ошибка при использовании SELECT", err)
	}
	defer rSites.Close()

	var sites []Site

	for rSites.Next() {
		var s Site
		rSites.Scan(&s.SiteName, &s.Traffic)
		sites = append(sites, s)
	}

	return sites
}

func (m *model) GetRawStats(limit int) []Log {
	rStats, err := db.Query("SELECT sites.name, users.ip, users.name, stats.bytes, datetime(stats.date, 'unixepoch') AS dt "+
		"FROM sites, users, stats WHERE sites.id=stats.site_id AND users.id=stats.user_id ORDER BY dt DESC LIMIT ?", limit)
	fatalError("Ошибка при использовании SELECT", err)
	defer rStats.Close()

	var log []Log

	for rStats.Next() {
		var l Log
		rStats.Scan(&l.SiteName, &l.IP, &l.Name, &l.Traffic, &l.Date)
		log = append(log, l)
	}

	return log
}

/*

SELECT sites.name, users.ip, stats.bytes/1000 FROM sites, users, stats WHERE sites.id=stats.site_id AND users.id=stats.user_id AND strftime('%m', stats.date,'unixepoch')=strftime('%m', 'now')

*/
