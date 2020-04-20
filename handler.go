package main

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type handler struct {
	pattern *regexp.Regexp
	db      *bolt.DB
}

func NewHandler(db *bolt.DB) http.Handler {
	return &handler{
		db:      db,
		pattern: regexp.MustCompile("^/[A-Za-z0-9-_]+/[A-Za-z0-9-_]+/visits.svg$"),
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodGet:
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	switch {
	case h.pattern.MatchString(r.URL.Path):
		h.handleVisit(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *handler) handleVisit(w http.ResponseWriter, r *http.Request) {
	var count int64
	key := strings.TrimSuffix(r.URL.Path[1:], "/visits.svg")

	err := h.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("visits"))
		bcount := b.Get([]byte(key))
		if bcount == nil {
			bcount = []byte("0")
		}

		count, _ = strconv.ParseInt(string(bcount), 10, 64)

		if !strings.HasPrefix(r.Referer(), fmt.Sprintf("https://github.com/%s", key)) {
			return nil
		}

		count++

		err := b.Put([]byte(key), []byte(strconv.FormatInt(count, 10)))
		if err != nil {
			return errors.Wrap(err, "error on put in transaction")
		}

		return nil
	})
	if err != nil {
		log.Printf("%+v", errors.Wrap(err, "error on update transaction"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")

	canvas := svg.New(w)
	canvas.Start(70, 20)
	canvas.Rect(0, 0, 40, 20, `fill="#555"`)
	canvas.Rect(40, 0, 30, 20, `fill="#08C"`)
	canvas.Text(21, 15, "visits", `fill="#000"`, `font-family="Verdana,DejaVu Sans,sans-serif"`, `font-size="12"`, `text-anchor="middle"`, `opacity="0.1"`)
	canvas.Text(20, 14, "visits", `fill="#fff"`, `font-family="Verdana,DejaVu Sans,sans-serif"`, `font-size="12"`, `text-anchor="middle"`)
	canvas.Text(56, 15, number(count), `fill="#000"`, `font-family="Verdana,DejaVu Sans,sans-serif"`, `font-size="12"`, `text-anchor="middle"`, `opacity="0.1"`)
	canvas.Text(55, 14, number(count), `fill="#fff"`, `font-family="Verdana,DejaVu Sans,sans-serif"`, `font-size="12"`, `text-anchor="middle"`)
	canvas.End()
}

func number(i int64) string {
	switch {
	case i == 0:
		return "none"
	case i < 1000:
		return strconv.FormatInt(i, 10)
	case i < 10000:
		return strconv.FormatFloat(float64(i)/float64(1000), 'f', 1, 64) + "K"
	case i < 1000000:
		return strconv.FormatInt(i/1000, 10) + "K"
	case i < 10000000:
		return strconv.FormatFloat(float64(i)/float64(1000000), 'f', 1, 64) + "M"
	case i < 1000000000:
		return strconv.FormatInt(i/1000000, 10) + "M"
	default:
		return "∞"
	}
}
