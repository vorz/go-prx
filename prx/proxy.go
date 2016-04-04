package prx

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var tr = &http.Transport{
	DisableCompression: true,
	DisableKeepAlives:  true, //TODO need testing
}

//счетчик TCP-соединений
var numCon int

//ProxyServ структура реализующая интерфейс Handler
type ProxyServ struct {
	Log   *log.Logger
	Debug bool
	sync.Mutex
}

var users = make(map[string]string)

func (p *ProxyServ) manageUsers(r *http.Request) string {
	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if _, ok := users[clientIP]; !ok {
			users[clientIP] = "unknown"
			go func(ip string) {
				name, err := net.LookupAddr(ip)
				if err == nil && len(name) > 0 {
					users[ip] = name[0]
				}
			}(clientIP)
		}
		return clientIP
	}
	return ""
}

//Стандартная функция ServeHTTP интерфейса Handler
func (p *ProxyServ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Debugf("-------Получен запрос %v %v %v %v", r.URL.Path, r.Host, r.Method, r.URL.String())

	p.Lock()
	numCon++
	p.Unlock()

	if r.Method == "CONNECT" {
		p.handleConnect(w, r)
	} else {
		p.handleHTTP(w, r)
	}

}

//Функция обработки обычных http запросов
func (p *ProxyServ) handleHTTP(w http.ResponseWriter, r *http.Request) {
	p.Debugf("handleHttp. Запущена обработка http-соединения с методом %v", r.Method)

	start := time.Now()

	/*
		for key, val := range r.Header {
			p.Log.Print(key, " : ", val)
		}
		p.Log.Print(r.Body)

		p.Log.Print("Host: ", r.Host, " ; ", " Is ABS: ", r.URL.IsAbs())
		p.Log.Print("Remote: ", r.RemoteAddr, " ; ", " URL RequestURI: ", r.URL.RequestURI())
		p.Log.Print("Content Length: ", r.ContentLength)
	*/

	r.RequestURI = ""
	r.Header.Del("Accept-Encoding")
	//if _, ok := r.Header["Proxy-Connection"]; ok {
	//	r.Header.Set("Connection", r.Header["Proxy-Connection"][0])
	//}
	r.Header.Del("Proxy-Authenticate")
	r.Header.Del("Proxy-Authorization")
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Connection")

	resp, err := tr.RoundTrip(r)
	if err != nil {
		p.Warnf("ОШИБКА при пересылке: %s", err.Error())
		http.Error(w, "Ошибка доступа к удаленному сайту: "+err.Error(), 404)
		return
	}
	if resp == nil {
		p.Warnf("ОШИБКА получения ответа от сервера %v %v:", r.URL.Host, err.Error())
		http.Error(w, "ОШИБКА получения ответа от сервера: "+err.Error(), 500)
		return
	}
	p.Debugf("Получен ответ от %v: %v", r.URL.Host, resp.Status)

	//если получен "трейлер" заголовок
	if len(resp.Trailer) > 0 {
		p.Warnf("ОШИБКА Trayler хедер от: %v %v:", r.URL.Host, err.Error())
	}

	for k := range w.Header() {
		w.Header().Del(k)
	}

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.Warnf("ОШИБКА чтения тела ответа %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := resp.Body.Close(); err != nil {
		p.Warnf("ОШИБКА: невозможно закрыть http.Body. %v", err)
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	finish := time.Now()
	duration := finish.Sub(start)

	p.Lock()
	numCon--
	p.Unlock()

	ip := p.manageUsers(r)

	p.Warnf("[%d][%s:%s] %v: %v. %d байт за %v", numCon, ip, users[ip], r.Method, r.URL.Host, len(body), duration)

}

//Функция туннелирования CONNECT-запросов (https)
func (p *ProxyServ) handleConnect(w http.ResponseWriter, r *http.Request) {
	p.Debugf("=====handleConnect. Запущена обработка CONNECT")

	start := time.Now()

	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		p.Panicf("Proxy hijacking not supporting %s: %s", r.RemoteAddr, err.Error())
		return
	}

	_, err = io.WriteString(conn, "HTTP/1.0 200 Connection established\r\n\r\n")
	if err != nil {
		p.Warnf("ОШИБКА: Невозможно отправить ответ %s: %s", r.RemoteAddr, err.Error())
		return
	}

	dstConn, err := net.Dial("tcp", r.URL.Host)
	if err != nil {
		p.Warnf("ОШИБКА: Невозможно соединиться с %s: %s", r.RequestURI, err.Error())
		return
	}

	var done = make(chan int64)

	//Функция копирования и закрытия соединения
	fCopy := func(from, to net.Conn) {
		n, err := io.Copy(to, from)
		if err != nil {
			p.Warnf("ОШИБКА копирования %s", err.Error())
		}
		done <- n
	}

	go fCopy(dstConn, conn)
	go fCopy(conn, dstConn)

	num := <-done
	num += <-done

	if err := conn.Close(); err != nil {
		p.Warnf("ОШИБКА закрытия оригинального соединения %s", err.Error())
	}
	if err := dstConn.Close(); err != nil {
		p.Warnf("ОШИБКА закрытия соединения с удаленным узлом %s", err.Error())
	}

	finish := time.Now()
	duration := finish.Sub(start)

	p.Lock()
	numCon--
	p.Unlock()

	ip := p.manageUsers(r)

	p.Warnf("[%d][%s:%s] %v:  %v. %d байт за %v", numCon, ip, users[ip], r.Method, r.URL.Host, num, duration)
}

//Debugf вывод для отладки если задан параметр Debug
func (p *ProxyServ) Debugf(msg string, argv ...interface{}) {
	if p.Debug {
		p.Log.Printf(msg, argv...)
	}
}

//Warnf вывод при ошибках
func (p *ProxyServ) Warnf(msg string, argv ...interface{}) {
	p.Log.Printf(msg, argv...)
}

//Panicf вывод при критических ошибках (завершение программы)
func (p *ProxyServ) Panicf(msg string, argv ...interface{}) {
	p.Log.Fatalf("КРИТИЧЕСКАЯ ОШИБКА:"+msg, argv...)
}
