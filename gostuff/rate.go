//Author:  Josh Hoak aka Kashomon
package gostuff

import (
	"fmt"
	"math"

	"golang.org/x/net/websocket"
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
func ComputeRating(name string, gameID int, gameType string, result float64) {

	var bullet, blitz, standard, correspondence, bulletRD, blitzRD, standardRD, correspondenceRD float64
	var oBullet, oBlitz, oStandard, oCorrespondence, oBulletRD, oBlitzRD, oStandardRD, oCorrespondenceRD float64

	//update player's rating and notify them of rating change, also determine player color to assign correct rating
	if All.Games[gameID].WhitePlayer == name {
		_, bullet, blitz, standard, correspondence, bulletRD, blitzRD, standardRD,
			correspondenceRD = GetRatingAndRD(name)
		_, oBullet, oBlitz, oStandard, oCorrespondence, oBulletRD, oBlitzRD,
			oStandardRD, oCorrespondenceRD = GetRatingAndRD(PrivateChat[name])
	} else {
		_, bullet, blitz, standard, correspondence, bulletRD, blitzRD, standardRD,
			correspondenceRD = GetRatingAndRD(PrivateChat[name])
		_, oBullet, oBlitz, oStandard, oCorrespondence, oBulletRD, oBlitzRD,
			oStandardRD, oCorrespondenceRD = GetRatingAndRD(name)
	}

	var whiteRating float64
	var blackRating float64
	var whiteRD float64
	var blackRD float64
	const (
		bulletString         = "bullet"
		blitzString          = "blitz"
		standardString       = "standard"
		correspondenceString = "correspondence"
	)

	if gameType == bulletString {

		whiteRating, whiteRD = grabRating(bullet, bulletRD, oBullet, oBulletRD, result)
		blackRating, blackRD = grabRating(oBullet, oBulletRD, bullet, bulletRD, 1.0-result)

		//updates database with players new rating and RD
		if All.Games[gameID].WhitePlayer == name {
			updateRating(bulletString, name, whiteRating, whiteRD, PrivateChat[name], blackRating, blackRD)
			updateRatingHistory(name, bulletString, whiteRating)
			updateRatingHistory(PrivateChat[name], bulletString, blackRating)
		} else {
			updateRating(bulletString, PrivateChat[name], whiteRating, whiteRD, name, blackRating, blackRD)
			updateRatingHistory(PrivateChat[name], bulletString, whiteRating)
			updateRatingHistory(name, bulletString, blackRating)
		}

	} else if gameType == blitzString {

		whiteRating, whiteRD = grabRating(blitz, blitzRD, oBlitz, oBlitzRD, result)
		blackRating, blackRD = grabRating(oBlitz, oBlitzRD, blitz, blitzRD, 1.0-result)

		//updates both players rating
		if All.Games[gameID].WhitePlayer == name {
			updateRating(blitzString, name, whiteRating, whiteRD, PrivateChat[name], blackRating, blackRD)
			updateRatingHistory(name, blitzString, whiteRating)
			updateRatingHistory(PrivateChat[name], blitzString, blackRating)
		} else {
			updateRating(blitzString, PrivateChat[name], whiteRating, whiteRD, name, blackRating, blackRD)
			updateRatingHistory(PrivateChat[name], blitzString, whiteRating)
			updateRatingHistory(name, blitzString, blackRating)
		}

	} else if gameType == standardString {

		whiteRating, whiteRD = grabRating(standard, standardRD, oStandard, oStandardRD, result)
		blackRating, blackRD = grabRating(oStandard, oStandardRD, standard, standardRD, 1.0-result)
		//updates database with players new rating and RD
		if All.Games[gameID].WhitePlayer == name {
			updateRating(standardString, name, whiteRating, whiteRD, PrivateChat[name], blackRating, blackRD)
			updateRatingHistory(name, standardString, whiteRating)
			updateRatingHistory(PrivateChat[name], standardString, blackRating)
		} else {
			updateRating(standardString, PrivateChat[name], whiteRating, whiteRD, name, blackRating, blackRD)
			updateRatingHistory(PrivateChat[name], standardString, whiteRating)
			updateRatingHistory(name, standardString, blackRating)
		}
	} else if gameType == correspondenceString {
		whiteRating, whiteRD = grabRating(correspondence, correspondenceRD, oCorrespondence, oCorrespondenceRD, result)
		blackRating, blackRD = grabRating(oCorrespondence, oCorrespondenceRD, correspondence, correspondenceRD, 1.0-result)
		//updates database with players new rating and RD
		if All.Games[gameID].WhitePlayer == name {
			updateRating(correspondenceString, name, whiteRating, whiteRD, PrivateChat[name], blackRating, blackRD)
			updateRatingHistory(name, correspondenceString, whiteRating)
			updateRatingHistory(PrivateChat[name], correspondenceString, blackRating)
		} else {
			updateRating(correspondenceString, PrivateChat[name], whiteRating, whiteRD, name, blackRating, blackRD)
			updateRatingHistory(PrivateChat[name], correspondenceString, whiteRating)
			updateRatingHistory(name, correspondenceString, blackRating)
		}
	} else {
		fmt.Println("Not a valid game type rate.go 1")
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
