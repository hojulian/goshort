package goshort

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const addr = "127.0.0.1:8080"

type Body struct {
	Body string `json:"body"`
}

func run() {
	var mux = http.NewServeMux()
	var gs = New(addr, mux)
	var srv = &http.Server{
		Addr:         gs.Addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	gs.SetupRoutes()
	_ = srv.ListenAndServe()
}

func TestShorten(t *testing.T) {
	t.Run(
		"shortening and expanding",
		func(t *testing.T) {
			go run()
			// shorten
			var body = Body{
				Body: "http://www.apple.com",
			}
			req, err := json.Marshal(body)
			require.Nil(t, err)
			var dest = "http://" + addr + "/s"
			res, err := http.Post(dest, "application/json", bytes.NewReader(req))
			require.Nil(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			// expand
			b, err := ioutil.ReadAll(res.Body)
			require.Nil(t, err)
			var s = Body{}
			err = json.Unmarshal(b, &s)
			require.Nil(t, err)
			res, err = http.Get(s.Body)
			require.Nil(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
		},
	)
}
