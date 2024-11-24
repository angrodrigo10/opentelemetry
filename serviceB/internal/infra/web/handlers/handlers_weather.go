package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Temperature struct {
	City   string  `json:"city"`
	Temp_c float64 `json:"temp_c"`
	Temp_f float64 `json:"temp_f"`
	Temp_k float64
}

type CurrentWeather struct {
	TemperatureC float64 `json:"temp_c"`
	Name         string  `json:"name"`
}

type WeatherResponse struct {
	Current  CurrentWeather `json:"current"`
	Location CurrentWeather `json:"location"`
}

func GetTemperature(ctx context.Context, city string) (*Temperature, error) {
	apiKey := "1da14e1d67344108ab9194641242010"

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	resp, err := client.Get("http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + city)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch weather data")
	}

	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return nil, err
	}

	temp := &Temperature{
		City:   weatherResp.Location.Name,
		Temp_c: weatherResp.Current.TemperatureC,
		Temp_f: celsiusToFahrenheit(weatherResp.Current.TemperatureC),
		Temp_k: celsiusToKelvin(weatherResp.Current.TemperatureC),
	}

	return temp, nil
}

func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273.15
}
