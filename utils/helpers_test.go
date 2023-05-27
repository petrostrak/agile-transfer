package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Test_ReadIDParams(t *testing.T) {
	testCases := []struct {
		ID       string
		expected uuid.UUID
	}{
		{"121f03cd-ce8c-447d-8747-fb8cb7aa3a52", uuid.MustParse("121f03cd-ce8c-447d-8747-fb8cb7aa3a52")},
		{"2f0141f8-f325-4b15-9973-e7b34852e298", uuid.MustParse("2f0141f8-f325-4b15-9973-e7b34852e298")},
		{"3659fbd7-7ba2-4151-95a2-b977ebf79307", uuid.MustParse("3659fbd7-7ba2-4151-95a2-b977ebf79307")},
	}

	for _, tt := range testCases {
		var req *http.Request
		req, _ = http.NewRequest("GET", "/", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", tt.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		handler.ServeHTTP(rr, req)

		result := ReadIDParam(req)

		if result != tt.expected {
			t.Errorf("Expected %v but got %v", tt.expected, result)
		}
	}
}

func Test_HumanDate(t *testing.T) {
	testCases := []struct {
		inputTime time.Time
		expected  string
	}{
		{time.Unix(1405544146, 10), "16 Jul 2014 at 23:55"},
		{time.Unix(1505543477, 10), "16 Sep 2017 at 09:31"},
		{time.Unix(1504583477, 10), "05 Sep 2017 at 06:51"},
	}

	for _, tt := range testCases {
		result := HumanDate(tt.inputTime)

		if result != tt.expected {
			t.Errorf("Expected %v but got %v", tt.expected, result)
		}
	}
}

func Test_WriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := make(map[string]any)
	payload["foo"] = false

	headers := make(http.Header)
	headers.Add("FOO", "BAR")
	err := WriteJSON(rr, http.StatusOK, payload, headers)
	if err != nil {
		t.Errorf("failed to write JSON: %v", err)
	}
}

func Test_ReadJSON(t *testing.T) {
	sampleJSON := map[string]interface{}{
		"foo": "bar",
	}
	body, _ := json.Marshal(sampleJSON)

	var decodedJSON struct {
		Foo string `json:"foo"`
	}

	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Log("Error", err)
	}

	rr := httptest.NewRecorder()
	defer req.Body.Close()

	err = ReadJSON(rr, req, &decodedJSON)
	if err != nil {
		t.Error("failed to decode json", err)
	}

	badJSON := `
		{
			"foo": "bar"
		}
		{
			"alpha": "beta"
		}`

	req, err = http.NewRequest("POST", "/", bytes.NewReader([]byte(badJSON)))
	if err != nil {
		t.Log("Error", err)
	}

	err = ReadJSON(rr, req, &decodedJSON)
	if err == nil {
		t.Error("did not get an error with bad json")
	}
}
