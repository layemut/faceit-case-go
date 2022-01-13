package service

import (
	"log"

	"github.com/layemut/faceit-case-go/notify"
)

type NotificationService struct {
	PubSub *notify.Pubsub
}

func (ns *NotificationService) SubscribeUserCreateEvent() {
	ch := ns.PubSub.Subscribe("create")
	log.Println("Subscribed to user save event")
	go func() {
		for {
			val := <-ch
			log.Printf("User created: ID %v, Name: %v, sending mail notification to %v", val.ID, val.FirstName, val.Email)
		}
	}()
}

func (ns *NotificationService) SubscribeUserUpdateEvent() {
	ch := ns.PubSub.Subscribe("update")
	log.Println("Subscribed to user update event")
	go func() {
		for {
			val := <-ch
			log.Printf("User update: ID %v, Name: %v, sending mail notification to %v", val.ID, val.FirstName, val.Email)
		}
	}()
}
