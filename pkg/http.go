package cheek

import (
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
	router.GET("/healthz/", getHealthCheck)
	router.GET("/jobs/:jobId/:jobRunId", getJobDetailPage(s))
	router.GET("/", getHomePage())

	// api endpoints
	router.GET("/api/jobs", getJobs(s))
	router.GET("/api/jobs/:jobId", getJob(s))
	router.GET("/api/jobs/:jobId/runs/:jobRunId", getJobRun(s))
	router.POST("/api/jobs/:jobId/trigger", postTrigger(s))

	fileServer := http.FileServer(http.FS(fsys()))
	router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fileServer.ServeHTTP(w, r)
	})

	return router
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
			j.loadRunsFromDb(10, false)
		}

		if err := json.NewEncoder(w).Encode(s.Jobs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getJob(s *Schedule) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		jobId := ps.ByName("jobId")
		job, ok := s.Jobs[jobId]

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

		if !ok {
			status := Response{Job: jobId, Status: "error: can't find job to get runs", Type: "runs"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
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

		job.execCommandWithRetry("ui") // trigger

		status := Response{Job: jobId, Status: "ok", Type: "trigger"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// endpoints needed
// /jobs
// /jobs/{jobid}
// /jobs/{jobid}/logs/{logid}

// func ui(s *Schedule) func(w http.ResponseWriter, r *http.Request) {

// 	return func(w http.ResponseWriter, r *http.Request) {

// 		// get jobid from url
// 		var jobId string
// 		var job *JobSpec
// 		var ok bool
// 		job = &JobSpec{}

// 		if strings.HasPrefix(r.URL.Path, "/job/") {
// 			jobId = strings.TrimPrefix(r.URL.Path, "/job/")
// 			job, ok = s.Jobs[jobId]
// 			if !ok {
// 				http.Error(w, fmt.Errorf("job %s not found", jobId).Error(), http.StatusNotFound)
// 				return
// 			} else {
// 				job.loadRunsFromDb(50, false)
// 			}
// 		}

// 		// get job ids
// 		jobNames := make([]string, 0)
// 		for k := range s.Jobs {
// 			jobNames = append(jobNames, k)
// 		}
// 		sort.Strings(jobNames)

// 		// add custom functions to template
// 		funcMap := template.FuncMap{
// 			"roundToSeconds": func(d time.Duration) int {
// 				return int(d.Seconds())
// 			},
// 		}

// 		// parse template files
// 		tmpl, err := template.New("index.html").Funcs(funcMap).ParseFS(fsys(), "*.html")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		data := struct {
// 			SelectedJobName string
// 			JobNames        []string
// 			JobSpecs        map[string]*JobSpec
// 			SelectedJobSpec JobSpec
// 		}{SelectedJobName: jobId, JobNames: jobNames, SelectedJobSpec: *job}

// 		if jobId == "" {
// 			for _, j := range s.Jobs {
// 				j.loadRunsFromDb(5, false)
// 			}
// 			data.JobSpecs = s.Jobs
// 		}

// 		err = tmpl.Execute(w, data)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }
