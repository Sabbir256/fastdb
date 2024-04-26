package main

import (
	"fmt"
	"net"
)

const REDIS_DEFAULT_PORT = ":6379"
func main() {
	fmt.Println("Listening on port", REDIS_DEFAULT_PORT)

	ln, err := net.Listen("tcp", REDIS_DEFAULT_PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := NewRespReader(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
