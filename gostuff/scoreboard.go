package gostuff

import (
	"log"
	"os"
	"strconv"
	"sync"
)

const (
	total = 10 //max number of queries returned for high score board
)

type ScoreBoard struct {
	Bullet         [total]TopRating
	Blitz          [total]TopRating
	Standard       [total]TopRating
	Correspondence [total]TopRating
	Recent         [total]RecentPlayer //ten most recently registered players
}

//used for bullet, blitz and standard
type TopRating struct {
	Name   string
	Rating int
	Index  int
}

type RecentPlayer struct {
	Name  string
	Date  string
	Index int
}

var LeaderBoard = struct {
	sync.RWMutex
	Scores ScoreBoard
}{}

//fetches top ten bullet, blitz and standard ratings as well as the most recent 10 registered players
func UpdateHighScore() {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	var score ScoreBoard

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return
	}
	var query string

	query = "SELECT username, bullet FROM rating order by bullet DESC limit " + strconv.Itoa(total)

	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	i := 0

	for rows.Next() {

		if err := rows.Scan(&score.Bullet[i].Name, &score.Bullet[i].Rating); err != nil {
			log.Println(err)
		}
		//fmt.Printf("name is %s rating is %d\n", score.Bullet[i].Name, score.Bullet[i].Rating)
		score.Bullet[i].Index = i + 1
		i++
	}
	i = 0

	query = "SELECT username, blitz FROM rating order by blitz DESC limit " + strconv.Itoa(total)

	rows, err = db.Query(query)

	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&score.Blitz[i].Name, &score.Blitz[i].Rating); err != nil {
			log.Println(err)
		}
		score.Blitz[i].Index = i + 1
		i++
	}
	i = 0

	query = "SELECT username, standard FROM rating order by standard DESC limit " + strconv.Itoa(total)

	rows, err = db.Query(query)

	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&score.Standard[i].Name, &score.Standard[i].Rating); err != nil {
			log.Println(err)
		}
		score.Standard[i].Index = i + 1
		i++
	}
	i = 0

	query = "SELECT username, date FROM userinfo order by date desc, time desc limit " + strconv.Itoa(total)

	rows, err = db.Query(query)

	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&score.Recent[i].Name, &score.Recent[i].Date); err != nil {
			log.Println(err)
		}
		score.Recent[i].Index = i + 1
		i++
	}
	//secure mutex lock before modifying global leaderboard in memory
	LeaderBoard.Lock()
	LeaderBoard.Scores = score
	LeaderBoard.Unlock()
}
