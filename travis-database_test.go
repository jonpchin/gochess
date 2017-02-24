package gostuff

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icrowley/fake"
)

// Travis CI default MySQL username and pass is public information

func TestTravisConnect(t *testing.T) {

	// only run this test in Travis
	if isEnvironmentTravis() == false {
		return
	}

	//if database ping fails here that means connection is alive but database is missing
	if db.Ping() != nil {
		t.Fatal("Can't ping MySQL in Travis")
	}

	// registers a random person to the database
	var userInfo UserInfo
	userInfo.Username = fake.UserName()
	userInfo.Password = fake.Password(5, 32, true, true, false)

	userInfo.Email = fake.EmailAddress()
	userInfo.IpAddress = fake.IPv4()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := userInfo.Register(w, r)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)

	found := CheckUserNameInDb(userInfo.Username)
	if found == false {
		t.Fatal("Username was not found in the database")
	}
}
