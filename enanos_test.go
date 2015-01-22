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
		g.BeforeEach(func() {
			go func() {
				defaultHappy := func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}
				mux := http.NewServeMux()
				mux.HandleFunc("/default/happy", defaultHappy)
				err := http.ListenAndServe(":8000", mux)
				if err != nil {
					fmt.Errorf("error encountered %v", err)
				}
			}()
		})
		g.Describe("Happy :", func() {
			g.It("GET returns 200", func() {
				resp, _ := http.Get("http://localhost:8000/default/happy")
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})
		})
	})
}
