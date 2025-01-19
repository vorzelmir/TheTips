package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const (
	bash = "/usr/bin/bash"
	arg  = "-c"
	cmd  = "/usr/bin/cat /etc/os-release | /usr/bin/grep --ignore pretty_name"
	name = "PRETTY_NAME"
	os   = "Operational System "
)

type CustomHandler struct {
	message string
}

func (ch CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, ch.message)
}

// get current date/time
func getDate() (handler CustomHandler) {
	date := time.Now()
	handler.message = date.Format(time.RFC850)
	return
}

// get Operational System name of the backend server
func getInfo() (handler CustomHandler) {
	out, err := exec.Command(bash, arg, cmd).Output()

	if err != nil {
		log.Fatal(err)
	}

	replaced := strings.Replace(string(out), name, os, 1)
	handler.message = replaced
	return
}

func main() {
	customHandlerTime := getDate()

	customHandlerInfo := getInfo()

	sm := http.NewServeMux()
	server := http.Server{
		//port number have to be needed when configuring proxy server
		Addr:              ":8888",
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		Handler:           sm,
	}

	sm.Handle("/time", &customHandlerTime)
	sm.Handle("/os", &customHandlerInfo)

	log.Fatal(server.ListenAndServe())
}
