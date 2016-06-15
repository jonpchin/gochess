//Author:  Josh Hoak aka Kashomon
package gostuff

import (
	"fmt"
	"golang.org/x/net/websocket"
	"math"
)

//fetches player's new rating by passing both player's rating and their deviation and game result and returns their rating and deviation
func grabRating(pRating float64, pDeviation float64, oRating float64, oDeviation float64, results float64) (float64, float64) {

	player := &Rating{pRating, pDeviation, 0.06}
	opponents := &Rating{oRating, oDeviation, DefaultVol}

	newRating, _ := CalculateRating(player, opponents, results)
	return Round(newRating.Rating), RoundPlus(newRating.Deviation, 4)
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}

//	fmt.Println(Round(123.4999))
//	fmt.Println(RoundPlus(123.558, 2))

//computes the rating for one player and the other player and updates the database and notifies both players, result can be white, black or draw
func ComputeRating(name string, gameID int16, gameType string, result float64) {

	var bullet, blitz, standard, bulletRD, blitzRD, standardRD float64
	var oBullet, oBlitz, oStandard, oBulletRD, oBlitzRD, oStandardRD float64

	//update player's rating and notify them of rating change, also determine player color to assign correct rating
	if All.Games[gameID].WhitePlayer == name {
		_, bullet, blitz, standard, bulletRD, blitzRD, standardRD = GetRatingAndRD(name)
		_, oBullet, oBlitz, oStandard, oBulletRD, oBlitzRD, oStandardRD = GetRatingAndRD(PrivateChat[name])
	} else {
		_, bullet, blitz, standard, bulletRD, blitzRD, standardRD = GetRatingAndRD(PrivateChat[name])
		_, oBullet, oBlitz, oStandard, oBulletRD, oBlitzRD, oStandardRD = GetRatingAndRD(name)
	}

	var whiteRating float64
	var blackRating float64
	var whiteRD float64
	var blackRD float64
	if gameType == "bullet" {

		whiteRating, whiteRD = grabRating(bullet, bulletRD, oBullet, oBulletRD, result)
		blackRating, blackRD = grabRating(oBullet, oBulletRD, bullet, bulletRD, 1.0-result)
		//updates database with players new rating and RD

		if All.Games[gameID].WhitePlayer == name {
			updateRating("bullet", name, whiteRating, whiteRD, PrivateChat[name], blackRating, blackRD)
		} else {
			updateRating("bullet", PrivateChat[name], whiteRating, whiteRD, name, blackRating, blackRD)
		}

	} else if gameType == "blitz" {

		whiteRating, whiteRD = grabRating(blitz, blitzRD, oBlitz, oBlitzRD, result)
		blackRating, blackRD = grabRating(oBlitz, oBlitzRD, blitz, blitzRD, 1.0-result)
		//updates database with players new rating and RD
		if All.Games[gameID].WhitePlayer == name {
			updateRating("blitz", name, whiteRating, whiteRD, PrivateChat[name], blackRating, blackRD)
		} else {
			updateRating("blitz", PrivateChat[name], whiteRating, whiteRD, name, blackRating, blackRD)
		}

	} else if gameType == "standard" {

		whiteRating, whiteRD = grabRating(standard, standardRD, oStandard, oStandardRD, result)
		blackRating, blackRD = grabRating(oStandard, oStandardRD, standard, standardRD, 1.0-result)
		//updates database with players new rating and RD
		if All.Games[gameID].WhitePlayer == name {
			updateRating("standard", name, whiteRating, whiteRD, PrivateChat[name], blackRating, blackRD)
		} else {
			updateRating("standard", PrivateChat[name], whiteRating, whiteRD, name, blackRating, blackRD)
		}
	} else {
		fmt.Println("Not a valid game type")
	}

	var r Nrating
	r.Type = "rating"
	r.WhiteRating = whiteRating
	r.BlackRating = blackRating

	if _, ok := Active.Clients[name]; ok { // send data if other guy is still connected
		websocket.JSON.Send(Active.Clients[name], &r)
	}

	if _, ok := Active.Clients[PrivateChat[name]]; ok { // send data if other guy is still connected
		websocket.JSON.Send(Active.Clients[PrivateChat[name]], &r)
	}

}
