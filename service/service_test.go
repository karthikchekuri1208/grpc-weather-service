package service

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	pb "github.com/karthikchekuri1208/grpc-weather-service/gen/github.com/karthikchekuri1208/grpc-weather-service/proto"
	"github.com/stretchr/testify/assert"
)

// Helper function to mock HTTP client responses
type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// Test Cases

func TestGetWeather_Success(t *testing.T) {
	server := NewWeatherServiceServer("valid_api_key")
	// Mock the HTTP client to return a valid response
	server.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"weather": [{"main": "Clear"}],
					"main": {"temp": 280.0}
				}`)),
			}
		}),
	}
	req := &pb.WeatherRequest{Lat: 37.7749, Long: -12.4194} // San Francisco coordinates

	resp, err := server.GetWeather(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, []string{"Cold", "Moderate", "Hot"}, resp.Temperature)
	assert.Equal(t, "Clear", resp.Condition)
}

func TestGetWeather_InvalidAPIKey(t *testing.T) {
	server := NewWeatherServiceServer("invalid_api_key")
	// Mock the HTTP client to return an unauthorized response
	server.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body: io.NopCloser(strings.NewReader(`{
					"weather": [{"test": "Clear"}],
					"main": {"temp": 280.0}
				}`)),
			}
		}),
	}
	req := &pb.WeatherRequest{Lat: 37.7749, Long: -122.4194}

	resp, err := server.GetWeather(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "weather condition not found")
}

func TestGetWeather_InvalidCoordinates(t *testing.T) {
	server := NewWeatherServiceServer("valid_api_key")
	// Mock the HTTP client to return a bad request response
	server.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: io.NopCloser(strings.NewReader(`{
					"weather": [{"main": "Clear"}],
					"mai": {"temp": 280.0}
				}`)),
			}
		}),
	}
	req := &pb.WeatherRequest{Lat: 1000, Long: 2000} // Invalid coordinates

	resp, err := server.GetWeather(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "temperature data not found in API response")
}

func TestGetWeather_InvalidAPIResponse(t *testing.T) {
	server := NewWeatherServiceServer("valid_api_key")
	// Mock the HTTP client to return invalid JSON
	server.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"invalid": "data"}`)),
			}
		}),
	}
	req := &pb.WeatherRequest{Lat: 37.7749, Long: -122.4194}

	resp, err := server.GetWeather(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to get the data from the API")
}

func TestGetWeather_EmptyWeatherData(t *testing.T) {
	server := NewWeatherServiceServer("valid_api_key")
	// Mock the HTTP client to return empty weather data
	server.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"weather": []}`)),
			}
		}),
	}
	req := &pb.WeatherRequest{Lat: 37.7749, Long: -122.4194}

	resp, err := server.GetWeather(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to get the data from the API")
}

func TestTempClass(t *testing.T) {
	assert.Equal(t, "Cold", tempClass(5.0))
	assert.Equal(t, "Moderate", tempClass(15.0))
	assert.Equal(t, "Hot", tempClass(35.0))
}
