package cheek

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"sort"
	"strings"
)

type Response struct {
	Job    string `json:"jobs,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

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

	http.HandleFunc("/healthz/", func(w http.ResponseWriter, r *http.Request) {
		status := Response{Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/schedule/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/trigger/", trigger(s))
	http.HandleFunc("/", ui(s))

	fs := http.FileServer(http.FS(fsys()))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	s.log.Info().Msgf("Starting HTTP server on %v", httpAddr)
	s.log.Fatal().Err(http.ListenAndServe(httpAddr, nil))
}

func trigger(s *Schedule) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		jobId := strings.TrimPrefix(r.URL.Path, "/trigger/")
		job, ok := s.Jobs[jobId]

		if !ok {
			http.Error(w, errors.New("cant find job to trigger").Error(), http.StatusNotFound)
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
			return
		}

		data := struct {
			SelectedJobName string
			JobNames        []string
			JobSpecs        map[string]*JobSpec
			SelectedJobSpec JobSpec
		}{SelectedJobName: jobId, JobNames: jobNames, SelectedJobSpec: *job}

		if jobId == "" {
			// pass along all job specs only when in overview
			// takes a lot of I/O
			for _, j := range s.Jobs {
				j.loadRuns()
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
