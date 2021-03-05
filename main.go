package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type Connection struct {
	IP   string `json:ip`
	Port string `json: port`
	Nick string `json: nick`
}

var connections []Connection

func remove(slice []Connection, s int) []Connection {
	return append(slice[:s], slice[s+1:]...)
}

func handleConnection(conn net.Conn) {
	data := make([]byte, 256)
	var newConnection Connection
	connAddr := conn.RemoteAddr().String()
	index := strings.LastIndex(connAddr, ":")

	newConnection.IP = connAddr[:index]
	newConnection.Port = connAddr[index+1:]
	fmt.Println(newConnection.IP, newConnection.Port)

	conn.Read(data)
	newConnection.Nick = strings.Trim(string(data)[:], "\u0000")
	fmt.Println(newConnection.Nick)

	for i := range connections {
		json.NewEncoder(conn).Encode(connections[i])
	}
	connections = append(connections, newConnection)
	fmt.Println(connections)

	var lastReply string
	lastReply = "ok"
	for lastReply != "bye" {
		lastReply = "bye"
		//conn.Read(data)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		lastReply = strings.Trim(message, "\n")
	}
	conn.Close()
	for i := range connections {
		if newConnection.IP == connections[i].IP && newConnection.Nick == connections[i].Nick {
			connections = remove(connections, i)
			break
		}
	}
	/*for v := range connections {

		for i := range connections {
			json.NewEncoder(tmpConn).Encode(connections[i])
		}
	}*/
	fmt.Println(connections)
	//conn.Write([]byte("dance you fuckerr"))
}

func main() {
	/*	var mockData Connection
		mockData.IP = "umutaev.ru"
		mockData.Nick = "max"
		mockData.Port = "1488"
		connections = append(connections, mockData)*/
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handleConnection(conn)
	}
}
