package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	fmt.Printf("Start service %s\n", os.Getenv("SERVICE_NAME"))
	http.HandleFunc("/service/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

// HelloServer ...
func HelloServer(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(
		w,
		"Hello from behind Envoy (service %s)! hostname: %s resolved hostname: %s\n",
		os.Getenv("SERVICE_NAME"),
		hostname,
		addrs,
	)
}
