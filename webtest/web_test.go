package webtest

import (
	"bufio"
	"log"
	"os"
	"testing"

	"github.com/sclevine/agouti"
)

//localhost testing
func TestLoginDev(t *testing.T) {

	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		t.Fatal("Failed to start Chrome Driver:", err)
	}
	page, err := driver.NewPage(agouti.Browser("Chrome"))
	if err != nil {
		t.Fatal("Failed to open page:", err)
	}

	if err := page.Navigate("https://localhost:443"); err != nil {
		t.Fatal("Failed to navigate:", err)
	}

	if err := page.Navigate("https://localhost/login"); err != nil {
		t.Fatal("Failed to navigate:", err)
	}

	loginURL, err := page.URL()
	if err != nil {
		t.Fatal("Failed to get page URL:", err)
	}

	expectedLoginURL := "https://localhost/login"
	if loginURL != expectedLoginURL {
		t.Fatal("Expected URL to be", expectedLoginURL, "but got", loginURL)
	}
	user := "can"
	err = page.FindByID("user").Fill(user)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}
	pass := readPass(user)
	err = page.FindByID("password").Fill(pass)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}

	err = page.FindByID("login").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}

	//if err := driver.Stop(); err != nil {
	//	t.Fatal("Failed to close pages and stop WebDriver:", err)
	//}
}

func TestLoginProduction(t *testing.T) {

	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		t.Fatal("Failed to start Chrome Driver:", err)
	}
	page, err := driver.NewPage(agouti.Browser("Chrome"))
	if err != nil {
		t.Fatal("Failed to open page:", err)
	}

	if err := page.Navigate("https://goplaychess.com:443"); err != nil {
		t.Fatal("Failed to navigate:", err)
	}

	if err := page.Navigate("https://goplaychess.com/login"); err != nil {
		t.Fatal("Failed to navigate:", err)
	}

	loginURL, err := page.URL()
	if err != nil {
		t.Fatal("Failed to get page URL:", err)
	}

	expectedLoginURL := "https://goplaychess.com/login"
	if loginURL != expectedLoginURL {
		t.Fatal("Expected URL to be", expectedLoginURL, "but got", loginURL)
	}
	user := "foo"
	err = page.FindByID("user").Fill(user)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}
	pass := readPass(user)
	err = page.FindByID("password").Fill(pass)
	if err != nil {
		t.Fatal("Couldn't fill login info:", err)
	}

	err = page.FindByID("login").Click()
	if err != nil {
		t.Fatal("Couldn't submit:", err)
	}

	//if err := driver.Stop(); err != nil {
	//	t.Fatal("Failed to close pages and stop WebDriver:", err)
	//}
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
