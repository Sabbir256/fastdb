package main

import (
	"fmt"
	"net"
	"strings"
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

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array.");
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, empty array found!")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Unknown command: " + command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
