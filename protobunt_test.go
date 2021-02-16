package protobunt

import (
	pb "protobunt/proto"
	"testing"
	"time"
)

import "github.com/tidwall/buntdb"


func TestServerClient(t *testing.T) {
	db, _ := buntdb.Open(":memory:")
	go StartBuntServer("127.0.0.1", "8080", db)
	time.Sleep(2 * time.Second) // Pause to be sure that service is up even on weak hardware
	client, ctx, cancel := CreateBuntClient("127.0.0.1", "8080")
	defer cancel()

	testKey := "Alice"
	testValue1 := "Bob"
	testValue2 := "Charlie"

	r1 := pb.ViewRequest{Key: testKey, Action: GET}
	v, _ := client.View(ctx, &r1)
	res1 := v.GetVal()
	if res1 != "" {
		t.Error("FAIL 1")
	}

	r2 := pb.UpdateRequest{Key: testKey, Value: testValue1, Action: SET}
	_, err := client.Update(ctx, &r2)
	if err != nil {
		t.Error("FAIL 2" + err.Error())
	}

	r3 := pb.ViewRequest{Key: testKey, Action: GET}
	initValue, _ := client.View(ctx, &r3)
	res3 := initValue.GetVal()
	if res3 != testValue1 {
		t.Error("FAIL 3")
	}

	r4 := pb.UpdateRequest{Key: testKey, Value: testValue2, Action: SET}
	updateValue, _ := client.Update(ctx, &r4)
	res4 := updateValue.GetReplaced()
	if !res4 {
		t.Error("FAIL 4")
	}

	r5 := pb.ViewRequest{Key: testKey, Action: GET}
	updatedValue, _ := client.View(ctx, &r5)
	res5 := updatedValue.GetVal()
	if res5 != testValue2 {
		t.Error("FAIL 5")
	}

	r6 := pb.UpdateRequest{Key: testKey, Action: DELETE}
	deleteKey, _ := client.Update(ctx, &r6)
	res6 := deleteKey.GetPreviousValue()
	if res6 != testValue2 {
		t.Error("FAIL 6")
	}

	r7 := pb.ViewRequest{Key: testKey, Action: GET}
	deletedValue, _ := client.View(ctx, &r7)
	res7 := deletedValue.GetVal()
	if res7 != "" {
		t.Error("FAIL 7")
	}
}