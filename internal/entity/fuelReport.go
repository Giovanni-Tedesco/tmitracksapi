package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FuelReport struct {
	ID          primitive.ObjectID `json:"id" validate:"required"`
	Date        string             `json:"date" validate:"required"`
	Fuel        int64              `json:"fuel" validate:"required"`
	Location    string             `json:"location" validate:"required"`
	Oil         int64              `json:"oil" validate:"required"`
	Hydraulic   int64              `json:"hydraulic" validate:"required"`
	Coolant     int64              `json:"coolant" validate:"required"`
	AirFilter   int64              `json:"airfilter" validate:"required"`
	TubesGrease int64              `json:"tubesgrease" validate:"required"`
	WasherFluid int64              `json:"washerfluid" validate:"required"`
}
