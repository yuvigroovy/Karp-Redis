package main

import (
	"fmt"
	"net"
	"strings"
)

func startServer(port string) net.Listener {
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println(err.Error())

		return nil
	}

	return listener
}

func main() {
	port := "6379"
	listener := startServer(port)
	fmt.Println("started server on port: " + port)

	connection, err := listener.Accept()
	fmt.Println("accepting connection")

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	defer connection.Close()

	for {
		resp := newResp(connection)
		writer := NewWriter(connection)

		value, err := read(resp)

		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		writer.write(result)
	}
}
