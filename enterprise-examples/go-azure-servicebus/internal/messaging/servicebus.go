package messaging

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aaron-devsecops/go-azure-servicebus/internal/models"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

var connStr = os.Getenv("SERVICEBUS_CONNECTION_STRING")

func SendToQueue(order models.Order) error {
	client, err := azservicebus.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return err
	}
	sender, err := client.NewSender("order-queue", nil)
	if err != nil {
		return err
	}
	defer sender.Close(context.Background())

	data, _ := json.Marshal(order)
	msg := &azservicebus.Message{
		Body: data,
	}
	return sender.SendMessage(context.Background(), msg, nil)
}

func PublishToTopic(order models.Order) error {
	client, err := azservicebus.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return err
	}
	sender, err := client.NewSender("order-topic", nil)
	if err != nil {
		return err
	}
	defer sender.Close(context.Background())

	data, _ := json.Marshal(order)
	msg := &azservicebus.Message{
		Body: data,
	}
	return sender.SendMessage(context.Background(), msg, nil)
}

func ReceiveFromSubscription(subName string, handler func(models.Order)) {
	client, _ := azservicebus.NewClientFromConnectionString(connStr, nil)
	receiver, _ := client.NewReceiverForSubscription("order-topic", subName, nil)

	for {
		messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		if err != nil {
			log.Println("Receive failed:", err)
			continue
		}

		for _, msg := range messages {
			var order models.Order
			_ = json.Unmarshal(msg.Body, &order)
			handler(order)
			_ = receiver.CompleteMessage(context.Background(), msg, nil)
		}
	}
}
