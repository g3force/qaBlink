package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

type JenkinsResponseJobHealthReport struct {
	Score uint8 `json:"score"`
}

type JenkinsResponse struct {
	Color        string   `json:"color"`
	InQueue      bool `json:"inQueue"`
	HealthReport [] JenkinsResponseJobHealthReport `json:"healthReport"`
}

type JenkinsConfigConnection struct {
	Id      uint8 `json:"id"`
	User    string `json:"user"`
	Token   string `json:"token"`
	BaseUrl string `json:"baseUrl"`
}

type JenkinsConfigJob struct {
	Id            uint8 `json:"id"`
	JobName       string `json:"jobName"`
	ConnectionRef uint8 `json:"connectionRef"`
}

type JenkinsConfig struct {
	Connections [] JenkinsConfigConnection `json:"connections"`
	Jobs        [] JenkinsConfigJob `json:"jobs"`
}

type JenkinsJob struct {
	url   string
	state QaBlinkState
}

func (job *JenkinsJob) Update() {
	// Build the request
	req, err := http.NewRequest("GET", job.url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Request failed: ", err)
		job.state.StatusCode = UNKNOWN
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var jenkinsResponse JenkinsResponse

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&jenkinsResponse); err != nil {
		log.Print(err)
		job.state.StatusCode = UNKNOWN
		return
	}

	if len(jenkinsResponse.HealthReport) > 0 {
		job.state.Score = jenkinsResponse.HealthReport[0].Score
	}

	switch jenkinsResponse.Color {
	case "blue":
		job.state.StatusCode = STABLE
	case "yellow":
		job.state.StatusCode = UNSTABLE
	case "red":
		job.state.StatusCode = FAILED
	case "disabled":
		job.state.StatusCode = DISABLED
	default:
		job.state.StatusCode = UNKNOWN
	}
}

func findJob(jobs [] JenkinsConfigJob, jobId uint8) JenkinsConfigJob {
	for _, job := range jobs {
		if job.Id == jobId {
			return job
		}
	}
	panic("")
}

func findConnection(connections [] JenkinsConfigConnection, id uint8) JenkinsConfigConnection {
	for _, connection := range connections {
		if connection.Id == id {
			return connection
		}
	}
	panic("")
}

func NewJenkinsJob(config *JenkinsConfig, jobId uint8) *JenkinsJob {
	jobStatus := new(JenkinsJob)
	job := findJob(config.Jobs, jobId)
	connection := findConnection(config.Connections, job.ConnectionRef)
	jobStatus.url = fmt.Sprintf("https://%s:%s@%s/job/%s/api/json",
		connection.User, connection.Token, connection.BaseUrl, job.JobName)
	jobStatus.state.StatusCode = UNKNOWN
	return jobStatus
}
