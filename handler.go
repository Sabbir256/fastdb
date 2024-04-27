/*
	Handler for processing commands from Redis clients.
*/

package main

import "sync"

var db = map[string]string{}
var dbMutex = sync.RWMutex{}
var hashdb = map[string]map[string]string{}
var hashdbMutex = sync.RWMutex{}

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
	"HGETALL": hgetall,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	dbMutex.Lock()
	db[key] = value
	dbMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	dbMutex.RLock()
	value, ok := db[key]
	dbMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	hashdbMutex.Lock()
	if _, ok := hashdb[hash]; !ok {
		hashdb[hash] = map[string]string{}
	}
	hashdb[hash][key] = value
	hashdbMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	hashdbMutex.RLock()
	value, ok := hashdb[hash][key]
	hashdbMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	hashdbMutex.RLock()
	value, ok := hashdb[hash]
	hashdbMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "array", array: mapToValues(value)}
}

func mapToValues(m map[string]string) []Value {
	values := make([]Value, 0)
	for k, v := range m {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}
	return values
}
