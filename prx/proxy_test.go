package prx

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

var logger = log.New(ioutil.Discard, "log: ", log.Ltime) //"пустой" лог
var testProxy = &ProxyServ{Log: logger, Debug: false}
var bufferedImage *bufio.Reader

var testServ *httptest.Server

func init() {

	file, err := os.Open("test/Koala.jpg")
	if err != nil {
		log.Fatal(err)
	}
	bufferedImage, _ := ioutil.ReadAll(file)
	file.Close()

	testServ = httptest.NewServer(nil)
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("HELLO"))
	})
	http.HandleFunc("/someHtml", func(w http.ResponseWriter, req *http.Request) {
		file, err := os.Open("test/SomeHtml.html")
		if err != nil {
			log.Fatal(err)
		}
		rd := bufio.NewReader(file)
		w.WriteHeader(200)
		w.Header().Add("Content Type", "text/html")
		rd.WriteTo(w)
	})
	http.HandleFunc("/image", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Header().Add("Content Type", "image/jpg")
		w.Write(bufferedImage)
	})
}

//Простейший тест, на выходе должны получить ответ "200 OK HELLO"
func TestSimpleTextServer(t *testing.T) {
	proxy := httptest.NewServer(testProxy)
	defer proxy.Close()
	pURL, _ := url.Parse(proxy.URL)

	cl := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(pURL)}}
	req, _ := http.NewRequest("GET", testServ.URL+"/hello", nil)
	resp, err := cl.Do(req)
	if err != nil {
		t.Fatalf("Ошибка получения ответа от тествого сервера: %v", err)
	}
	st := resp.Status
	if st != "200 OK" {
		t.Errorf("%s: Неправильный статус (ожидалось 200 OK)", st)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Невозможно прочитать тело ответа: %v", err)
	}
	if string(body) != "HELLO" {
		t.Errorf("%v: Неправильное тело ответа (ожидалось HELLO)", body)
	}
	resp.Body.Close()
}

//Функция для получения ответа от сервера с через прокси и без него
func GetResponce(path string, t *testing.T) {
	proxy := httptest.NewServer(testProxy)
	defer proxy.Close()
	pURL, _ := url.Parse(proxy.URL)

	//Создаем 2 хттп-клиента, один пускаем через прокси-сервер
	clNonProxy := &http.Client{}
	clProxy := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(pURL)}}

	req, _ := http.NewRequest("GET", path, nil)

	resp1, _ := clNonProxy.Do(req)
	resp2, _ := clProxy.Do(req)

	//используем strings.ToLower т.к. регистр ответа зависит от сервера
	if strings.ToLower(resp1.Status) != "200 ok" || strings.ToLower(resp2.Status) != "200 ok" {
		t.Fatalf("%s && %s: Неправильные статусы ответа (ожидалось 200 OK)", resp1.Status, resp2.Status)
	}

	//Если используется тестовый сервер тестируем хедеры и тела ответа на идентичность
	if strings.HasPrefix(path, testServ.URL) {
		for k, v := range resp1.Header {
			for kk, vv := range v {
				if resp2.Header[k][kk] != vv {
					t.Errorf("Отличаются хедеры:\n %v (без прокси)\n и \n%v (с прокси)", resp1.Header[k][kk], resp2.Header[k][kk])
				}
			}
		}

		body1, _ := ioutil.ReadAll(resp1.Body)
		body2, _ := ioutil.ReadAll(resp2.Body)

		if string(body1) != string(body2) {
			t.Error("Тело ответа различается для разных запросов")
		}

		resp1.Body.Close()
		resp2.Body.Close()
	}
}

func TestHtmlResponces(t *testing.T)     { GetResponce(testServ.URL+"/someHtml", t) }
func TestImageResponces(t *testing.T)    { GetResponce(testServ.URL+"/image", t) }
func TestRealLifeResponces(t *testing.T) { GetResponce("http://yandex.ru", t) }

//#################################################################################
//###########################    BENCHMARKS  ######################################
//#################################################################################

func BenchmarkSimpleTextResp(b *testing.B) {
	proxy := httptest.NewServer(testProxy)
	defer proxy.Close()
	pURL, _ := url.Parse(proxy.URL)

	cl := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(pURL)}}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cl.Get(testServ.URL + "/hello")
	}

}

func BenchmarkSimpleTextRespWithoutProxy(b *testing.B) {

	cl := &http.Client{}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cl.Get(testServ.URL + "/hello")
	}

}

func BenchmarkImageResp(b *testing.B) {
	proxy := httptest.NewServer(testProxy)
	defer proxy.Close()
	pURL, _ := url.Parse(proxy.URL)

	cl := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(pURL)}}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cl.Get(testServ.URL + "/image")
	}

}

func BenchmarkImageRespWithoutProxy(b *testing.B) {

	cl := &http.Client{}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cl.Get(testServ.URL + "/image")
	}

}
