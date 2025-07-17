package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProvisionRequest struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CPU     int                `json:"cpu"`
	RAM     int                `json:"ram"`
	OS      string             `json:"os"`
	Project string             `json:"project"`
	Status  string             `json:"status"`
	Steps   []string           `json:"steps"`
}
