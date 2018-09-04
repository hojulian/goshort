package goshort

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type GoShort struct {
	Addr  string
	Mux   *http.ServeMux
	mutex *sync.Mutex
	src   rand.Source
}

func New(addr string) *GoShort {
	return &GoShort{
		Addr:  addr,
		Mux:   http.NewServeMux(),
		mutex: &sync.Mutex{},
		src:   rand.NewSource(time.Now().UnixNano()),
	}
}

func (g *GoShort) SetupRoutes() {
	g.Mux.HandleFunc("/s", g.Shorten)
}

func (g *GoShort) addHandler(pattern string, handler http.Handler) {
	g.Mux.Handle(pattern, handler)
}

func (g *GoShort) Shorten(w http.ResponseWriter, r *http.Request) {
	var body = struct {
		Body string `json:"body"`
	}{}
	var req, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot read body"))
		return
	}
	// unmarshalJSON
	err = json.Unmarshal(req, &body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid body JSON"))
		return
	}
	// check url
	_, err = url.Parse(body.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid url"))
		return
	}
	// add handler
	var id = g.randomize()
	g.addHandler("/"+id, http.RedirectHandler(body.Body, http.StatusFound))
	// return shortened url
	body.Body = "http://" + g.Addr + "/" + id
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to marshal body"))
	}
}

func (g *GoShort) randomize() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	// thread safe generator
	var int63 = func() int64 {
		g.mutex.Lock()
		v := g.src.Int63()
		g.mutex.Unlock()
		return v
	}
	// length
	var n = 8
	var b = make([]byte, n)

	for i, cache, remain := n-1, int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
