package mechanic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Giovanni-Tedesco/tmitracksapi/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MechanicHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mechanics Writer")
}

func CreateReport(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

	reports := db.Collection("Reports")

	var report entity.Report

	// report.Duration = "00:20"
	// report.Date = primitive.Timestamp{T: uint32(time.Now().Unix())}
	// report.Equipment = "Cat 235"
	// report.Notes = "No notes 3"
	// report.Technician = "Andrei"
	// report.ID = primitive.NewObjectID()
	err := json.NewDecoder(r.Body).Decode(&report)

	if err != nil {
		log.Fatal(err)
	}

	report.ID = primitive.NewObjectID()

	insertResults, err := reports.InsertOne(context.TODO(), report)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Inserted Id: %v", insertResults.InsertedID)
}

func DeleteReport(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

}

func GetReportById(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

}

func GetReportByDateRange(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

}

type DateRequest struct {
	Date string `json:"Date,omitempty" validate:"required"`
}

// Request body of a date, returns all reports generated on that date.
func GetReportByDate(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

	collection := db.Collection("Reports")

	var req DateRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
	if err != nil {
		fmt.Println("Error here")
		log.Fatal(err)
	}

	fmt.Println(req.Date)

	query := bson.M{"date": req.Date}

	curr, err := collection.Find(context.TODO(), query)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	var reports []entity.Report

	for curr.Next(ctx) {
		var r entity.Report
		err := curr.Decode(&r)
		if err != nil {
			log.Fatal(err)
		}

		reports = append(reports, r)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)

}
