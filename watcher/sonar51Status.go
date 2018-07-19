package watcher

import (
	"encoding/json"
	"fmt"
	"github.com/g3force/qaBlink/config"
	"log"
	"net/http"
)

type Sonar51Msr struct {
	Key     string  `json:"key"`
	Val     float64 `json:"val"`
	FrmtVal string  `json:"frmt_val"`
}

type Sonar51Resource struct {
	Msr []Sonar51Msr `json:"msr"`
}

type Sonar51Response struct {
	Resource []Sonar51Resource `json:"projectStatus"`
}

type Sonar51Job struct {
	url   string
	state QaBlinkState
	id    string
	QaBlinkJob
}

func (job *Sonar51Job) State() QaBlinkState {
	return job.state
}

func (job *Sonar51Job) Id() string {
	return job.id
}

func (job *Sonar51Job) Update() {
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
	var resources []Sonar51Resource

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		log.Print("Could not decode json: ", resp.Body, err)
		job.state.StatusCode = UNKNOWN
		return
	}

	violations := 0.0
	for _, r := range resources {
		for _, msr := range r.Msr {
			violations += msr.Val
		}
	}

	if violations > 0 {
		job.state.StatusCode = UNSTABLE
	} else {
		job.state.StatusCode = STABLE
	}
}

func NewSonar51Job(config *config.SonarConfig, jobId string) *Sonar51Job {
	jobStatus := new(Sonar51Job)
	job, err := findSonarJob(config.Jobs, jobId)
	if err != nil {
		return nil
	}
	connection := findSonarConnection(config.Connections, job.ConnectionRef)
	jobStatus.url = fmt.Sprintf("%s/api/resources/index?metrics=violations&format=json&resource=%s",
		connection.BaseUrl, job.ProjectKey)
	jobStatus.state.StatusCode = UNKNOWN
	jobStatus.id = jobId
	return jobStatus
}
