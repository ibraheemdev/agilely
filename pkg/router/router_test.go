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

	requests := [][]interface{}{
		[]interface{}{"PUT", "testput", &put, rt.PUT},
		[]interface{}{"GET", "testget", &get, rt.GET},
		[]interface{}{"HEAD", "testhead", &head, rt.HEAD},
		[]interface{}{"OPTIONS", "testoptions", &options, rt.OPTIONS},
		[]interface{}{"POST", "testpost", &post, rt.POST},
		[]interface{}{"DELETE", "testdelete", &del, rt.DELETE},
		[]interface{}{"PATCH", "testpatch", &patch, rt.PATCH},
	}

	for _, r := range requests {
		r[3].(func(path string, handler http.Handler))("/test", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}

			*r[2].(*string) = string(b)
		}))

		wr := httptest.NewRecorder()
		req := httptest.NewRequest(r[0].(string), "/test", strings.NewReader(r[1].(string)))
		rt.ServeHTTP(wr, req)
		if *r[2].(*string) != r[1].(string) {
			t.Error("want:", r[1], "got:", *r[2].(*string))
		}
	}
}

func TestRouterMiddleware(t *testing.T) {
	t.Parallel()

	r := NewRouter(httprouter.New())
	r.Use(alice.New(helloWorldMiddleware))
	r.GET("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val := r.Context().Value("called").(string); val != "yup" {
			t.Error("expected middleware to be called")
		}
	}))
}

func helloWorldMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "called", "yup")
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
