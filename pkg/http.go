package cheek

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"sort"
	"strings"
)

//go:embed public
var files embed.FS

func fsys() fs.FS {
	fsys, err := fs.Sub(files, "public")
	if err != nil {
		log.Fatal(err)
	}

	return fsys
}

func server(s *Schedule) {
	var httpAddr string = fmt.Sprintf(":%s", s.cfg.Port)
	type Healthz struct {
		Jobs   int    `json:"jobs"`
		Status string `json:"status"`
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		status := Healthz{Jobs: len(s.Jobs), Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/", ui(s))

	fs := http.FileServer(http.FS(fsys()))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	s.log.Info().Msgf("Starting HTTP server on %v", httpAddr)
	s.log.Fatal().Err(http.ListenAndServe(httpAddr, nil))
}

func ui(s *Schedule) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get jobid from url
		var jobId string
		var job *JobSpec
		var ok bool
		job = &JobSpec{}

		if strings.HasPrefix(r.URL.Path, "/job/") {
			jobId = strings.TrimPrefix(r.URL.Path, "/job/")
			job, ok = s.Jobs[jobId]
			if !ok {
				jobId = ""
			} else {
				job.loadRuns()
			}
		}

		// get job ids
		jobNames := make([]string, 0)
		for k := range s.Jobs {
			jobNames = append(jobNames, k)
		}
		sort.Strings(jobNames)

		// parse template files
		tmpl, err := template.ParseFS(fsys(), "*.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data := struct {
			SelectedJob string
			JobNames    []string
			JobSpec     JobSpec
		}{SelectedJob: jobId, JobNames: jobNames, JobSpec: *job}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
