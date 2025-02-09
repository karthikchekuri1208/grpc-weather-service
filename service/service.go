package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	pb "github.com/karthikchekuri1208/grpc-weather-service/gen/github.com/karthikchekuri1208/grpc-weather-service/proto"
)

type WeatherServiceServer struct {
	pb.UnimplementedWeatherServiceServer
	apiKey string
	client *http.Client
}

func NewWeatherServiceServer(apiKey string) *WeatherServiceServer {
	return &WeatherServiceServer{
		apiKey: apiKey,
		client: http.DefaultClient,
	}
}

func (s *WeatherServiceServer) GetWeather(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	// Calling OpenWeather API
	//https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid={API key}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", req.Lat, req.Long, s.apiKey)
	fmt.Println("Here is the URL: ", url)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenWeather API: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	weatherData, ok := result["weather"].([]interface{})
	if !ok || len(weatherData) == 0 {
		return nil, fmt.Errorf("failed to get the data from the API")
	}

	weatherMap, ok := weatherData[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid weather data format")
	}

	weather, ok := weatherMap["main"].(string)
	if !ok {
		return nil, fmt.Errorf("weather condition not found")
	}

	mainData, ok := result["main"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("temperature data not found in API response")
	}

	tempVal, ok := mainData["temp"].(float64)
	if !ok {
		return nil, fmt.Errorf("temperature value is missing or incorrect format")
	}

	// Converting temperature from Kelvin to Celsius
	temp := tempVal - 273.15
	tempClass := tempClass(temp)

	return &pb.WeatherResponse{
		Condition:   weather,
		Temperature: tempClass,
	}, nil
}

// tempClass classifies temperature class based on the temperatures.
func tempClass(temp float64) string {

	var tstring string

	if temp < 10 {
		tstring = "Cold"
	} else if temp >= 10 && temp < 30 {
		tstring = "Moderate"
	} else {
		tstring = "Hot"
	}
	return tstring
}
