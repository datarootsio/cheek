package cheek

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v3"
)

type TemplateData struct {
	Name string
}
type Response struct {
	Job    string `json:"jobs,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

// This will be injected at build time
var (
	version   string
	commitSHA string
)

type VersionResponse struct {
	Version   string `json:"version"`
	CommitSHA string `json:"commit_sha"`
}

type ScheduleStatusResponse struct {
	Status         map[string]int `json:"status,omitempty"`
	FailedRunCount int            `json:"failed_run_count,omitempty"`
	HasFailedRuns  bool           `json:"has_failed_runs,omitempty"`
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

func setupRouter(s *Schedule) *httprouter.Router {
	router := httprouter.New()

	// ui endpoints
	router.GET("/jobs/:jobId/:jobRunId", getJobDetailPage(s))
	router.GET("/core/logs", getCoreLogsPage())
	router.GET("/", getHomePage())

	// api endpoints
	router.GET("/healthz/", getHealthCheck)
	router.GET("/api/jobs", getJobs(s))
	router.GET("/api/jobs/:jobId", getJob(s))
	router.GET("/api/jobs/:jobId/runs/:jobRunId", getJobRun(s))
	router.POST("/api/jobs/:jobId/trigger", postTrigger(s))
	router.GET("/api/core/logs", getCoreLogs(s))
	router.GET("/api/schedule/status", getScheduleStatus(s))
	router.GET("/api/version", getVersion) // Add version endpoint

	fileServer := http.FileServer(http.FS(fsys()))
	router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fileServer.ServeHTTP(w, r)
	})

	return router
}

func getCoreLogsPage() httprouter.Handle {
	tmpl, err := template.ParseFS(fsys(), "templates/corelogs.html", "templates/base.html")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := tmpl.ExecuteTemplate(w, "base.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getHomePage() httprouter.Handle {

	tmpl, err := template.ParseFS(fsys(), "templates/overview.html", "templates/base.html")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := tmpl.ExecuteTemplate(w, "base.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getJobDetailPage(s *Schedule) httprouter.Handle {
	tmpl, err := template.ParseFS(fsys(), "templates/jobview.html", "templates/base.html")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		jobId := ps.ByName("jobId")
		_, ok := s.Jobs[jobId]
		if !ok {
			http.Error(w, fmt.Errorf("job %s not found", jobId).Error(), http.StatusNotFound)
			return
		}

		err := tmpl.ExecuteTemplate(w, "base.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func server(s *Schedule) {
	httpAddr := fmt.Sprintf(":%s", s.cfg.Port)
	router := setupRouter(s)

	s.log.Info().Msgf("Starting HTTP server on %v", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, router))
}

func getHealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := Response{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getJobs(s *Schedule) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		for _, j := range s.Jobs {
			j.loadRunsFromDb(20, false)
		}

		if err := json.NewEncoder(w).Encode(s.Jobs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getScheduleStatus(s *Schedule) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		ssr := ScheduleStatusResponse{
			Status: make(map[string]int, len(s.Jobs)),
		}

		for _, j := range s.Jobs {
			j.loadRunsFromDb(1, false)
			lastRunStatus := j.Runs[0].Status
			ssr.Status[j.Name] = *lastRunStatus
			if *lastRunStatus == 1 {
				ssr.FailedRunCount++
			}
		}

		ssr.HasFailedRuns = ssr.FailedRunCount > 0

		if err := json.NewEncoder(w).Encode(ssr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func getCoreLogs(s *Schedule) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		logs, err := getCoreLogsFromDB(s.cfg.DB, 120)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(logs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getJob(s *Schedule) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		jobId := ps.ByName("jobId")
		job, ok := s.Jobs[jobId]

		if !ok {
			status := Response{Job: jobId, Status: "error: can't find job to get runs", Type: "runs"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// convert job to YAML
		jobYaml, err := yaml.Marshal(job)
		if err != nil {
			status := Response{Job: jobId, Status: "error: can't convert job to yaml", Type: "yaml"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// get job runs from db
		job.loadRunsFromDb(50, false)

		job.Yaml = string(jobYaml)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(job); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getJobRun(s *Schedule) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//jobId is needed because db instance is on job level
		jobId := ps.ByName("jobId")
		runId := ps.ByName("jobRunId")

		runIdInt, err := strconv.Atoi(runId)
		job, ok := s.Jobs[jobId]

		if !ok || err != nil {
			status := Response{Job: jobId, Status: "error: can't find job / id to get runs", Type: "runs"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		jr, err := job.loadLogFromDb(runIdInt)
		if err != nil {
			status := Response{Job: jobId, Status: "error: can't find job / id to get runs", Type: "runs"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func postTrigger(s *Schedule) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		jobId := ps.ByName("jobId")
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

		job.execCommandWithRetry(context.Background(), "ui", nil) // trigger

		status := Response{Job: jobId, Status: "ok", Type: "trigger"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getVersion(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	versionResponse := VersionResponse{Version: version, CommitSHA: commitSHA}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versionResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
