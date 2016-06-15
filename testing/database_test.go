package testing

import(
	"testing"
	"github.com/jonpchin/GoChess/gostuff"
)

func TestDbsetup(t *testing.T){
	
	var connect bool
	connect = gostuff.DbSetup("../secret/config.txt")
	if connect == false{
		t.Error("Database failed to connect database_test.go TestingDbsetup()")
	}

}