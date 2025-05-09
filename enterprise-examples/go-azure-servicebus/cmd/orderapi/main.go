package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aaron-devsecops/go-azure-servicebus/internal/messaging"
	"github.com/aaron-devsecops/go-azure-servicebusinternal/models"

	"github.com/google/uuid"
)

func main() {
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}
		var order models.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		order.ID = uuid.NewString()

		// Send to Queue & Publish to Topic
		_ = messaging.SendToQueue(order)
		_ = messaging.PublishToTopic(order)

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(order)
	})

	log.Println("Order API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
