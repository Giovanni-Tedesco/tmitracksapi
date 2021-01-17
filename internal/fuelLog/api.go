package fuelLog

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Giovanni-Tedesco/tmitracksapi/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOneFuelReport(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	fuellogs := db.Collection("FuelLogs")

	var report entity.FuelReport

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&report)

	if err != nil {
		log.Fatal(err)
	}

	report.ID = primitive.NewObjectID()

	res, err := fuellogs.InsertOne(context.TODO(), report)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}

// Create batch reports
func CreateBatchReports(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

}

func ReadAllReport(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

}

func ReadSinlgeReport(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

}

func ReadReportsInRange(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

}

func UpdateReport(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

}

func DeleteReport(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

}
