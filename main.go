package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRes struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type LoginRoundTripper struct {
	Tripper http.RoundTripper
}

type ArticleRoundTripper struct {
	Tripper http.RoundTripper
}

var (
	endpoint string          = "https://apingweb.com/api"
	email    string          = "superman@gmail.com"
	password string          = "123456"
	response UserLoginRes    = UserLoginRes{}
	client   *http.Client    = &http.Client{}
	ctx      context.Context = context.Background()
)

func main() {
	authTripper := NewAuthRoundTripper()
	authTripper.Tripper = &ArticleRoundTripper{
		Tripper: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			IdleConnTimeout: time.Duration(time.Second * 5),
		},
	}

	client.Transport = authTripper

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, bytes.NewReader(nil)) // overwrite default request

	if err != nil {
		log.Fatalf("NEW REQUEST %s", err.Error())
	}

	res, err := authTripper.RoundTrip(req)
	if err != nil {
		log.Fatalf("CLIENT DO %s", err.Error())
	}
	defer res.Body.Close()

	HttpDump(res, req)

	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("IO READ %s", err.Error())
	}

	fmt.Println(string(resByte))
}

/**
#########################
# AUTH ROUND TRIPPER
#########################
*/

func NewAuthRoundTripper() *LoginRoundTripper {
	return &LoginRoundTripper{}
}

func (h *LoginRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	body := make(map[string]interface{})
	body["email"] = email
	body["password"] = password

	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		log.Fatal(err.Error())
	}

	loginReq, err := http.NewRequest(http.MethodPost, endpoint+"/login", bytes.NewReader(bodyBytes))
	loginReq.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatal(err.Error())
	}

	res, err := http.DefaultClient.Do(loginReq)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+response.Token)

	return h.Tripper.RoundTrip(req)
}

/**
#########################
# ARTICLE ROUND TRIPPER
#########################
*/

func NewArticleRoundTripper() *ArticleRoundTripper {
	return &ArticleRoundTripper{}
}

func (h *ArticleRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL, _ = url.Parse(endpoint + "/articles")
	return h.Tripper.RoundTrip(req)
}

/**
############################
# DUMP REQUEST & RESPONSE
############################
*/

func HttpDump(res *http.Response, req *http.Request) {
	dumpRes, err := httputil.DumpResponse(res, true)
	defer log.Printf("METHOD: %s - URL: %s - RESPONSE: %s", req.Method, req.URL, string(dumpRes))

	if err != nil {
		log.Fatal(err.Error())
	}

	dumpReq, err := httputil.DumpRequestOut(req, false)
	defer log.Printf("METHOD: %s - URL: %s - REQUEST: %s", req.Method, req.URL, string(dumpReq))

	if err != nil {
		log.Fatal(err.Error())
	}
}
