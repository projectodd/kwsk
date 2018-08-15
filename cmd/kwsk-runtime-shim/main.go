package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var initialized int32
var serverHostAndPort string
var actionCode string
var actionBinary string
var actionParams map[string]interface{}

const (
	KwskActionCode   string = "KWSK_ACTION_CODE"
	KwskActionBinary string = "KWSK_ACTION_BINARY"
	KwskActionParams string = "KWSK_ACTION_PARAMS"
	PrintLogs        bool   = false
)

type ActionInitMessage struct {
	Value ActionInitValue `json:"value,omitempty"`
}

type ActionInitValue struct {
	Main   string `json:"main,omitempty"`
	Code   string `json:"code,omitempty"`
	Binary string `json:"binary,omitempty"`
}

type ActionRunMessage struct {
	Value interface{} `json:"value"`
}

func main() {
	if !PrintLogs {
		log.SetOutput(ioutil.Discard)
	}
	serverHostAndPort = "localhost:8081"
	actionCode = os.Getenv(KwskActionCode)
	actionBinary = os.Getenv(KwskActionBinary)
	if _, exists := os.LookupEnv(KwskActionParams); exists {
		err := json.Unmarshal([]byte(os.Getenv(KwskActionParams)), &actionParams)
		if err != nil {
			log.Fatal("Failed to load action params")
		}
	}

	http.HandleFunc("/", handler)
	addr := ":8080"
	log.Printf("Starting http server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Got request %+v\n", req)
	if atomic.LoadInt32(&initialized) == 0 {
		log.Println("Initializing action")
		initBody := &ActionInitMessage{
			Value: ActionInitValue{
				Main:   "main",
				Code:   actionCode,
				Binary: actionBinary,
			},
		}
		res, err := actionRequest(serverHostAndPort, "init", initBody)
		log.Printf("Response: %+v\n", res)
		if err != nil {
			log.Printf("Error initializing action: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if res.StatusCode == http.StatusForbidden {
			log.Println("Action already initialized")
		} else if res.StatusCode != http.StatusOK {
			log.Printf("Action initializer returned a HTTP error: %+v\n", res)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		atomic.AddInt32(&initialized, 1)
	}

	log.Println("Running action")

	decoder := json.NewDecoder(req.Body)
	var params map[string]interface{}
	err := decoder.Decode(&params)
	if err != nil {
		if err == io.EOF {
			// This means we had an empty request body so just make
			// empty parameters
			params = map[string]interface{}{}
		} else {
			log.Printf("Error decoding request body: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	runBody := &ActionRunMessage{
		Value: combineActionParameters(params),
	}
	res, err := actionRequest(serverHostAndPort, "run", runBody)
	if err != nil {
		log.Printf("Error running action: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	log.Printf("Response: %+v\n", res)
	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func actionRequest(hostAndPort string, path string, requestBody interface{}) (*http.Response, error) {
	url := fmt.Sprintf("http://%s/%s", hostAndPort, path)
	log.Printf("Sending POST to url %s\n", url)

	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshalling action request body: %s\n", err)
		return nil, err
	}
	log.Printf("Request Body: %s\n", body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating http request for action: %s\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func combineActionParameters(params map[string]interface{}) map[string]interface{} {
	combinedParams := map[string]interface{}{}
	for k, v := range actionParams {
		combinedParams[k] = v
	}
	for k, v := range params {
		combinedParams[k] = v
	}
	return combinedParams
}
