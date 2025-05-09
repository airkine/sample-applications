package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aaron-devsecops/go-azure-servicebusinternal/models"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
	client, _ := azservicebus.NewClientFromConnectionString(os.Getenv("SERVICEBUS_CONNECTION_STRING"), nil)
	receiver, _ := client.NewReceiverForQueue("order-queue", nil)

	log.Println("Queue Worker running...")

	for {
		messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		if err != nil {
			log.Println("Receive error:", err)
			continue
		}

		for _, msg := range messages {
			var order models.Order
			_ = json.Unmarshal(msg.Body, &order)
			log.Printf("[Queue Worker] Processing Order ID %s for %s (%d units)\n", order.ID, order.Item, order.Amount)
			_ = receiver.CompleteMessage(context.Background(), msg, nil)
		}
	}
}
