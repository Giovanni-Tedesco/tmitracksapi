package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Date       string             `json:"Date" validate:"omitempty,required"`
	Duration   string             `json:"Duration" validate:"omitempty,required"`
	Equipment  string             `json:"Equipment" validate:"omitempty,required"`
	Technician string             `json:"Technician,omitempty" validate:"omitempty,required"`
	Notes      string             `json:"Notes" validate:"omitempty,required"`
}
