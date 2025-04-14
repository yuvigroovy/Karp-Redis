package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func newResp(reader io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(reader)}
}

func readLine(resp *Resp) (line []byte, n int, err error) {
	for {
		bytes, err := resp.reader.ReadByte()

		if err != nil {
			return nil, 0, err
		}

		n += 1

		line = append(line, bytes)

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	return line[:len(line)-2], n, nil
}

func readInteger(resp *Resp) (number int, n int, err error) {
	line, n, err := readLine(resp)

	if err != nil {
		return 0, 0, err
	}

	number64, err := strconv.ParseInt(string(line), 10, 64)

	if err != nil {
		return 0, 0, err
	}

	return int(number64), n, nil
}

func read(resp *Resp) (Value, error) {
	_type, err := resp.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return readArray(resp)
	case BULK:
		return readBulk(resp)
	default:
		fmt.Println("Invalid syntax")
		return Value{}, nil
	}
}

func readArray(resp *Resp) (Value, error) {
	value := Value{}
	value.typ = "array"

	length, _, err := readInteger(resp)

	if err != nil {
		return Value{}, err
	}

	value.array = make([]Value, length)

	for index := 0; index < length; index++ {
		elementValue, err := read(resp)

		if err != nil {
			return value, err
		}

		value.array[index] = elementValue
	}

	return value, nil
}

func readBulk(resp *Resp) (Value, error) {
	value := Value{}
	value.typ = "bulk"

	length, _, err := readInteger(resp)

	if err != nil {
		return Value{}, err
	}

	bulkBytes := make([]byte, length)

	resp.reader.Read(bulkBytes)
	value.bulk = string(bulkBytes)

	readLine(resp)

	return value, nil
}

func (value Value) marshal() []byte {
	switch value.typ {
	case "array":
		return value.marshalArray()
	case "bulk":
		return value.marshalBulk()
	case "string":
		return value.marshalString()
	case "error":
		return value.marshallError()
	case "null":
		return value.marshallNull()
	default:
		return []byte{}
	}
}

func (value Value) marshalString() []byte {
	var bytes []byte

	bytes = append(bytes, STRING)
	bytes = append(bytes, value.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (value Value) marshalBulk() []byte {
	var bytes []byte

	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(value.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, value.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (value Value) marshalArray() []byte {
	var bytes []byte

	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(value.array))...)
	bytes = append(bytes, '\r', '\n')

	for index := 0; index < len(value.array); index++ {
		bytes = append(bytes, value.array[index].marshal()...)
	}

	return bytes
}

func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}
