package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type TestCoordinates struct {
	Lat  string
	Long string
}

func isCoordinates(t interface{}) bool {
	switch t.(type) {
	case TestCoordinates:
		return true
	default:
		return false
	}
}

func TestLatitudeLongitude(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	latitudeLongitude(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) == "" {
		t.Errorf("expected status got %v", string(data))
	}
}

func TestMain(m *testing.M) {

	t := TestCoordinates{"41.40338", "2.17403"}
	fmt.Println(isCoordinates(t))

	exitEval := m.Run()
	os.Exit(exitEval)
}
