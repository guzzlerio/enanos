package main

import (
	"fmt"
	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_Enanos(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Enanos Server:", func() {
		url := func(path string) (fullPath string) {
			fullPath = "http://localhost:8000" + path
			return
		}

		g.BeforeEach(func() {
			go func() {
				defaultHappy := func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}
				defaultGrumpy := func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
				mux := http.NewServeMux()
				mux.HandleFunc("/default/happy", defaultHappy)
				mux.HandleFunc("/default/grumpy", defaultGrumpy)
				err := http.ListenAndServe(":8000", mux)
				if err != nil {
					fmt.Errorf("error encountered %v", err)
				}
			}()
		})
		g.Describe("Happy :", func() {
			g.It("GET returns 200", func() {
				resp, _ := http.Get(url("/default/happy"))
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})
		})

		g.Describe("Grumpy :", func() {
			g.It("GET returns 500", func() {
				resp, _ := http.Get(url("/default/grumpy"))
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})
		})
	})
}
