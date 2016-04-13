package main

import (
	"errors"
	"html/template"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var templates *template.Template

type Pages struct {
	PagesNum   int
	Current    int
	IndexStart int
	IndexEnd   int
	Prev       int
	Next       int
}

//Ининциализация роутера
func routerInit() *httprouter.Router {
	router := httprouter.New()

	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false

	templates = parseTemplates()
	router.GET("/", IndexHandle)
	router.GET("/overall", StatsHandler)
	router.GET("/stat/:id", SiteStatHandler)
	router.ServeFiles("/css/*filepath", http.Dir("templates/css"))

	return router
}

//GetUserIP - Получить с какого адреса стучатся на сервер
func GetUserIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); len(ip) > 0 {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

//Парсинг и кеширование html файлов из папки templates
func parseTemplates() *template.Template {
	result := template.New("templates")

	basePath := "templates"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()

	templatePathsRaw, _ := templateFolder.Readdir(-1)
	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
				basePath+"/"+pathInfo.Name())
		}
	}

	template.Must(result.ParseFiles(*templatePaths...))
	return result
}

//IndexHandle Стандартная функция обработчик GET /
func IndexHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var data struct {
		IP      string
		Traffic string
		UserID  int
	}

	data.IP = GetUserIP(r)

	data.UserID = base.GetUserId(data.IP)
	if data.UserID > 0 {
		data.Traffic = strconv.FormatInt(int64(base.GetTraffic(data.UserID)/1000), 10)
	} else {
		data.Traffic = "Пользователь не зарегистрирован"
	}

	err := templates.Lookup("index.html").Execute(w, data)
	if err != nil {
		http.NotFound(w, r)
	}

}

//StatsHandler обработчик GET /overall
func StatsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data struct {
		Conns    int
		UsersNum int
		Users    []User
		Page     Pages
	}

	data.Conns = pServ.GetConnNumber()
	data.Users = base.GetUsers()
	data.UsersNum = len(data.Users)

	err := pageInit(data.UsersNum, 10, &data.Page, r)
	if err != nil {
		http.NotFound(w, r)
	}

	err = templates.Lookup("stats.html").Execute(w, data)
	if err != nil {
		http.NotFound(w, r)
	}
}

//SiteStatHandler обработчик GET /stat/id
func SiteStatHandler(w http.ResponseWriter, r *http.Request, id httprouter.Params) {
	var data struct {
		IP     string
		UserID int
		Sites  []Site
		Page   Pages
	}
	data.IP = GetUserIP(r)
	var err error
	data.UserID, err = strconv.Atoi(id.ByName("id"))
	if err != nil {
		http.NotFound(w, r)
	}

	data.Sites = base.GetSitesStats(data.UserID)
	err = pageInit(len(data.Sites), 10, &data.Page, r)
	if err != nil {
		http.NotFound(w, r)
	}

	err = templates.Lookup("sites.html").Execute(w, data)
	if err != nil {
		http.NotFound(w, r)
	}
}

//Слепленный "на коленке" разделитель страниц
func pageInit(num int, size int, page *Pages, r *http.Request) error {
	p := r.URL.Query().Get("page")
	if p == "" {
		page.Current = 1
	} else {
		if n, e := strconv.Atoi(p); e == nil {
			page.Current = n
		} else {
			return errors.New("Bad page number")
		}
	}
	page.PagesNum = num/size + 1
	if page.Current > page.PagesNum {
		return errors.New("Empty page")
	}
	page.IndexStart = (page.Current - 1) * size
	page.IndexEnd = page.Current * size
	page.Prev = page.Current - 1
	page.Next = page.Current + 1
	return nil
}
