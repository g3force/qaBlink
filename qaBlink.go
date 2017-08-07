package main

import (
	"time"
	"log"
)

type QaBlinkSlot struct {
	Status *JenkinsJob
	Slot
}

type QaBlink struct {
	UpdateInterval uint32
	Slots          [] QaBlinkSlot
}

func NewQaBlink(config *QaBlinkConfig) *QaBlink {
	qaBlink := new(QaBlink)
	qaBlink.UpdateInterval = config.UpdateInterval
	for _, slot := range config.Slots {

		var qaSlot QaBlinkSlot
		qaSlot.Id = slot.Id
		qaSlot.RefId = slot.RefId
		for _, refId := range slot.RefId {
			qaSlot.Status = NewJenkinsJob(config.Jenkins, refId)
		}
		qaBlink.Slots = append(qaBlink.Slots, qaSlot)
	}
	return qaBlink
}

func (qaBlink *QaBlink) Update() {
	for {
		log.Printf("Updating %d slots\n", len(qaBlink.Slots))
		for _, slot := range qaBlink.Slots {
			slot.Status.Update()
			log.Printf("%d: %v [%v]", slot.Id, slot.Status.state.Score, slot.Status.state.StatusCode)
		}
		time.Sleep(time.Duration(qaBlink.UpdateInterval) * time.Second)
	}
}

func main() {
	config := NewQaBlinkConfig("config.json")
	qaBlink := NewQaBlink(config)

	qaBlink.Update()
}
