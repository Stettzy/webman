package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Body       []byte `json:"body"`
}

func prepAndExecRequest(w http.ResponseWriter, r *http.Request) *http.Response {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return nil
	}

	var req Request

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return nil
	}

	// Define client as http.NewRequest must use client to execute the request
	client := &http.Client{}

	// Transform the body to io.Reader type as http.NewRequest requires this type
	bodyReader := bytes.NewReader(req.Body)

	newReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return nil
	}

	// Append the headers
	for k, v := range req.Headers {
		newReq.Header.Add(k, v)
	}

	// Fire the request
	resp, err := client.Do(newReq)
	if err != nil {
		http.Error(w, "Error executing request", http.StatusBadRequest)
		return nil
	}

	return resp
}

func returnResponse(w http.ResponseWriter, r *http.Response) {
	var res Response

	// We map the values to our response structure
	res.StatusCode = r.StatusCode
	res.Status = r.Status

	// Body must be parsed here
	parsedBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
	}
	res.Body = parsedBody

	// We have to convert to json in order to return it as response
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
	}

	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Write back with the response writer
	_, err = w.Write(resJson)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}

func startServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Main request entry point
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle outbound requests
		resp := prepAndExecRequest(w, r)
		if resp == nil {

			r, err := json.Marshal(map[string]string{"message": "fail"})
			if err != nil {
				return
			}

			w.Write(r)

			if err != nil {
				return
			}

			return
		}
		// Return result response
		returnResponse(w, resp)
	})

	fmt.Printf("Server listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startServer()
}
