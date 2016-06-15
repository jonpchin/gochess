package gostuff

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"os"
)

func StartCron() {
	c := cron.New()
	c.AddFunc("@daily", updateRD)
	c.AddFunc("@hourly", UpdateHighScore)
	c.Start()
}

func updateRD() { //increase rating RD by one in database if its less then 250, default is 350

	problems, err := os.OpenFile("logs/mainLog.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//check if database connection is open
	if db.Ping() != nil {
		fmt.Println("DATABASE DOWN! @updateRating() ping")
		return
	}

	stmt, err := db.Prepare("update rating set bulletRD=bulletRD+1 where bulletRD < 250")

	if err != nil {
		log.Println("cron.go 1", err)
		return
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Println("cron.go 2 ", err)
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println("cron.go 3 ", err)
	}
	log.Printf("%d rows were updated in bullet ratingRD table.\n", affect)

	stmt, err = db.Prepare("update rating set blitzRD=blitzRD+1 where blitzRD < 250")

	if err != nil {
		log.Println("cron.go 4 ", err)
		return
	}

	res, err = stmt.Exec()
	if err != nil {
		log.Println("cron.go 5 ", err)
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Println("cron.go 6 ", err)
	}
	log.Printf("%d rows were updated in blitz ratingRD table.\n", affect)

	stmt, err = db.Prepare("update rating set standardRD=standardRD+1 where standardRD < 250")

	if err != nil {
		log.Println("cron.go 7 ", err)
		return
	}

	res, err = stmt.Exec()
	if err != nil {
		log.Println("cron.go 7 ", err)
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Println("cron.go 8 ", err)
	}
	log.Printf("%d rows were updated in standard ratingRD table.\n", affect)
}

//remove games older then 30 days to clean up profile page, activated only on server startup
func RemoveOldGames() {
	//DELETE FROM games WHERE date < NOW() - INTERVAL 100 DAY;

	problems, err := os.OpenFile("logs/mainLog.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN! @RemoveOldGames() ping")
		return
	}

	stmt, err := db.Prepare("DELETE FROM games WHERE date < NOW() - INTERVAL 30 DAY")

	if err != nil {
		log.Println("cron.go 9 ", err)
		return
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Println("cron.go 10 ", err)
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println("cron.go 11 ", err)
	}
	log.Printf("%d rows were deleted from games because they were older then 30 days.\n", affect)
}
