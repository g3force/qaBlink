package main

type SonarConfigConnection struct {
	Id      uint8  `json:"id"`
	Token   string `json:"token"`
	BaseUrl string `json:"baseUrl"`
}

type SonarConfigJob struct {
	Id            uint8  `json:"id"`
	ProjectKey    string `json:"projectKey"`
	ConnectionRef uint8  `json:"connectionRef"`
}

type SonarConfig struct {
	Connections []SonarConfigConnection `json:"connections"`
	Jobs        []SonarConfigJob        `json:"jobs"`
}
