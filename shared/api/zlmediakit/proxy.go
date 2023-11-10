package zlmediakitapi

//将部分接口代理出去
import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func (f *MediaConfig) ProxyInit() {
	addr, err := url.Parse(f.PreUrl)

	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(addr)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Receive a reuest:%s %s %s ", r.Method, r.Host, r.URL)
		proxy.ServeHTTP(w, r)
	})

	log.Println("Listening on %s", f.PreUrl)
	if err := http.ListenAndServe(strconv.FormatInt(f.Port, 10), nil); err != nil {
		log.Fatal(err)
	}
}
