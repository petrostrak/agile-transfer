package utils

import (
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
