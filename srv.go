package protobunt

import (
	"bufio"
	"errors"
	"github.com/tidwall/buntdb"
	"log"
	"net"
	"strings"
)


var (
	BadRequest = errors.New("bad request")
)

type BuntServer struct {
	db *buntdb.DB
	protobuntVersion string
}

type Action func(fn func(tx *buntdb.Tx) error) error

func (buntServer *BuntServer) processRequest(raw string) string{
	var response string
	raw = strings.Trim(raw, "\n")

	reqData := strings.Split(raw, "\t")

	if len(reqData) < 2 {
		log.Print("Unable to decode request")
		return BadRequest.Error()
	}

	actionName := reqData[0]

	var action Action

	switch actionName {

	case "Test":
		return buntServer.protobuntVersion

	case "View":
		action = buntServer.db.View

	case "Update":
		action =  buntServer.db.Update
	}

	txName := reqData[1]

	var key, value string
	contextData := strings.Trim(reqData[2], "{}")

	if strings.Contains(contextData, ":") {
		context := strings.Split(contextData, ":")
		key = context[0]
		value = context[1]
	} else {
		key = contextData
		value = ""
	}

	_ = action(func(tx *buntdb.Tx) error {
		if actionName == "View" {
			if txName == "Get" {
				val, _ := tx.Get(key)
				response = val
			}


		}
		if actionName == "Update" {
			if txName == "Set" {
				val, _, _ := tx.Set(key, value, nil)
				response = val
			}
			if txName == "Delete" {
				val, _:= tx.Delete(key)
				response = val
			}
		}
		return nil
	})

	return response + "\n"

}

func (buntServer *BuntServer) handleConnection(c net.Conn) {
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		if err != nil {
			log.Fatal(err)
		}

		result := buntServer.processRequest(netData)

		c.Write([]byte(result + "\n"))
	}
	c.Close()
}

func StartBuntServer(host string, port string, db *buntdb.DB) {

	buntServer := BuntServer{}
	buntServer.protobuntVersion = VERSION
	buntServer.db = db
	defer buntServer.db.Close()

	link := host + ":" + port
	listener, _ := net.Listen("tcp", link)
	log.Println("Starting BuntDB Server...")
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go buntServer.handleConnection(c)
	}
}
