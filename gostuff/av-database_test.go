package gostuff

import (
	"testing"
)

// App Veyor default MySQL username and pass is public information
func TestAppVeyorConnect(t *testing.T) {

	// only run this test in AppVeyor
	if IsEnvironmentAppVeyor() == false {
		return
	}

	//if database ping fails here that means connection is alive but database is missing
	if db.Ping() != nil {
		t.Fatal("Can't ping MySQL in App Veyor")
	}
}
