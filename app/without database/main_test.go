package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	uh := newsUserhandler()
	request, _ := http.NewRequest("GET", "/users", nil)
	response := httptest.NewRecorder()
	uh.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestUserswithid(t *testing.T) {
	uh := newsUserhandler()
	request, _ := http.NewRequest("GET", "/users/", nil)
	response := httptest.NewRecorder()
	uh.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestPosts(t *testing.T) {
	ph := newspostHandler()
	request, _ := http.NewRequest("GET", "/posts", nil)
	response := httptest.NewRecorder()
	ph.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestPost(t *testing.T) {
	ph := newspostHandler()
	http.Handle("/posts/", ph)
	request, _ := http.NewRequest("GET", "/posts/", nil)
	response := httptest.NewRecorder()
	ph.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestPostsperuser(t *testing.T) {
	ph := newspostHandler()
	http.Handle("/posts/users/", ph)
	request, _ := http.NewRequest("GET", "/posts/users/", nil)
	response := httptest.NewRecorder()
	ph.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

}
