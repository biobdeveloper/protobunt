package protobunt

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	GET = "Get"
	SET = "Set"
	DELETE = "Delete"
)


type BuntClient struct {
	conn net.Conn
	dbName string
	protobuntVersion string
}


func (cli *BuntClient) call(action string, tx string, txData string) string {
	rawRequestString := strings.Join([]string{action, tx, txData}, "\t") + "\n"
	cli.conn.Write([]byte(rawRequestString))
	message, _ := bufio.NewReader(cli.conn).ReadString('\n')
	return strings.Trim(message, "\n")
}

func (cli *BuntClient) testConnection() {
	testConn := cli.call("Test", cli.protobuntVersion, "")
	if testConn != cli.protobuntVersion {
		fmt.Printf("Warning! Different versions: client %s, server %s", cli.protobuntVersion, testConn)
	}
}

func (cli *BuntClient) View(tx string, key string) string {
	if tx == "" {
		tx = GET
	}
	return cli.call("View", tx, "{" + key + "}")
}

func (cli *BuntClient) Update(tx, key, value string) string {
	if tx == "" {
		tx = SET
	}
	return cli.call("Update", tx,"{" + key + ":" + value + "}")
}

func CreateBuntClient(host string, port string) *BuntClient {

	client := new(BuntClient)
	var link = host + ":" + port
	client.conn, _ = net.Dial("tcp", link)

	log.Println("BuntDB client is ready...")

	client.protobuntVersion = VERSION
	client.testConnection()

	return client
}

