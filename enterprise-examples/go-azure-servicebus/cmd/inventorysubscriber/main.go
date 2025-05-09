package main

import (
	"log"

	"github.com/aaron-devsecops/go-azure-servicebus/internal/messaging"
	"github.com/aaron-devsecops/go-azure-servicebus/internal/models"
)

func main() {
	log.Println("Email Subscriber running...")
	messaging.ReceiveFromSubscription("email-subscription", func(order models.Order) {
		log.Printf("[Email Service] Order ID %s: Email sent for item %s\n", order.ID, order.Item)
	})
}
