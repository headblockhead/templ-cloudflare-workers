package main

import (
	"net/http"
	"os"
	"strconv"

	"time"

	"github.com/a-h/templ"
	"github.com/headblockhead/templwasm/session"
	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
)

func main() {
	kv, err := cloudflare.NewKVNamespace("templ_counter")
	if err != nil {
		os.Stderr.WriteString("error creating KV namespace: " + err.Error())
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		global_count, err := get(kv, "global")
		if err != nil {
			http.Error(w, "error fetching global counter", http.StatusInternalServerError)
			return
		}
		session_count := 0
		if session.ID(req) != "" {
			session_count, err = get(kv, session.ID(req))
			if err != nil {
				http.Error(w, "error fetching session counter", http.StatusInternalServerError)
				return
			}
		}
		templ.Handler(page(global_count, session_count)).ServeHTTP(w, req)
	})
	mux.HandleFunc("/increment/global", func(w http.ResponseWriter, req *http.Request) {
		count, err := increment(kv, "global")
		if err != nil {
			http.Error(w, "error incrementing global counter", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(strconv.Itoa(count)))
	})
	mux.HandleFunc("/increment/session", func(w http.ResponseWriter, req *http.Request) {
		count, err := increment(kv, session.ID(req))
		if err != nil {
			http.Error(w, "error incrementing session counter", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(strconv.Itoa(count)))
	})
	mux.HandleFunc("/sse", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")
		for {
			globalCount, err := get(kv, "global")
			if err != nil {
				http.Error(w, "error fetching global counter", http.StatusInternalServerError)
				return
			}
			w.Write([]byte("event: global\ndata: " + strconv.Itoa(globalCount) + "\n\n"))

			sessionCount, err := get(kv, session.ID(req))
			if err != nil {
				http.Error(w, "error fetching session counter", http.StatusInternalServerError)
				return
			}
			w.Write([]byte("event: session\ndata: " + strconv.Itoa(sessionCount) + "\n\n"))

			time.Sleep(time.Millisecond * 500)
		}
	})
	withCookie := session.NewMiddleware(mux)
	workers.Serve(withCookie)
}

func get(kv *cloudflare.KVNamespace, key string) (count int, err error) {
	countStr, err := kv.GetString(key, nil)
	if err != nil {
		return count, err
	}
	count, _ = strconv.Atoi(countStr)
	return count, nil
}

func put(kv *cloudflare.KVNamespace, key string, count int) error {
	return kv.PutString(key, strconv.Itoa(count), nil)
}

func increment(kv *cloudflare.KVNamespace, key string) (count int, err error) {
	count, err = get(kv, key)
	if err != nil {
		return count, err
	}
	count++
	err = put(kv, key, count)
	if err != nil {
		return count, err
	}
	return count, nil
}
