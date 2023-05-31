package utils

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func Test_LogError(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	LogError(ErrInsufficientBalance)
	w.Close()
	os.Stdout = stdout

	result, _ := io.ReadAll(r)
	expected := fmt.Sprintf("[ERROR] %v %v\n", time.Now().Format("2006/01/02 15:04:05"), ErrInsufficientBalance)

	if string(result) != expected {
		t.Errorf("incorrect log error: expected %s but got %s", expected, string(result))
	}
}
