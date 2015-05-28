package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	Message   string
	Templates *template.Template
	Version   string
}

type Server struct {
	Config Config
}

func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", s.rootHandler())
	mux.Handle("/version", s.versionHandler())
	mux.Handle("/prestop", s.prestopHandler())
	mux.Handle("/poststart", s.poststartHandler())
	mux.Handle("/mount/", s.mountHandler())
	mux.Handle("/environment", s.environmentHandler())
	mux.Handle("/static/", s.staticHandler())
	return http.Handler(mux)
}

func (s *Server) rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Config.Templates.ExecuteTemplate(w, "index.html", s.Config); err != nil {
			log.Printf("error parsing index.htmlm template: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) staticHandler() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
}

func (s *Server) versionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Version: %s\n", s.Config.Version)
	}
}

func (s *Server) poststartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("post-start")
		fmt.Fprint(w, "post-start")
	}
}

func (s *Server) environmentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("environment")
		fmt.Fprint(w, strings.Join(os.Environ(), "\n"))
	}
}

func (s *Server) prestopHandler() http.HandlerFunc {
	waitSecs := 5
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "initiated pre-stop\n")
		log.Println("initiated pre-stop")
		for i := waitSecs; i > 0; i-- {
			log.Printf("shutting down in: %ds", i)
			fmt.Fprintf(w, "shutting down in: %ds\n", i)
			time.Sleep(time.Second * 1)
		}
		log.Println("pre-stop complete")
		fmt.Fprint(w, "pre-stop complete")
	}
}

// mountHandler Prints the contents of the file specified in the 'file' query param.
func (s *Server) mountHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		filename := q.Get("file")
		if filename == "" {
			fmt.Fprintf(w, "No file param specified")
			return
		}

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "file not mounted: %s\n", filename)
			return
		}
		fmt.Fprintf(w, "file contents for: %s\n", filename)
		w.Write(b)
	}
}
