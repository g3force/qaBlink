package config

import (
	"bufio"
	"encoding/json"
	"os"
)

type Slot struct {
	Id    string   `json:"id"`
	RefId []string `json:"refs"`
}

type QaBlinkConfig struct {
	UpdateInterval  uint32         `json:"updateInterval"`
	FadeTime        uint32         `json:"fadeTime"`
	PerSlotDuration uint32         `json:"perSlotDuration"`
	Slots           []Slot         `json:"slots"`
	Jenkins         *JenkinsConfig `json:"jenkins"`
	Sonar           *SonarConfig   `json:"sonar"`
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
	if config.UpdateInterval == 0 {
		config.UpdateInterval = 30
	}
	if config.FadeTime == 0 {
		config.FadeTime = 50
	}
	if config.PerSlotDuration == 0 {
		config.PerSlotDuration = 100
	}
	return config
}
