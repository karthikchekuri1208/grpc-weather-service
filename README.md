# grpc-weather-service
This project exposes a gRPC service as an HTTP API to fetch weather data using the OpenWeather API. It uses gRPC and gRPC Gateway to translate RESTful HTTP/JSON requests into gRPC calls.

---

## Features
- Fetch weather conditions for a given latitude and longitude by calling OpenWeather API.
- Classify the temperature as Hot(>30 Celsius), Moderate(10-30 Celsius), or Cold(<10 Celsius) based on predefined ranges.
- Expose the gRPC service as an HTTP API using the gRPC Gateway.

---

## Prerequisites
1. **Go** (version 1.23.4)
2. **Protocol Buffers Compiler** (`protoc`) (libprotoc 29.3)
3. **OpenWeather API Key** (sign up at [OpenWeather](https://openweathermap.org/api))

---

## Setup
**Install Dependencies**
Install the required tools and plugins:
```bash
# Install protoc (Protocol Buffers Compiler)
# Download from: https://github.com/protocolbuffers/protobuf/releases

# Install Go plugins for protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```
---

## Test the Project by following:

1. Clone the project in your local.
   - git clone https://github.com/karthikchekuri1208/grpc-weather-service.git 


2. Set your OpenWeather API key as the environment variable:
   - For powershell : $env:OPENWEATHER_API_KEY = "your_api_key"
   - For bash       : export OPENWEATHER_API_KEY="your_api_key"


3.  Run the grpc server and http gateway:
    - go run main.go

4. One should see a message of grpc server and HTTP gateway running in port 8080.

5. One can test the API call via browser or curl from the command prompt.
   - Example: http://localhost:8080/v1/weather?lat=44&long=-25
   - Response : { "condition": "Clouds", "temperature": "Moderate"}
   