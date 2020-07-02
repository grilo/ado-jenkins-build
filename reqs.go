package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

var (
	cookieJar, _ = cookiejar.New(nil)
	tr           = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: tr,
		Jar:       cookieJar,
	}
)

func Get(url string, params map[string]string) string {

	req, _ := http.NewRequest("GET", url, nil)

	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	log.Printf("GET: %s", req.URL.String())
	response, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	for _, cookie := range response.Cookies() {
		log.Printf("Found a cookie %s: %s", cookie.Name, cookie.Value)
	}

	defer response.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	bodyString := string(bodyBytes)

	if os.Getenv("TRACE") != "" {
		log.Printf(bodyString)
	}

	return bodyString
}

func Post(url string, headers map[string]string, params map[string]string) *http.Response {

	req, _ := http.NewRequest("POST", url, nil)

	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	for key, value := range headers {
		log.Printf("Adding Header %s: %s", key, value)
		req.Header.Add(key, value)
	}

	if os.Getenv("TRACE") != "" {
		dump, _ := httputil.DumpRequest(req, true)
		log.Println(string(dump))
	}

	log.Printf("POST: %s", req.URL.String())
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	for _, cookie := range response.Cookies() {
		log.Printf("Found a cookie %s: %s", cookie.Name, cookie.Value)
	}

	defer response.Body.Close()

	return response

}

func getJenkinsCrumb(jenkinsURL string) map[string]string {
	resp := Get(jenkinsURL+"/crumbIssuer/api/json", nil)
	json := gjson.Parse(resp)

	return map[string]string{
		json.Get("crumbRequestField").String(): json.Get("crumb").String(),
	}
}

func queryApiJson(url string) gjson.Result {
	buildResponse := Get(url+"api/json", nil)
	return gjson.Parse(buildResponse)
}

func TriggerBuild(url *url.URL, params map[string]string) string {

	var buildUrl string

	crumb := getJenkinsCrumb(url.Scheme + "://" + url.Host)

	log.Printf("Triggering build...")
	response := Post(url.String()+"/buildWithParameters", crumb, params)

	location := response.Header.Get("Location")
	if location == "" {
		log.Fatalf("Unable to obtain location after trigger build.")
	}
	log.Printf("Build queued: %s", location)

	for {
		time.Sleep(time.Second * 3)
		json := queryApiJson(location)
		buildUrl = json.Get("executable.url").String()
		if buildUrl != "" {
			break
		}

		if json.Get("cancelled").Bool() {
			log.Fatalf("Remote build cancelled: %s", location)
		}
	}

	log.Printf("Build url: %s", buildUrl)
	return buildUrl
}

func WaitForBuild(jobUrl string, timeout int) int {

	var buildResult string

	startTime := time.Now()

	for {
		if int(time.Now().Sub(startTime).Seconds()) > timeout {
			log.Fatalf("Timeout of (%d) seconds exceeded waiting for: %s", timeout, jobUrl)
		}
		time.Sleep(time.Second * 1)

		json := queryApiJson(jobUrl)
		buildResult = json.Get("result").String()
		if json.Get("building").Bool() != true {
			break
		}
	}

	log.Printf("State: %s", buildResult)

	switch buildResult { // https://javadoc.jenkins-ci.org/hudson/model/Result.html
	case "SUCCESS":
		return 0
	case "FAILURE":
		return 1
	case "ABORTED":
		return 2
	case "UNSTABLE":
		return 3
	case "NOT_BUILT":
		return 4
	default:
		return 1
	}
}
