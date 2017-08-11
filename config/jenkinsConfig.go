package config

type JenkinsConfigConnection struct {
	Id      string `json:"id"`
	BaseUrl string `json:"baseUrl"`
}

type JenkinsConfigJob struct {
	Id            string `json:"id"`
	JobName       string `json:"jobName"`
	ConnectionRef string `json:"connectionRef"`
}

type JenkinsConfig struct {
	Connections []JenkinsConfigConnection `json:"connections"`
	Jobs        []JenkinsConfigJob        `json:"jobs"`
}
