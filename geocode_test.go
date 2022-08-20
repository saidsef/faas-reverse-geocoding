package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestLatitudeLongitude(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	latitudeLongitude(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if req == nil {
		t.Errorf("the request object is empty %v", req)
	}
	if res.Body == nil {
		t.Errorf("expected body to have value, but got %d", res.StatusCode)
	}
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) == "" {
		t.Errorf("expected status got %v", string(data))
	}
}

func TestMain(t *testing.T) {
	tests := []struct {
		port string
	}{
		{"8080"},
	}
	for _, tt := range tests {
		t.Run(tt.port, func(t *testing.T) {
			if i, err := strconv.ParseInt(port, 10, 64); err == nil {
				t.Errorf("found error on port number: %d and %v", i, err)
			}
		})
	}
}

func Test_latitudeLongitude(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test get /",
			args: args{w: rr, r: req},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			latitudeLongitude(tt.args.w, tt.args.r)
		})
	}
}
