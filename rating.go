//Author:  Josh Hoak aka Kashomon
package gostuff

import (
	"fmt"
)

// Represents a player's rating and the confidence in a player's rating.
type Rating struct {
	Rating     float64 // Player's rating. Usually starts off at 1500.
	Deviation  float64 // Confidence/uncertainty in a player's rating
	Volatility float64 // Measures erratic performances
}

// Creates a default Rating using:
// 	Rating     = DefaultRat
// 	Deviation  = DefaultDev
// 	Volatility = DefaultVol
func DefaultRating() *Rating {
	return &Rating{DefaultRat, DefaultDev, DefaultVol}
}

// Creates a new custom Rating.
func NewRating(r, rd, s float64) *Rating {
	return &Rating{r, rd, s}
}

// Creates a new rating, converted from Glicko1 scaling to Glicko2 scaling.
// This assumes the starting rating value is 1500.
func (r *Rating) ToGlicko2() *Rating {
	return NewRating(
		(r.Rating-DefaultRat)/glicko2Scale,
		(r.Deviation)/glicko2Scale,
		r.Volatility)
}

// Creates a new rating, converted from Glicko2 scaling to Glicko1 scaling.
// This assumes the starting rating value is 1500.
func (r *Rating) FromGlicko2() *Rating {
	return NewRating(
		r.Rating*glicko2Scale+DefaultRat,
		r.Deviation*glicko2Scale,
		r.Volatility)
}

func (r *Rating) String() string {
	return fmt.Sprintf("{Rating[%.3f] Deviation[%.3f] Volatility[%.3f]}",
		r.Rating, r.Deviation, r.Volatility)
}

// Ensure that some other Rating is equal to this rating, given some epsilon. In
// other words, find the error between this rating's values and the other
// rating's values and make sure it's less than epsilon in absolute value.
func (r *Rating) MostlyEquals(o *Rating, epsilon float64) bool {
	return floatsMostlyEqual(r.Rating, o.Rating, epsilon) &&
		floatsMostlyEqual(r.Deviation, o.Deviation, epsilon) &&
		floatsMostlyEqual(r.Volatility, o.Volatility, epsilon)
}

// Create a duplicate rating with the same values.
func (r *Rating) Copy() *Rating {
	return &Rating{
		r.Rating,
		r.Deviation,
		r.Volatility,
	}
}
