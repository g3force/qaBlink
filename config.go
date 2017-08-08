package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type Slot struct {
	Id    uint8   `json:"id"`
	RefId []uint8 `json:"refs"`
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
