package api

import (
	"fmt"
	"net/http"
	"os/exec"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	out, err := exec.Command("ping", "-c", "1", host).Output()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out)
}

func UptimeHandler(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("uptime").Output()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out)
}

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "<h1>Hello %s</h1>", name)
}
