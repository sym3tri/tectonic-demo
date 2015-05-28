package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/sym3tri/tectonic-demo/server"
)

const version = "v2.0"

func main() {
	fs := flag.NewFlagSet("demo", flag.ExitOnError)
	listen := fs.String("listen", "0.0.0.0:8080", "address/port to listen on")
	message := fs.String("message", "", "a message to print")
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	tpls, err := template.ParseFiles("./static/index.html")
	if err != nil {
		fmt.Printf("index.html template not found: %v", err)
		os.Exit(1)
	}

	cfg := server.Config{
		Message:   *message,
		Templates: tpls,
		Version:   version,
	}

	srv := &server.Server{
		Config: cfg,
	}

	httpsrv := &http.Server{
		Addr:    *listen,
		Handler: srv.HTTPHandler(),
	}

	log.Printf("Binding to %s...", httpsrv.Addr)
	log.Printf("Version: %s", version)
	log.Fatal(httpsrv.ListenAndServe())
}
