package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bin struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Number    int                `json:"Number"`
	XPosition float64            `json:"xPosition"`
	YPosition float64            `json:"yPosition"`
}
