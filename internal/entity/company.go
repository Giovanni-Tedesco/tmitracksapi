package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID        primitive.ObjectID
	Name      string   `json:"name" validate:"required"`
	Employees []User   `json:"employees" validate:"required"`
	Reports   []Report `json:"reports" validate:"required"`
}
