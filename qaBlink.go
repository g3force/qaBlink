package main

import (
	"time"
	"log"
	"github.com/hink/go-blink1"
)

type QaBlinkSlot struct {
	Id   uint8
	Jobs []QaBlinkJob
}

type QaBlink struct {
	UpdateInterval uint32
	Slots          [] QaBlinkSlot
	Blink1Device   *blink1.Device
}

func (*QaBlinkState) Update() {}

func NewQaBlink(config *QaBlinkConfig) *QaBlink {
	qaBlink := new(QaBlink)
	qaBlink.UpdateInterval = config.UpdateInterval
	for _, slot := range config.Slots {
		var qaSlot QaBlinkSlot
		qaSlot.Id = slot.Id
		for _, refId := range slot.RefId {
			qaSlot.Jobs = append(qaSlot.Jobs, NewJenkinsJob(config.Jenkins, refId))
		}
		qaBlink.Slots = append(qaBlink.Slots, qaSlot)
	}
	return qaBlink
}

func (qaBlink *QaBlink) UpdateStatus() {
	for {
		log.Printf("Updating %d slots\n", len(qaBlink.Slots))
		for _, slot := range qaBlink.Slots {
			for _, job := range slot.Jobs {
				job.Update()
				log.Printf("%d: %v [%v]", slot.Id, job.State().Score, job.State().StatusCode)
			}
		}
		time.Sleep(time.Duration(qaBlink.UpdateInterval) * time.Second)
	}
}

func toState(state QaBlinkState) blink1.State {
	if state.Pending {
		return blink1.State{Red: 0, Green: 0, Blue: 255}
	}
	switch state.StatusCode {
	case STABLE:
		return blink1.State{Red: 0, Green: 255, Blue: 0}
	case UNSTABLE:
		return blink1.State{Red: 255, Green: 255, Blue: 0}
	case FAILED:
		return blink1.State{Red: 255, Green: 0, Blue: 0}
	case UNKNOWN:
		return blink1.State{Red: 0, Green: 0, Blue: 0}
	case DISABLED:
		return blink1.State{Red: 255, Green: 0, Blue: 255}
	}
	return blink1.State{}
}

func (qaBlink *QaBlink) UpdateBlink() {
	perSlotDuration := time.Duration(500) * time.Millisecond
	for {
		for _, slot := range qaBlink.Slots {
			for id, job := range slot.Jobs {
				state := toState(job.State())
				state.FadeTime = time.Duration(100) * time.Millisecond
				if len(slot.Jobs) == 1 {
					state.LED = blink1.LEDAll
				} else {
					switch id {
					case 0:
						state.LED = blink1.LED1
					case 1:
						state.LED = blink1.LED2
					default:
						state.LED = blink1.LEDAll
					}
				}
				qaBlink.Blink1Device.SetState(state)
			}
			time.Sleep(perSlotDuration)
		}
	}
}

func main() {

	config := NewQaBlinkConfig("config.json")
	qaBlink := NewQaBlink(config)

	device, err := blink1.OpenNextDevice()
	if err != nil {
		log.Fatal(err)
		return
	}
	qaBlink.Blink1Device = device

	//var state = blink1.State{Duration: time.Duration(10) * time.Millisecond, Red: 255, Blue: 0, Green: 0, LED: blink1.LED1, FadeTime: time.Duration(10) * time.Millisecond}
	//var state = blink1.State{Blue: 255}
	//device.SetState(state)
	//time.Sleep(time.Duration(1) * time.Second)

	go qaBlink.UpdateStatus()
	qaBlink.UpdateBlink()

	device.Close()
}
