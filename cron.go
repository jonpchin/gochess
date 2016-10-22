package gostuff

import (
	"fmt"
	"log"
	"os"

	"github.com/robfig/cron"
)

func StartCron() {
	c := cron.New()
	c.AddFunc("@daily", updateRD)
	//	c.AddFunc("@weekly", ExportDatabase)
	c.AddFunc("@hourly", UpdateHighScore)
	c.Start()
}

func updateRD() { //increase rating RD by one in database if its less then 250, default is 350

	problems, err := os.OpenFile("logs/mainLog.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		fmt.Println("DATABASE DOWN! @updateRating() ping")
		return
	}

	stmt, err := db.Prepare("update rating set bulletRD=bulletRD+1 where bulletRD < 250")
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Println(err)
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Printf("%d rows were updated in bullet ratingRD table.\n", affect)

	stmt, err = db.Prepare("update rating set blitzRD=blitzRD+1 where blitzRD < 250")

	if err != nil {
		log.Println(err)
		return
	}

	res, err = stmt.Exec()
	if err != nil {
		log.Println(err)
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Printf("%d rows were updated in blitz ratingRD table.\n", affect)

	stmt, err = db.Prepare("update rating set standardRD=standardRD+1 where standardRD < 250")
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return
	}

	res, err = stmt.Exec()
	if err != nil {
		log.Println(err)
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Printf("%d rows were updated in standard ratingRD table.\n", affect)
}

//remove games older then 30 days to clean up profile page, activated only on server startup
func RemoveOldGames(days string) {

	problems, err := os.OpenFile("logs/mainLog.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return
	}

	stmt, err := db.Prepare("DELETE FROM games WHERE date < NOW() - INTERVAL " + days + " DAY")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Println(err)
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Printf("%d rows were deleted from games because they were older then "+days+" days.\n", affect)
}

// Remove old entries in activate table in database
func RemoveOldActivate(days string) {
	problems, err := os.OpenFile("logs/mainLog.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return
	}

	stmt, err := db.Prepare("DELETE FROM activate WHERE expire < NOW() - INTERVAL " + days + " DAY")
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Println(err)
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Printf("%d rows were deleted from activate table because they were older then "+days+" days.\n", affect)
}

// Remove old entries in the forgot table in the database
// If the entry is older then the days parameter then it will be deleted
func RemoveOldForgot(days string) {
	problems, err := os.OpenFile("logs/mainLog.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return
	}

	stmt, err := db.Prepare("DELETE FROM forgot WHERE expire < NOW() - INTERVAL " + days + " DAY")
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Println(err)
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Printf("%d rows were deleted from forgot table because they were older then "+days+" days.\n", affect)
}
