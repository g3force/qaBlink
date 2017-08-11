package watcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/g3force/qaBlink/config"
	"log"
	"net/http"
)

type SonarResponseProjectStatus struct {
	Status string `json:"status"`
}

type SonarResponse struct {
	ProjectStatus SonarResponseProjectStatus `json:"projectStatus"`
}

type SonarJob struct {
	url   string
	state QaBlinkState
	QaBlinkJob
}

func (job *SonarJob) State() QaBlinkState {
	return job.state
}

func (job *SonarJob) Update() {
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
	var sonarResponse SonarResponse

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&sonarResponse); err != nil {
		log.Print("Could not decode json: ", err)
		job.state.StatusCode = UNKNOWN
		return
	}

	switch sonarResponse.ProjectStatus.Status {
	case "OK":
		job.state.StatusCode = STABLE
	case "WARN":
		job.state.StatusCode = UNSTABLE
	case "ERROR":
		job.state.StatusCode = FAILED
	case "NONE":
		job.state.StatusCode = DISABLED
	default:
		job.state.StatusCode = UNKNOWN
	}
}

func findSonarJob(jobs []config.SonarConfigJob, jobId string) (config.SonarConfigJob, error) {
	for _, job := range jobs {
		if job.Id == jobId {
			return job, nil
		}
	}
	return config.SonarConfigJob{}, errors.New("Job not found")
}

func findSonarConnection(connections []config.SonarConfigConnection, id string) config.SonarConfigConnection {
	for _, connection := range connections {
		if connection.Id == id {
			return connection
		}
	}
	panic("Sonar connection not found: " + id)
}

func NewSonarJob(config *config.SonarConfig, jobId string) *SonarJob {
	jobStatus := new(SonarJob)
	job, err := findSonarJob(config.Jobs, jobId)
	if err != nil {
		return nil
	}
	connection := findSonarConnection(config.Connections, job.ConnectionRef)
	jobStatus.url = fmt.Sprintf("%s/api/qualitygates/project_status?projectKey=%s",
		connection.BaseUrl, job.ProjectKey)
	jobStatus.state.StatusCode = UNKNOWN
	return jobStatus
}
