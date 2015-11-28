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

	createRequest := func(path, method string) *http.Request {
		req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", server.URL, path), nil)
		Expect(err).NotTo(HaveOccurred())
		return req
	}

	BeforeEach(func() {
		router := mux.NewRouter()

		router.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/text", http.StatusFound)
		}).Methods("POST")

		router.HandleFunc("/text", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Pork and beans"))
		}).Methods("GET")

		server = httptest.NewServer(router)
	})

	JustBeforeEach(func() {
		client = herottp.New(config)
		resp, respErr = client.Do(req)
	})

	AfterEach(func() {
		server.Close()
	})

	Context("default configuration", func() {
		BeforeEach(func() {
			config = herottp.Config{}
			req = createRequest("redirect", "POST")
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

	Context("when following redirects is disabled", func() {
		BeforeEach(func() {
			config = herottp.Config{NoFollowRedirect: true}
			req = createRequest("redirect", "POST")
		})

		It("returns the redirect response", func() {
			Expect(respErr).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusFound))
		})
	})
})
