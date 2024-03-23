package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rohanhonnakatti/go-graphql-mongodb/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file: %w", err)
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		log.Fatal("API KEY NOT PRESENT IN .env...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database...")

	return &DB{
		client: client,
	}
}

func (db *DB) GetJob(jobId string) *model.JobListing {

	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}

	var jobListing model.JobListing
	err := jobCollec.FindOne(ctx, filter).Decode(&jobListing)
	if err != nil {
		log.Fatal(err)
	}

	return &jobListing
}

func (db *DB) GetJobs() []*model.JobListing {

	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var jobListings []*model.JobListing
	cursor, err := jobCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &jobListings); err != nil {
		panic(err)
	}

	return jobListings
}

func (db *DB) CreateJobListing(inputModel model.CreateJobListing) *model.JobListing {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	inserted, err := jobCollec.InsertOne(ctx, bson.M{
		"title":       inputModel.Title,
		"description": inputModel.Description,
		"company":     inputModel.Company,
		"url":         inputModel.URL,
	})
	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserted.InsertedID.(primitive.ObjectID).Hex()
	jobListing := model.JobListing{
		ID:      insertedID,
		Title:   inputModel.Title,
		Company: inputModel.Company,
		URL:     inputModel.URL,
	}

	return &jobListing
}

func (db *DB) UpdateJobListing(jobId string, inputModel model.UpdateJobListing) *model.JobListing {

	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	updateJobInfo := bson.M{}

	if inputModel.Title != nil {
		updateJobInfo["title"] = inputModel.Title
	}

	if inputModel.Description != nil {
		updateJobInfo["description"] = inputModel.Description
	}

	if inputModel.URL != nil {
		updateJobInfo["url"] = inputModel.URL
	}

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo}

	var jobListing model.JobListing
	results := jobCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
	if err := results.Decode(&jobListing); err != nil {
		log.Fatal(err)
	}

	return &jobListing
}

func (db *DB) DeleteJobListing(jobId string) *model.DeleteJobResponse {

	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	_, err := jobCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	return &model.DeleteJobResponse{
		DeleteJobID: jobId,
	}
}
