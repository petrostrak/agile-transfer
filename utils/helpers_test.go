package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
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
