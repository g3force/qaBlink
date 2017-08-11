package config

import (
	"bufio"
	"encoding/json"
	"os"
)

type Slot struct {
	Id       string   `json:"id"`
	RefId    []string `json:"refs"`
	DeviceId uint8    `json:"deviceId"`
}

type QaBlinkConfig struct {
	UpdateInterval uint32         `json:"updateInterval"`
	Slots          []Slot         `json:"slots"`
	Jenkins        *JenkinsConfig `json:"jenkins"`
	Sonar          *SonarConfig   `json:"sonar"`
}

func NewQaBlinkConfig(fileName string) *QaBlinkConfig {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	config := new(QaBlinkConfig)
	if err := json.NewDecoder(reader).Decode(&config); err != nil {
		panic(err)
	}
	return config
}
