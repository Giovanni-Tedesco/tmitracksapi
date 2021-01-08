package mechanic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/Giovanni-Tedesco/tmitracksapi/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MechanicHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mechanics Writer")
}

func CreateReport(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

	reports := db.Collection("Reports")

	var report entity.Report

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
	type RequestBody struct {
		ID primitive.ObjectID `json:"id" validate:"required"`
	}

	var req RequestBody
	v := validator.New()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Request Disalloed. Invalid Field")
		return
	}

	err = v.Struct(req)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintf(w, "%v", e)
		}
		return
	}

	collection := db.Collection("Reports")

	var report entity.Report

	err = collection.FindOne(context.TODO(), bson.M{"_id": req.ID}).Decode(&report)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

func GetReportByDateRange(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		StartDate string `json:"start_date" validate:"required"`
		EndDate   string `json:"end_date" validate:"required"`
	}

	var req RequestBody
	v := validator.New()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Request Disalloed. Invalid Field")
	}

	err = v.Struct(req)

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Fprintf(w, "%v", e)
		}
		return
	}

	collection := db.Collection("Reports")

	query := bson.D{
		{"$and", []bson.M{
			bson.M{"date": bson.M{"$lte": req.EndDate}},
			bson.M{"date": bson.M{"$gte": req.StartDate}},
		}},
	}

	options := options.Find()
	options.SetSort(bson.D{{"date", 1}, {"duration", -1}})

	curr, err := collection.Find(context.TODO(), query, options)

	if err != nil {
		log.Fatal(err)
	}

	var reports []entity.Report
	for curr.Next(context.TODO()) {
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

func GetAllReports(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	collection := db.Collection("Reports")

	curr, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	var reports []entity.Report

	ctx := context.Background()

	for curr.Next(ctx) {
		var r entity.Report
		curr.Decode(&r)

		reports = append(reports, r)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)
}
