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
	"time"
)

type Response struct {
	Job    string `json:"jobs,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

//go:embed web_assets
var files embed.FS

func fsys() fs.FS {
	fsys, err := fs.Sub(files, "web_assets")
	if err != nil {
		log.Fatal(err)
	}

	return fsys
}

func setupMux(s *Schedule) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz/", func(w http.ResponseWriter, r *http.Request) {
		status := Response{Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/schedule/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/trigger/", trigger(s))
	mux.HandleFunc("/", ui(s))

	fs := http.FileServer(http.FS(fsys()))
	mux.Handle("/static/", fs)

	return mux

}

func server(s *Schedule) {
	var httpAddr string = fmt.Sprintf(":%s", s.cfg.Port)

	mux := setupMux(s)

	s.log.Info().Msgf("Starting HTTP server on %v", httpAddr)
	s.log.Fatal().Err(http.ListenAndServe(httpAddr, mux))

}

func trigger(s *Schedule) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		jobId := strings.TrimPrefix(r.URL.Path, "/trigger/")
		job, ok := s.Jobs[jobId]

		if !ok {
			status := Response{Job: jobId, Status: "error: can't find job to trigger", Type: "trigger"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		job.execCommandWithRetry("ui") // trigger

		status := Response{Job: jobId, Status: "ok", Type: "trigger"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
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
				http.Error(w, fmt.Errorf("job %s not found", jobId).Error(), http.StatusNotFound)
				return
			} else {
				job.loadRunsFromDb(50)
			}
		}

		// get job ids
		jobNames := make([]string, 0)
		for k := range s.Jobs {
			jobNames = append(jobNames, k)
		}
		sort.Strings(jobNames)

		// add custom functions to template
		funcMap := template.FuncMap{
			"roundToSeconds": func(d time.Duration) int {
				return int(d.Seconds())
			},
		}

		// parse template files
		tmpl, err := template.New("index.html").Funcs(funcMap).ParseFS(fsys(), "*.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			SelectedJobName string
			JobNames        []string
			JobSpecs        map[string]*JobSpec
			SelectedJobSpec JobSpec
		}{SelectedJobName: jobId, JobNames: jobNames, SelectedJobSpec: *job}

		if jobId == "" {
			for _, j := range s.Jobs {
				j.loadRunsFromDb(5)
			}
			data.JobSpecs = s.Jobs
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
