package main

import (
	"github.com/g3force/qaBlink/config"
	"github.com/g3force/qaBlink/watcher"
	"github.com/hink/go-blink1"
	"log"
	"time"
)

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
	for {
		log.Printf("Updating %d slots\n", len(qaBlink.Slots))
		for _, slot := range qaBlink.Slots {
			for jobId, job := range slot.Jobs {
				job.Update()
				log.Printf("%20s(job:%d): %8v [pending: %5v,score: %3v]", slot.Id, jobId, job.State().StatusCode, job.State().Pending, job.State().Score)
			}
		}
		time.Sleep(time.Duration(qaBlink.UpdateInterval) * time.Second)
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
	for {
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
}

func main() {

	blinkConfig := config.NewQaBlinkConfig("config.json")
	qaBlink := NewQaBlink(blinkConfig)

	for {
		device, err := blink1.OpenNextDevice()
		if device == nil {
			break
		}
		if err != nil {
			log.Print(err)
			break
		}
		device.SetState(blink1.State{Red: 255, Blue: 255})
		qaBlink.blink1Devices = append(qaBlink.blink1Devices, device)
	}

	log.Printf("Found %d blink1 devices.\n", len(qaBlink.blink1Devices))

	go qaBlink.UpdateStatus()
	go qaBlink.UpdateBlink()

	for {
		time.Sleep(time.Hour)
	}
}
