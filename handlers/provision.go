package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"mini-poa/db"
	"mini-poa/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProvisionRequest(w http.ResponseWriter, r *http.Request) {
	var request models.ProvisionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	request.Status = "pending"
	request.Steps = []string{"Request received"}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ProvisionCollection.InsertOne(ctx, request)
	if err != nil {
		http.Error(w, "Failed to insert", http.StatusInternalServerError)
		return
	}

	// Fix: properly assign ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Failed to parse ID", http.StatusInternalServerError)
		return
	}
	request.ID = objectID

	go simulateProvisioningWorkflow(objectID)

	// Set response header and return JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}

func simulateProvisioningWorkflow(id primitive.ObjectID) {
	update := func(step, status string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Push a new step and update status
		update := bson.M{
			"$push": bson.M{"steps": step},
			"$set":  bson.M{"status": status},
		}
		_, err := db.ProvisionCollection.UpdateByID(ctx, id, update)
		if err != nil {
			log.Printf("⚠️ Failed to update status for ID %s: %v\n", id.Hex(), err)
		} else {
			log.Printf("✅ %s → %s\n", step, status)
		}
	}

	time.Sleep(10 * time.Second)
	update("Provisioning started", "in_progress")

	time.Sleep(10 * time.Second)
	update("Resource allocated", "in_progress")

	time.Sleep(10 * time.Second)
	update("Provisioning completed", "completed")
}

func GetAllProvisionRequests(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.ProvisionCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var requests []models.ProvisionRequest
	if err = cursor.All(ctx, &requests); err != nil {
		http.Error(w, "Failed to decode records", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}
