package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/subosito/gotenv"
)

var darkSkyAPIKey string
var hereMapsAppID string
var hereMapsAppCode string

var city string

type clientInfo struct {
	City string `json:"city"`
}

type geocodeInfo struct {
	Response struct {
		View []struct {
			Result []struct {
				Location struct {
					DisplayPosition struct {
						Latitude  float64 `json:"Latitude"`
						Longitude float64 `json:"Longitude"`
					} `json:"DisplayPosition"`
				} `json:"Location"`
			} `json:"Result"`
		} `json:"View"`
	} `json:"Response"`
}

type weatherInfo struct {
	Currently struct {
		Time                 float64 `json:"time"`
		Summary              string  `json:"summary"`
		Icon                 string  `json:"icon"`
		NearestStormDistance float64 `json:"nearestStormDistance"`
		NearestStormBearing  float64 `json:"nearestStormBearing"`
		PrecipIntensity      float64 `json:"precipIntensity"`
		PrecipProbability    float64 `json:"precipProbability"`
		Temperature          float64 `json:"temperature"`
		ApparentTemperature  float64 `json:"apparentTemperature"`
		DewPoint             float64 `json:"dewPoint"`
		Humidity             float64 `json:"humidity"`
		Pressure             float64 `json:"pressure"`
		WindSpeed            float64 `json:"windSpeed"`
		WindGust             float64 `json:"windGust"`
		WindBearing          float64 `json:"windBearing"`
		CloudCover           float64 `json:"cloudCover"`
		UvIndex              float64 `json:"uvIndex"`
		Visibility           float64 `json:"visibility"`
		Ozone                float64 `json:"ozone"`
	} `json:"currently"`
}

func init() {
	if err := gotenv.Load(); err != nil {
		log.Fatal(err)
	}
	darkSkyAPIKey = os.Getenv("DARKSKY_API_KEY")
	hereMapsAppID = os.Getenv("HEREMAPS_APP_ID")
	hereMapsAppCode = os.Getenv("HEREMAPS_APP_CODE")
	flag.StringVar(&city, "city", "", "Address to fetch weather for. If omitted the program will use location based on IP address.")
}

func main() {
	flag.Parse()

	var client clientInfo

	if city == "" {
		resp, err := http.Get("http://ipinfo.io")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&client)
		if err != nil {
			log.Fatal(err)
		}
		city = client.City
	}

	encodedCity := url.QueryEscape(city)
	locationURL := "https://geocoder.api.here.com/6.2/geocode.json?" +
		"app_id=" + hereMapsAppID +
		"&app_code=" + hereMapsAppCode +
		"&searchtext=" + encodedCity

	var geocode geocodeInfo

	resp, err := http.Get(locationURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&geocode)
	if err != nil {
		log.Fatal(err)
	}

	latitude := fmt.Sprintf("%f", geocode.Response.View[0].Result[0].Location.DisplayPosition.Latitude)
	longitude := fmt.Sprintf("%f", geocode.Response.View[0].Result[0].Location.DisplayPosition.Longitude)

	weatherURL := "https://api.darksky.net/forecast/" + darkSkyAPIKey + "/" + latitude + "," + longitude + "?units=si"

	var weather weatherInfo

	resp, err = http.Get(weatherURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		log.Fatal(err)
	}

	result := fmt.Sprintf("Current temperature is %v\u00b0C. Although it feels like %v\u00b0C. It's mostly %v with humidity at %v.",
		weather.Currently.Temperature, weather.Currently.ApparentTemperature, weather.Currently.Summary, weather.Currently.Humidity)

	fmt.Println(result)
}
