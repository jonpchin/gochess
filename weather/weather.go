package weather

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jonpchin/gochess/gostuff"
	geo "github.com/kellydunn/golang-geo"
)

// Contains information for the weather at this weather station
type WeatherAtStation struct {
	Weather_stn_id   int     // Weather station ID
	Weather_stn_name string  // Weather station name
	Weather_stn_lat  float64 // Latitude
	Weather_stn_long float64 // Longtitude
}

type RawWeatherData struct {
	Next  Reference
	Items []WeatherAtStation
}

type Reference struct {
	Ref string // Contains url to next page of weather stations if it can't fit on the current page
}

type weatherAtStations []WeatherAtStation

func FetchWeather() {
	const (
		weatherStations = "https://apex.oracle.com/pls/apex/raspberrypi/weatherstation/getallstations"
	)

	problems, _ := os.OpenFile("weather/data/weather_stations.json", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	client := gostuff.TimeOutHttp(5)
	response, err := client.Get(weatherStations)
	if response == nil {
		log.Println("fetchWeather URL time out for ", weatherStations)
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var weatherData RawWeatherData
	err = json.Unmarshal(responseData, &weatherData)
	if err != nil {
		log.Println("Just receieved a message I couldn't decode:", string(responseData), err)
	}

	geo.SetGoogleAPIKey(getGoogleMapsApiKey())

	for _, value := range weatherData.Items {
		point := geo.NewPoint(value.Weather_stn_lat, value.Weather_stn_long)
		var googleGeocoder geo.GoogleGeocoder
		result, err := googleGeocoder.ReverseGeocode(point)
		if err != nil {
			fmt.Println("Could not reverse geocode", err)
		} else {
			fmt.Println(result)
		}
	}
}

func (allWeather weatherAtStations) processWeatherStations() {

}

// returns news API key
func getGoogleMapsApiKey() string {

	// The file path where the news API key is
	const apiKeyPath = "weather/keys/google_maps_api_key.txt"
	config, err := os.Open(apiKeyPath)
	defer config.Close()

	if err != nil {
		fmt.Println("news.go getApiKey 1", err)
	}

	scanner := bufio.NewScanner(config)

	// Google maps API Key
	scanner.Scan()
	return scanner.Text()
}
