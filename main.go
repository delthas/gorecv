//go:generate packer index.html index.go index
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := flag.Int("port", 8080, "http server listen port")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Body != http.NoBody {
			name := r.URL.Query().Get("name")
			if i := strings.LastIndexAny(name, "/\\"); i >= 0 {
				name = name[i+1:]
			}
			if name == "" {
				r.Body.Close()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("no file name set!"))
				return
			}
			f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
			if err != nil {
				r.Body.Close()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("file already exists on target computer: " + err.Error()))
				return
			}
			_, err = io.Copy(f, r.Body)
			f.Close()
			r.Body.Close()
			if err != nil {
				os.Remove(f.Name())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("uploading file failed: " + err.Error()))
				return
			}
		} else {
			w.Write(index)
		}
	})
	fmt.Println("listening on http://localhost:" + strconv.Itoa(*port) + "/")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
