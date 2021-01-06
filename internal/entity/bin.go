package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bin struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Number    int                `json:Number`
	xPosition float64            `json:xPosition`
	yPosition float64            `json:yPosition`
}
