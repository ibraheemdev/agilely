package router

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func TestRouter(t *testing.T) {
	t.Parallel()

	rt := NewRouter(httprouter.New())
	var get, post, del, head, options, put, patch string

	rt.GET("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		get = string(b)
	}))
	rt.POST("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		post = string(b)
	}))
	rt.DELETE("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		del = string(b)
	}))
	rt.HEAD("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		head = string(b)
	}))
	rt.OPTIONS("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		options = string(b)
	}))
	rt.PUT("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		put = string(b)
	}))
	rt.PATCH("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		patch = string(b)
	}))

	requests := [][]string{
		[]string{"PUT", "testput"},
		[]string{"GET", "testget"},
		[]string{"HEAD", "testhead"},
		[]string{"OPTIONS", "testoptions"},
		[]string{"POST", "testpost"},
		[]string{"DELETE", "testdelete"},
		[]string{"PATCH", "testpatch"},
	}

	writers := []*string{&put, &get, &head, &options, &post, &del, &patch}

	for k, r := range requests {
		wr := httptest.NewRecorder()
		req := httptest.NewRequest(r[0], "/test", strings.NewReader(r[1]))
		rt.ServeHTTP(wr, req)
		if *writers[k] != r[1] {
			t.Error("want:", r[1], "got:", writers[k])
		}
	}
}

func TestRouterMiddleware(t *testing.T) {
	// TODO : Test all request types for middlewares
	t.Parallel()

	r := NewRouter(httprouter.New())
	r.Use(alice.New(helloWorldMiddleware))
	r.GET("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, ok := r.Context().Value("called").(string)
		if val != "yup" || !ok {
			t.Error("expected middleware to be called")
		}
	}))
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
}

func helloWorldMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "called", "yup")
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
