package gostuff

import(
	"sync"
	"os"
	"log"
	"fmt"
	"strconv"
)

const (
	total = 10 //max number of queries returned for high score board
)

type ScoreBoard struct{
	Bullet   [total]TopRating
	Blitz    [total]TopRating
	Standard [total]TopRating
	Recent   [total]RecentPlayer //ten most recently registered players
}

//used for bullet, blitz and standard
type TopRating struct{
	Name string
	Rating int
	Index int
}

type RecentPlayer struct{
	Name string
	Date string
	Index int
}

var LeaderBoard = struct {
	sync.RWMutex
	Scores ScoreBoard
}{}

//fetches top ten bullet, blitz and standard ratings as well as the most recent 10 registered players
func UpdateHighScore(){
	
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)
	
	var score ScoreBoard
	
	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN! scoreboard.go @updateHighScore() ping")
		return
	}
	var query string
	
	query = "SELECT username, bullet FROM rating order by bullet DESC limit " + strconv.Itoa(total)
	
	rows, err := db.Query(query)

	if err != nil {
		log.Println("scoreBoard.go updateHighScore() 1 ", err)
		return
	}
	defer rows.Close()
	i := 0
	
    for rows.Next() {
		

        if err := rows.Scan(&score.Bullet[i].Name, &score.Bullet[i].Rating); err != nil {
                fmt.Println("scoreBoard.go updateHighScore() 2 ", err)
        }
//		fmt.Printf("name is %s rating is %d\n", score.Bullet[i].Name, score.Bullet[i].Rating)
		score.Bullet[i].Index = i+1
		i++
        
    }
	i=0
	
	query = "SELECT username, blitz FROM rating order by blitz DESC limit " + strconv.Itoa(total)
	
	rows, err = db.Query(query)

	if err != nil {
		fmt.Println("scoreBoard.go updateHighScore() 2 ", err)
		return
	}

    for rows.Next() {

		
        if err := rows.Scan(&score.Blitz[i].Name, &score.Blitz[i].Rating); err != nil {
                fmt.Println("scoreBoard.go updateHighScore() 3 ", err)
        }
//		fmt.Printf("name is %s rating is %d\n", score.Blitz[i].Name, score.Blitz[i].Rating)
		score.Blitz[i].Index = i+1
		i++
        
    }
	i=0
	
	query = "SELECT username, standard FROM rating order by standard DESC limit " + strconv.Itoa(total)
	
	rows, err = db.Query(query)

	if err != nil {
		fmt.Println("scoreBoard.go updateHighScore() 3 ", err)
		return
	}

    for rows.Next() {

		
        if err := rows.Scan(&score.Standard[i].Name, &score.Standard[i].Rating); err != nil {
                fmt.Println("scoreBoard.go updateHighScore() 4 ", err)
        }
//		fmt.Printf("name is %s rating is %d\n", score.Standard[i].Name, score.Standard[i].Rating)
		score.Standard[i].Index = i+1
		i++
        
    }
	i=0
	
	query = "SELECT username, date FROM userinfo order by date desc, time desc limit " + strconv.Itoa(total)
	
	rows, err = db.Query(query)

	if err != nil {
		fmt.Println("scoreBoard.go updateHighScore() 3 ", err)
		return
	}

    for rows.Next() {
		
        if err := rows.Scan(&score.Recent[i].Name, &score.Recent[i].Date); err != nil {
                fmt.Println("scoreBoard.go updateHighScore() 4 ", err)
        }
//		fmt.Printf("name is %s date is %s\n", score.Recent[i].Name, score.Recent[i].Date)
		score.Recent[i].Index = i+1
		i++
        
    }
	//secure mutex lock before modifying global leaderboard in memory
	LeaderBoard.Lock()
	LeaderBoard.Scores = score
	LeaderBoard.Unlock()
}