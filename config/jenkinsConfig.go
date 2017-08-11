package config

type JenkinsConfigConnection struct {
	Id      uint8  `json:"id"`
	User    string `json:"user"`
	Token   string `json:"token"`
	BaseUrl string `json:"baseUrl"`
}

type JenkinsConfigJob struct {
	Id            uint8  `json:"id"`
	JobName       string `json:"jobName"`
	ConnectionRef uint8  `json:"connectionRef"`
}

type JenkinsConfig struct {
	Connections []JenkinsConfigConnection `json:"connections"`
	Jobs        []JenkinsConfigJob        `json:"jobs"`
}
