package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"

	//https://golang.org/pkg/net/http/pprof/
	//http://localhost:8082/debug/pprof/goroutine?debug=1
	_ "net/http/pprof"

	"os"
	"runtime"
	"text/template"

	"github.com/asalih/ca/ca"
)

var Manager *ca.Manager

func main() {
	startingGs := runtime.NumGoroutine()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		t, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Fprintf(w, "Unable to load template")
		}

		t.Execute(w, nil)
	})

	http.HandleFunc("/vulns", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if Manager == nil {
			json.NewEncoder(w).Encode(nil)
			return
		}

		json.NewEncoder(w).Encode(Manager.Attacker.FoundVulnerabilities)
	})

	http.HandleFunc("/attack", func(w http.ResponseWriter, r *http.Request) {

		if Manager != nil && !Manager.IsFinished() {
			w.Write([]byte("Attack in progress"))

			return
		}

		Manager = ca.NewManager(r.URL.Query().Get("attackUrl"))
		go Manager.CrawlAndAttack()

		w.Write([]byte("Attack started!"))
	})

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {

		if Manager != nil && Manager.IsFinished() {
			w.Write([]byte("finished"))
			return
		}
		w.Write([]byte("na"))
	})

	go http.ListenAndServe("localhost:8082", nil)

	endingGs := runtime.NumGoroutine()

	fmt.Println("Number of goroutines before:", startingGs)
	fmt.Println("Number of goroutines after :", endingGs)

	fmt.Println("Waiting an input for exit: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())
}
