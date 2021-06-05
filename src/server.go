package main

import (
	"fmt"
	"html"
	"net/http"
	"time"

	expiremap "github.com/nursik/go-expire-map"
)

var expireMap = expiremap.New()

func main() {
	Init()
	http.HandleFunc("/check", checkCompromise)
	http.ListenAndServe(":8989", nil)
}

func checkCompromise(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		email := html.EscapeString(r.URL.Query().Get("email"))

		result, _ := expireMap.Get(email)

		var exist bool

		if result != nil {
			exist = bool(result.(bool))
			fmt.Println("use cached value")
		} else {
			fmt.Println("cached not found, fetch from db")
			exist = email != "" && IsExist(email)
			expireMap.Set(email, exist, time.Duration(30)*time.Second)
		}

		w.Header().Set("Content-Type", "text/plain")

		if exist {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OUCH!!!! Email '%s' is COMPROMISED!", email)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Email '%s' is not detected to be compromised.", email)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}
