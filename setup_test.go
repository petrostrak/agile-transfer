package main

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestRoutes(t *testing.T) {
	var registered = []struct {
		route  string
		method string
	}{
		{"/accounts/", "GET"},
		{"/accounts/", "POST"},
		{"/accounts/{id}", "GET"},
		{"/accounts/{id}", "PATCH"},
		{"/accounts/{id}", "DELETE"},
		{"/transfer", "POST"},
		{"/transactions", "GET"},
	}

	mux := Routes()
	chiRoutes := mux.(chi.Routes)

	for _, route := range registered {
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})

	return found
}
