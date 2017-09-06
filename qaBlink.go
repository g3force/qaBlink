package main

import (
	"fmt"
	"github.com/g3force/go-blink1"
	"github.com/g3force/qaBlink/config"
	"github.com/g3force/qaBlink/watcher"
	"log"
	"os"
	"os/exec"
	"time"
)

var CONFIG_LOCATIONS = []string{"config.json", os.Getenv("HOME") + "/.qaBlink.json"}

type QaBlinkSlot struct {
	Id   string
	Jobs []watcher.QaBlinkJob
}

type QaBlink struct {
	UpdateInterval uint32
	Slots          []QaBlinkSlot
	blink1Devices  []*blink1.Device
}

func NewQaBlink(config *config.QaBlinkConfig) *QaBlink {
	qaBlink := new(QaBlink)
	qaBlink.UpdateInterval = config.UpdateInterval
	for _, slot := range config.Slots {
		var qaSlot QaBlinkSlot
		qaSlot.Id = slot.Id
		for _, refId := range slot.RefId {
			var jenkinsJob = watcher.NewJenkinsJob(config.Jenkins, refId)
			if jenkinsJob != nil {
				qaSlot.Jobs = append(qaSlot.Jobs, jenkinsJob)
			} else {
				var sonarJob = watcher.NewSonarJob(config.Sonar, refId)
				if sonarJob != nil {
					qaSlot.Jobs = append(qaSlot.Jobs, sonarJob)
				}
			}
		}
		qaBlink.Slots = append(qaBlink.Slots, qaSlot)
	}
	return qaBlink
}

func (qaBlink *QaBlink) UpdateStatus() {
	log.Printf("Updating %d slots\n", len(qaBlink.Slots))
	for _, slot := range qaBlink.Slots {
		for jobId, job := range slot.Jobs {
			job.Update()
			log.Printf("%40s(job:%d): %8v [pending: %5v,score: %3v]", job.Id(), jobId, job.State().StatusCode, job.State().Pending, job.State().Score)
		}
	}
}

func toState(state watcher.QaBlinkState) blink1.State {
	if state.Pending {
		return blink1.State{Red: 0, Green: 0, Blue: 255}
	}
	switch state.StatusCode {
	case watcher.STABLE:
		return blink1.State{Red: 0, Green: 255, Blue: 0}
	case watcher.UNSTABLE:
		return blink1.State{Red: 255, Green: 255, Blue: 0}
	case watcher.FAILED:
		return blink1.State{Red: 255, Green: 0, Blue: 0}
	case watcher.UNKNOWN:
		return blink1.State{Red: 0, Green: 0, Blue: 0}
	case watcher.DISABLED:
		return blink1.State{Red: 255, Green: 0, Blue: 255}
	}
	return blink1.State{}
}

func (qaBlink *QaBlink) UpdateBlink() {
	perSlotDuration := time.Duration(500) * time.Millisecond
	for _, slot := range qaBlink.Slots {
		slotId := 0
		for _, device := range qaBlink.blink1Devices {
			for ledId := 0; ledId < 2; ledId++ {
				var state blink1.State
				if slotId < len(slot.Jobs) {
					job := slot.Jobs[slotId]
					state = toState(job.State())
				} else {
					state = blink1.State{}
				}

				state.FadeTime = time.Duration(100) * time.Millisecond
				switch ledId {
				case 0:
					state.LED = blink1.LED1
				case 1:
					state.LED = blink1.LED2
				default:
					continue
				}

				device.SetState(state)
				slotId++
			}
		}
		time.Sleep(perSlotDuration)
	}
}

func (qaBlink *QaBlink) UpdateDevices() {

	newDevices := 0
	for i := 0; ; i++ {
		err := exec.Command("blink1-tool", "--red", "-d", fmt.Sprintf("%d", i)).Run()
		if err != nil && err.Error() != "exit status 1" {
			log.Print("Could not activate blink-devices by calling blink1-tool", err)
		}
		device, err := blink1.OpenNextDevice()
		if device == nil {
			break
		}
		if err != nil {
			log.Print(err)
			break
		}
		device.SetState(blink1.State{Red: 255, Blue: 255})
		newDevices++
	}

	qaBlink.blink1Devices = blink1.OpenDevices()

	if newDevices > 0 {
		log.Printf("Found %d new blink1 devices, %d now.\n", newDevices, len(qaBlink.blink1Devices))
	}
}

func repeat(f func(), duration time.Duration) {
	for {
		time.Sleep(duration)
		f()
	}
}

func main() {

	chosenConfig := ""
	for _, configLocation := range CONFIG_LOCATIONS {
		if _, err := os.Stat(configLocation); !os.IsNotExist(err) {
			chosenConfig = configLocation
			break
		}
	}
	blinkConfig := config.NewQaBlinkConfig(chosenConfig)
	qaBlink := NewQaBlink(blinkConfig)

	statusUpdateInterval := time.Duration(qaBlink.UpdateInterval) * time.Second
	deviceUpdateInterval := statusUpdateInterval

	go qaBlink.UpdateDevices()
	qaBlink.UpdateStatus()

	go repeat(qaBlink.UpdateStatus, statusUpdateInterval)
	go repeat(qaBlink.UpdateBlink, 0)
	go repeat(qaBlink.UpdateDevices, deviceUpdateInterval)

	for {
		time.Sleep(time.Hour)
	}
}
