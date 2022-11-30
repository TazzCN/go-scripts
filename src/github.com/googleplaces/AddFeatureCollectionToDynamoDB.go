package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// This script takes a feature collection of google reviews and puts them
// in a dynamodb table

const DYNAMO_DB_TABLE = "google-places"

// Feature Collection Types
type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Geometry   Geometry   `json:"geometry"`
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	GoogleMapsURL string     `json:"Google Maps URL"`
	Published     string     `json:"Published"`
	ReviewComment string     `json:"Review Comment"`
	StarRating    int        `json:"Star Rating"`
	Location      Location   `json:"Location"`
	Questions     []Question `json:"Questions"`
}

type Question struct {
	Question       string `json:"Question"`
	SelectedOption string `json:"Selected Option"`
	Rating         string `json:"Rating"`
}

type Location struct {
	Address        string         `json:"Address"`
	BusinessName   string         `json:"Business Name"`
	CountryCode    string         `json:"Country Code"`
	GeoCoordinates GeoCoordinates `json:"Geo Coordinates"`
}

type GeoCoordinates struct {
	Latitude  string `json:"Latitude"`
	Longitude string `json:"Longitude"`
}

type DynamoDBPlace struct {
	ID          string     `dynamodbav:"id"`
	Name        string     `dynamodbav:"name"`
	Published   string     `dynamodbav:"published"`
	Coordinates []float64  `dynamodbav:"coordinates"`
	URL         string     `dynamodbav:"url"`
	Latitude    string     `dynamodbav:"latitude"`
	Longitude   string     `dynamodbav:"longitude"`
	Location    Location   `dynamodbav:"location"`
	Questions   []Question `dynamodbav:"questions"`
	Comment     string     `dynamodbav:"comment"`
	Rating      int        `dynamodbav:"rating"`
}

func load_aws_profile(profile string) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func get_features_from_json_file() (Features []Feature) {
	jsonFile, err := os.Open("./src/github.com/googleplaces/Reviews.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var featureCollection FeatureCollection
	json.Unmarshal(byteValue, &featureCollection)
	return featureCollection.Features
}

// example input: https://www.google.com/maps/place//data=!4m2!3m1!1s0x0:0x2312ed3941bae624
// ouput: 0x2312ed3941bae624
func extract_unique_id_from_google_url(googleUrl string) (id string) {
	string_split_by_colon := strings.Split(googleUrl, ":")
	return string_split_by_colon[2]
}

func put_feature_into_dynamodb(svc *dynamodb.Client, feature Feature) (out *dynamodb.PutItemOutput, err error) {
	var id = extract_unique_id_from_google_url(feature.Properties.GoogleMapsURL)
	var name = feature.Properties.Location.BusinessName
	fmt.Println(id)
	d := DynamoDBPlace{
		ID:          id,
		Name:        name,
		Published:   feature.Properties.Published,
		URL:         feature.Properties.GoogleMapsURL,
		Latitude:    feature.Properties.Location.GeoCoordinates.Latitude,
		Longitude:   feature.Properties.Location.GeoCoordinates.Longitude,
		Coordinates: feature.Geometry.Coordinates,
		Location:    feature.Properties.Location,
		Questions:   feature.Properties.Questions,
		Comment:     feature.Properties.ReviewComment,
		Rating:      feature.Properties.StarRating,
	}
	av, av_err := attributevalue.MarshalMap(d)
	if av_err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}
	return svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(DYNAMO_DB_TABLE),
		Item:      av,
	})
}

func main() {
	features := get_features_from_json_file()
	fmt.Printf("Uploading %d features to dynamoDB table\n", len(features))
	// Loading personal profile for AWS
	cfg := load_aws_profile("personal")
	svc := dynamodb.NewFromConfig(cfg)

	for _, s := range features {
		_, err := put_feature_into_dynamodb(svc, s)

		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
	fmt.Println("Finished uploading items")
}
