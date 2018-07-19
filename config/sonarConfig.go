package config

type SonarConfigConnection struct {
	Id       string `json:"id"`
	BaseUrl  string `json:"baseUrl"`
	Selector string `json:"selector"`
}

type SonarConfigJob struct {
	Id            string `json:"id"`
	ProjectKey    string `json:"projectKey"`
	ConnectionRef string `json:"connectionRef"`
}

type SonarConfig struct {
	Connections []SonarConfigConnection `json:"connections"`
	Jobs        []SonarConfigJob        `json:"jobs"`
}
