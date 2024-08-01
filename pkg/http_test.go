package cheek

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestMux(t *testing.T) {

	s1 := Schedule{
		Jobs:       map[string]*JobSpec{},
		TZLocation: "Europe/Amsterdam",
		log:        zerolog.Logger{},
		cfg:        NewConfig(),
	}

	s2 := Schedule{
		Jobs:       map[string]*JobSpec{},
		TZLocation: "Europe/Amsterdam",
		log:        zerolog.Logger{},
		cfg:        NewConfig(),
	}

	j := &JobSpec{
		Cron:    "MooIAmACow",
		Name:    "bertha",
		Command: []string{"ls"},
		cfg:     NewConfig(),
	}
	s2.Jobs[j.Name] = j

	type args struct {
		req *http.Request
	}

	tests := []struct {
		schedule *Schedule
		name     string
		args     func(t *testing.T) args
		wantCode int
		wantBody string
	}{
		{
			schedule: &s2,
			name:     "/jobs/bertha/runs/1 must return xxx",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/api/jobs/bertha/runs/2", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusNotFound,
			wantBody: "error: can't find job",
		},
		{
			schedule: &s2,
			name:     "/jobs/bertha must return 200",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/api/jobs/bertha", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusOK,
			wantBody: "MooIAmACow",
		},
		{
			schedule: &s2,
			name:     "/jobs/bertha/1 must return 404",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/api/jobs/bertha/1", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusNotFound,
			wantBody: "not found",
		},

		{
			schedule: &s1,
			name:     "/healthz must return 200",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/healthz/", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusOK,
			wantBody: "ok",
		},
		{
			schedule: &s1,
			name:     "/api/jobs/ must return 200",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/api/jobs", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusOK,
			wantBody: "{}",
		},
		{
			schedule: &s1,
			name:     "trigger must return 401",
			args: func(*testing.T) args {
				req, err := http.NewRequest("POST", "/api/jobs/does_not_exist/trigger", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusNotFound,
			wantBody: "error:",
		},
		{
			schedule: &s2,
			name:     "/trigger/ must return 200",
			args: func(*testing.T) args {
				req, err := http.NewRequest("POST", "/api/jobs/bertha/trigger", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusOK,
			wantBody: "\"status\":\"ok\",\"type\":\"trigger\"",
		},
		{
			schedule: &s1,
			name:     "/ must return 200 with html content",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusOK,
			wantBody: "href=\"/\">cheek</a>",
		},
		{
			schedule: &s1,
			name:     "/jobs/not-exist must return 404",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/job/not-exist", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusNotFound,
			wantBody: "",
		},
		{
			schedule: &s1,
			name:     "/api/schedule/status must return 200",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/api/schedule/status", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusOK,
			wantBody: "{}",
		},
	}

	for _, tt := range tests {
		handler := setupRouter(tt.schedule)

		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, tArgs.req)

			if resp.Result().StatusCode != tt.wantCode {
				t.Fatalf("the status code should be [%d] but received [%d]", tt.wantCode, resp.Result().StatusCode)
			}

			if !strings.Contains(resp.Body.String(), tt.wantBody) {
				t.Fatalf("the response body should contain [%s] but received [%s]", tt.wantBody, resp.Body.String())
			}

		})
	}
}
