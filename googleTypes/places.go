package googleTypes

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
