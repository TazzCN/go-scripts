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
	"github.com/tazzcn/gocode/googleTypes"
)

// This script takes a feature collection of google reviews and puts them
// in a dynamodb table

const DYNAMO_DB_TABLE = "google-places"

func load_aws_profile(profile string) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func get_features_from_json_file() (Features []googleTypes.Feature) {
	jsonFile, err := os.Open("./add-googleplaces/Reviews.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var featureCollection googleTypes.FeatureCollection
	json.Unmarshal(byteValue, &featureCollection)
	return featureCollection.Features
}

// example input: https://www.google.com/maps/place//data=!4m2!3m1!1s0x0:0x2312ed3941bae624
// ouput: 0x2312ed3941bae624
func extract_unique_id_from_google_url(googleUrl string) (id string) {
	string_split_by_colon := strings.Split(googleUrl, ":")
	return string_split_by_colon[2]
}

func put_feature_into_dynamodb(svc *dynamodb.Client, feature googleTypes.Feature) (out *dynamodb.PutItemOutput, err error) {
	var id = extract_unique_id_from_google_url(feature.Properties.GoogleMapsURL)
	var name = feature.Properties.Location.BusinessName
	fmt.Println(id)
	d := googleTypes.DynamoDBPlace{
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
