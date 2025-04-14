package main

import (
	"fmt"
	"sync"
)

var handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hSet,
	"HGET": hGet,
}
var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hSet(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	key := args[0].bulk
	hash := args[1].bulk
	value := args[2].bulk

	fmt.Println("hset: params: ", key, hash, value)

	HSETsMu.Lock()
	if _, ok := HSETs[key]; !ok {
		HSETs[key] = map[string]string{}
	}
	HSETs[key][hash] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "ok"}
}

func hGet(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	key := args[0].bulk
	hash := args[1].bulk

	fmt.Println("hget: params: ", key, hash)

	HSETsMu.RLock()
	value, ok := HSETs[key][hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}
