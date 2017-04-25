package plot

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jonpchin/gochess/gostuff"
	chart "github.com/wcharczuk/go-chart"
)

var (
	UseGoogleCharts = false // Toggle to enable Google Charts or DrawChart()
)

func SetupCharts() {
	if UseGoogleCharts {
		gostuff.ReplaceString("var useGoogleCharts = false;", "var useGoogleCharts = true;",
			"js/profile.js")
	} else {
		gostuff.ReplaceString("var useGoogleCharts = true;", "var useGoogleCharts = false;",
			"js/profile.js")
	}
}

func DrawChart(w http.ResponseWriter, r *http.Request) {

	valid := gostuff.ValidateCredentials(w, r)

	if valid == false {
		return
	}

	const (
		bullet         = "bullet"
		blitz          = "blitz"
		standard       = "standard"
		correspondence = "correspondence"
	)

	username, _ := r.Cookie("username")
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	ratingList := []string{bullet, blitz, standard, correspondence}
	timeFormat := "20060102150405"

	// Index of these double arrays follows same order as ratingList
	var allRatingDates [][]time.Time
	var allRatings [][]float64

	for _, value := range ratingList {

		ratingHistory, pass, err := gostuff.GetRatingHistory(username.Value, value)
		if pass == false {
			log.Println(err)
			return
		}

		var ratingMemory []gostuff.RatingDate
		var ratingDates []time.Time
		var ratings []float64

		if ratingHistory != "" {
			if err := json.Unmarshal([]byte(ratingHistory), &ratingMemory); err != nil {
				log.Println("Just receieved a message I couldn't decode:", ratingHistory, err)
				return
			}

			for _, value := range ratingMemory {
				dateTime, err := time.Parse(timeFormat, value.DateTime)
				if err != nil {
					log.Println(err)
				}
				ratingDates = append(ratingDates, dateTime)
				ratings = append(ratings, value.Rating)
			}

			allRatingDates = append(allRatingDates, ratingDates)
			allRatings = append(allRatings, ratings)

		} else {
			// Then the player has no rating in this gametype
			allRatingDates = append(allRatingDates, nil)
			allRatings = append(allRatings, nil)
		}
	}

	graph := chart.Chart{
		Title: "Rating History",
		XAxis: chart.XAxis{
			Name:      "Time",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      "Rating",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    bullet,
				XValues: allRatingDates[0],
				YValues: allRatings[0],
			},
			chart.TimeSeries{
				Name:    blitz,
				XValues: allRatingDates[1],
				YValues: allRatings[1],
			},
			chart.TimeSeries{
				Name:    standard,
				XValues: allRatingDates[2],
				YValues: allRatings[2],
			},
			chart.TimeSeries{
				Name:    correspondence,
				XValues: allRatingDates[3],
				YValues: allRatings[3],
			},
		},
	}
	//w.Header().Set("Content-Type", "image/png")
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	output, err := os.Create("img/plots/" + username.Value + ".png")
	if err != nil {
		log.Println(err)
	}
	defer output.Close()
	fileWriter := bufio.NewWriter(output)
	graph.Render(chart.PNG, fileWriter)
	w.Write([]byte("pass"))
}
