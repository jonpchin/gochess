package webtest

import (
	"bufio"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sclevine/agouti"
)

//localhost testing
func TestLoginDev(t *testing.T) {

	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		t.Fatal("Failed to start Chrome Driver:", err)
	}
	page1, err := driver.NewPage(agouti.Browser("Chrome"))
	if err != nil {
		t.Fatal("Failed to open page:", err)
	}

	if err := page1.Navigate("https://localhost:443"); err != nil {
		t.Fatal("Failed to navigate index at localhost:", err)
	}

	if err := page1.Navigate("https://localhost/login"); err != nil {
		t.Fatal("Failed to navigate login at localhost:", err)
	}

	loginURL, err := page1.URL()
	if err != nil {
		t.Fatal("Failed to get page URL:", err)
	}

	expectedLoginURL := "https://localhost/login"
	if loginURL != expectedLoginURL {
		t.Fatal("Expected URL to be", expectedLoginURL, "but got", loginURL)
	}
	user1 := "can"
	err = page1.FindByID("user").Fill(user1)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}
	pass := readPass(user1)
	err = page1.FindByID("password").Fill(pass)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}

	err = page1.FindByID("login").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}

	time.Sleep(time.Second)
	if err := page1.Navigate("https://localhost/server/lobby"); err != nil {
		t.Fatal("Failed to navigate lobby at localhost:", err)
	}

	err = page1.FindByID("sendSeek").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}

	// start second browser
	page2, err := driver.NewPage(agouti.Browser("Chrome"))
	if err != nil {
		t.Fatal("Failed to open page:", err)
	}

	if err := page2.Navigate("https://localhost:443"); err != nil {
		t.Fatal("Failed to navigate index at localhost:", err)
	}

	if err := page2.Navigate("https://localhost/login"); err != nil {
		t.Fatal("Failed to navigate login at localhost:", err)
	}

	user2 := "ben"
	err = page2.FindByID("user").Fill(user2)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}
	pass = readPass(user2)
	err = page2.FindByID("password").Fill(pass)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}

	err = page2.FindByID("login").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}
	time.Sleep(time.Second)

	if err := page2.Navigate("https://localhost/server/lobby"); err != nil {
		t.Fatal("Failed to navigate lobby at localhost:", err)
	}

	err = page2.FindByID("sendSeek").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}
	time.Sleep(time.Second)
	var whitePlayer string
	page2.RunScript("return WhiteSide;", map[string]interface{}{}, &whitePlayer)
	var jsResult string

	if user1 == whitePlayer {
		page1.RunScript("sendMove('e2', 'e4');", map[string]interface{}{}, &jsResult)
		page2.RunScript("sendMove('c7', 'c5');", map[string]interface{}{}, &jsResult)
		page1.RunScript("sendMove('g1', 'f3');", map[string]interface{}{}, &jsResult)
		page1.RunScript("return board.fen();", map[string]interface{}{}, &jsResult)

		// check to make sure the position is what it should be
		if jsResult != "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R" {
			t.Error("board does not match user1")
		}

		// now try to resign the game
		err = page1.FindByID("resignButton").Click()
		if err != nil {
			t.Fatal("Couldn't resign user1:", err)
		}
		err = page1.ConfirmPopup()
		if err != nil {
			t.Fatal("Couldn't confirm resign popup user1:", err)
		}
		err = page1.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page2.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page1.FindByID("abortButton").Click()
		if err != nil {
			t.Fatal("Couldn't find abort button  user 1:", err)
		}
		err = page2.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page1.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 2:", err)
		}
		err = page1.FindByID("drawButton").Click()
		if err != nil {
			t.Fatal("Couldn't find draw button  user 1:", err)
		}
		err = page2.FindByID("drawButton").Click()
		if err != nil {
			t.Fatal("Couldn't find draw button  user 2:", err)
		}
		// TODO: Check if game really ended and check if the other player really won
		// Still need to test abort failure, abort sucess, draw, and checkmate

	} else if user2 == whitePlayer {
		page2.RunScript("sendMove('e2', 'e4');", map[string]interface{}{}, &jsResult)
		page1.RunScript("sendMove('c7', 'c5');", map[string]interface{}{}, &jsResult)
		page2.RunScript("sendMove('g1', 'f3');", map[string]interface{}{}, &jsResult)
		page2.RunScript("return board.fen();", map[string]interface{}{}, &jsResult)
		if jsResult != "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R" {
			t.Error("board does not match user2")
		}
		err = page2.FindByID("resignButton").Click()
		if err != nil {
			t.Fatal("Couldn't resign user2:", err)
		}
		err = page2.ConfirmPopup()
		if err != nil {
			t.Fatal("Couldn't confirm resign popup user2:", err)
		}
		err = page2.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 2:", err)
		}
		err = page1.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page1.FindByID("abortButton").Click()
		if err != nil {
			t.Fatal("Couldn't find abort button  user 1:", err)
		}
		err = page1.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page2.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 2:", err)
		}
		err = page1.FindByID("drawButton").Click()
		if err != nil {
			t.Fatal("Couldn't find draw button  user 1:", err)
		}
		err = page2.FindByID("drawButton").Click()
		if err != nil {
			t.Fatal("Couldn't find draw button  user 2:", err)
		}
	} else {
		// then navigate to chess page and try to terminate any possible games that are left over
		if err := page2.Navigate("https://localhost/chess/memberChess"); err != nil {
			t.Fatal("Failed to navigate login to chess page:", err)
		}
		err = page2.FindByID("abortButton").Click()
		if err != nil {
			t.Fatal("Couldn't find abort button  user 2:", err)
		}
		t.Fatal("No user matched as whitePlayer")
	}
	time.Sleep(time.Second)
	if err := driver.Stop(); err != nil {
		t.Error("Failed to close pages and stop WebDriver:", err)
	}
}

func TestLoginProduction(t *testing.T) {

	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		t.Fatal("Failed to start Chrome Driver:", err)
	}
	page1, err := driver.NewPage(agouti.Browser("Chrome"))
	if err != nil {
		t.Fatal("Failed to open page:", err)
	}

	if err := page1.Navigate("https://goplaychess.com:443"); err != nil {
		t.Fatal("Failed to navigate index:", err)
	}

	if err := page1.Navigate("https://goplaychess.com/login"); err != nil {
		t.Fatal("Failed to navigate login:", err)
	}

	loginURL, err := page1.URL()
	if err != nil {
		t.Fatal("Failed to get page URL:", err)
	}

	expectedLoginURL := "https://goplaychess.com/login"
	if loginURL != expectedLoginURL {
		t.Fatal("Expected URL to be", expectedLoginURL, "but got", loginURL)
	}
	user1 := "foo"
	err = page1.FindByID("user").Fill(user1)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}
	pass := readPass(user1)
	err = page1.FindByID("password").Fill(pass)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}

	err = page1.FindByID("login").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}
	time.Sleep(time.Second)
	if err := page1.Navigate("https://goplaychess.com/server/lobby"); err != nil {
		t.Fatal("Failed to navigate lobby:", err)
	}
	time.Sleep(time.Second)
	err = page1.FindByID("sendSeek").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}

	// start second browser
	page2, err := driver.NewPage(agouti.Browser("Chrome"))
	if err != nil {
		t.Fatal("Failed to open page:", err)
	}

	if err := page2.Navigate("https://goplaychess.com:443"); err != nil {
		t.Fatal("Failed to navigate index at localhost:", err)
	}

	if err := page2.Navigate("https://goplaychess.com/login"); err != nil {
		t.Fatal("Failed to navigate login at localhost:", err)
	}

	user2 := "Carl"
	err = page2.FindByID("user").Fill(user2)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}

	err = page2.FindByID("password").Fill(readPass(user2))
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}

	err = page2.FindByID("login").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}
	time.Sleep(time.Second)

	if err := page2.Navigate("https://goplaychess.com/server/lobby"); err != nil {
		t.Fatal("Failed to navigate lobby at localhost:", err)
	}
	time.Sleep(time.Second)
	err = page2.FindByID("sendSeek").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}

	time.Sleep(time.Second)
	var whitePlayer string
	page2.RunScript("return WhiteSide;", map[string]interface{}{}, &whitePlayer)
	var jsResult string
	time.Sleep(2 * time.Second)
	if user1 == whitePlayer {
		page1.RunScript("sendMove('e2', 'e4');", map[string]interface{}{}, &jsResult)
		page2.RunScript("sendMove('c7', 'c5');", map[string]interface{}{}, &jsResult)
		page1.RunScript("sendMove('g1', 'f3');", map[string]interface{}{}, &jsResult)
		time.Sleep(time.Second)
		page1.RunScript("return board.fen();", map[string]interface{}{}, &jsResult)

		// check to make sure the position is what it should be
		if jsResult != "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R" {
			t.Fatal("board does not match user1")
		}

		// now try to resign the game
		err = page1.FindByID("resignButton").Click()
		if err != nil {
			t.Fatal("Couldn't resign user1:", err)
		}
		err = page1.ConfirmPopup()
		if err != nil {
			t.Fatal("Couldn't confirm resign popup user1:", err)
		}
		err = page1.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page2.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page1.FindByID("abortButton").Click()
		if err != nil {
			t.Fatal("Couldn't find abort button  user 1:", err)
		}
		err = page2.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page1.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 2:", err)
		}
		err = page1.FindByID("drawButton").Click()
		if err != nil {
			t.Fatal("Couldn't find draw button  user 1:", err)
		}
		err = page2.FindByID("drawButton").Click()
		if err != nil {
			t.Fatal("Couldn't find draw button  user 2:", err)
		}
		// TODO: Check if game really ended and check if the other player really won
		// Still need to test abort failure, abort sucess, draw, and checkmate

	} else if user2 == whitePlayer {
		page2.RunScript("sendMove('e2', 'e4');", map[string]interface{}{}, &jsResult)
		page1.RunScript("sendMove('c7', 'c5');", map[string]interface{}{}, &jsResult)
		page2.RunScript("sendMove('g1', 'f3');", map[string]interface{}{}, &jsResult)
		time.Sleep(time.Second)
		page2.RunScript("return board.fen();", map[string]interface{}{}, &jsResult)

		if jsResult != "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R" {
			t.Error("board does not match user2")
		}
		err = page2.FindByID("resignButton").Click()
		if err != nil {
			t.Fatal("Couldn't resign user2:", err)
		}
		err = page2.ConfirmPopup()
		if err != nil {
			t.Fatal("Couldn't confirm resign popup user2:", err)
		}
		err = page2.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 2:", err)
		}
		err = page1.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page1.FindByID("abortButton").Click()
		if err != nil {
			t.Fatal("Couldn't find abort button  user 1:", err)
		}
		err = page1.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 1:", err)
		}
		err = page2.FindByID("rematchButton").Click()
		if err != nil {
			t.Fatal("Couldn't find rematch button  user 2:", err)
		}
		err = page1.FindByID("drawButton").Click()
		if err != nil {
			t.Fatal("Couldn't find draw button  user 1:", err)
		}
		err = page2.FindByID("drawButton").Click()
		if err != nil {
			t.Fatal("Couldn't find draw button  user 2:", err)
		}

	} else {
		// then navigate to chess page and try to terminate any possible games that are left over
		if err := page2.Navigate("https://goplaychess.com/chess/memberChess"); err != nil {
			t.Fatal("Failed to navigate login to chess page:", err)
		}
		err = page2.FindByID("abortButton").Click()
		if err != nil {
			t.Fatal("Couldn't find abort button  user 2:", err)
		}
		t.Fatal("No user matched as whitePlayer")
	}
	time.Sleep(time.Second)
	if err := driver.Stop(); err != nil {
		t.Error("Failed to close pages and stop WebDriver:", err)
	}
}

// returns pass of user's account
func readPass(user string) string {
	config, err := os.Open("data/" + user + ".txt")
	defer config.Close()
	if err != nil {
		log.Println("web_Test.go readAccount 1 ", err)
	}
	scanner := bufio.NewScanner(config)
	scanner.Scan()

	return scanner.Text()
}