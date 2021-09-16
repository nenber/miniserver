package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func notAllowed(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, req.Method+" is not allowed.")
}

func badRequest(w http.ResponseWriter, req *http.Request, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, message)
}

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		notAllowed(w, req)
	} else {
		fmt.Fprintf(w, time.Now().Format("15:04"))
	}
}

func add(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	if req.Method != "POST" {
		notAllowed(w, req)
	} else {
		author := req.Form.Get("author")
		message := req.Form.Get("message")

		if len(author) > 0 && len(message) > 0 {
			f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				log.Fatal(err)
			}

			defer f.Close()

			_, err2 := f.WriteString(author + ":" + message + "\n")

			if err2 != nil {
				log.Fatal(err2)
			}

			fmt.Fprintf(w, author+" say "+message)
		} else {
			badRequest(w, req, "Missing parameters")
		}
	}
}

func entries(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		notAllowed(w, req)
	} else {
		raw, err := os.ReadFile("data.txt")

		if err != nil {
			panic(err)
		}
		fmt.Println(raw)
		data := strings.Split(string(raw), "\n")

		for _, rawEntry := range data {
			entry := strings.Split(rawEntry, ":")

			fmt.Fprintf(w, entry[1]+"\n")
		}
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/add", add)
	http.HandleFunc("/entries", entries)

	fmt.Println("Server started on http://localhost:4567")
	http.ListenAndServe(":4567", nil)
}
