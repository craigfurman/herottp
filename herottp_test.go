package herottp_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/craigfurman/herottp"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HeroTTP", func() {
	var (
		client *herottp.Client
		config herottp.Config

		req *http.Request

		resp    *http.Response
		respErr error

		server *httptest.Server
	)

	BeforeEach(func() {
		router := mux.NewRouter()

		router.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/text", http.StatusFound)
		}).Methods("GET")

		router.HandleFunc("/text", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Pork and beans"))
		}).Methods("GET")

		server = httptest.NewServer(router)
	})

	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func() {
		client = herottp.New(config)
		resp, respErr = client.Do(req)
	})

	Describe("default configuration", func() {
		BeforeEach(func() {
			config = herottp.Config{}
			var err error
			req, err = http.NewRequest("GET", fmt.Sprintf("%s/redirect", server.URL), nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("follows redirects", func() {
			Expect(respErr).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal([]byte("Pork and beans")))
		})
	})
})
