package protobunt

import (
	"testing"
	"time"
)

import "github.com/tidwall/buntdb"


func TestServerClient(t *testing.T) {
	db, _ := buntdb.Open(":memory:")
	go StartBuntServer("127.0.0.1", "8080", db)
	time.Sleep(1 * time.Second)
	cli := CreateBuntClient("127.0.0.1", "8080")

	testKey := "Alice"
	testValue1 := "Bob"
	testValue2 := "Charlie"

	emptyValueOnGet := cli.View(GET, testKey)
	if emptyValueOnGet != "" {
		t.Error("FAIL")
	}
	emptyValue := cli.Update(SET, testKey, testValue1)
	if emptyValue != "" {
		t.Error("FAIL")
	}

	initValue := cli.View(GET, testKey)
	if initValue != testValue1 {
		t.Error("FAIL")
	}

	cli.Update(SET, testKey, testValue2)

	updatedValue := cli.View(GET, testKey)

	if updatedValue != testValue2 {
		t.Error("FAIL")
	}

	cli.Update(DELETE, testKey, "")

	deletedValue := cli.View(GET, testKey)

	if deletedValue != "" {
		t.Error("FAIL")
	}

}