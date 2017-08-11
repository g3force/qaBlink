package watcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/g3force/qaBlink/config"
	"log"
	"net/http"
	"strings"
)

type JenkinsResponseJobHealthReport struct {
	Score uint8 `json:"score"`
}

type JenkinsResponse struct {
	Color        string                           `json:"color"`
	InQueue      bool                             `json:"inQueue"`
	HealthReport []JenkinsResponseJobHealthReport `json:"healthReport"`
}

type JenkinsJob struct {
	url   string
	state QaBlinkState
	QaBlinkJob
}

func (job *JenkinsJob) State() QaBlinkState {
	return job.state
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

	color := jenkinsResponse.Color
	if strings.HasPrefix(color, "blue") {
		job.state.StatusCode = STABLE
	} else if strings.HasPrefix(color, "yellow") {
		job.state.StatusCode = UNSTABLE
	} else if strings.HasPrefix(color, "red") {
		job.state.StatusCode = FAILED
	} else if strings.HasPrefix(color, "disabled") {
		job.state.StatusCode = DISABLED
	} else {
		job.state.StatusCode = UNKNOWN
	}

	job.state.Pending = strings.HasSuffix(color, "anime")
}

func findJenkinsJob(jobs []config.JenkinsConfigJob, jobId uint8) (config.JenkinsConfigJob, error) {
	for _, job := range jobs {
		if job.Id == jobId {
			return job, nil
		}
	}
	return config.JenkinsConfigJob{}, errors.New("Job not found")
}

func findJenkinsConnection(connections []config.JenkinsConfigConnection, id uint8) config.JenkinsConfigConnection {
	for _, connection := range connections {
		if connection.Id == id {
			return connection
		}
	}
	panic("")
}

func NewJenkinsJob(config *config.JenkinsConfig, jobId uint8) *JenkinsJob {
	jobStatus := new(JenkinsJob)
	job, err := findJenkinsJob(config.Jobs, jobId)
	if err != nil {
		return nil
	}
	connection := findJenkinsConnection(config.Connections, job.ConnectionRef)
	jobStatus.url = fmt.Sprintf("%s/%s/api/json",
		connection.BaseUrl, job.JobName)
	jobStatus.state.StatusCode = UNKNOWN
	return jobStatus
}
