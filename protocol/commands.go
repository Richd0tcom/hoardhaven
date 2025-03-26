package protocol

import "sync"


var Handlers  = map[string]func([]Value) Value{
	"PING": ping,
	"SET" : set,
	"GET" : get,
}


var DB = struct {
	Map map[string]string
	lck sync.RWMutex
}{
	Map: map[string]string{},
	lck: sync.RWMutex{},
}
func ping(arg []Value) Value {
	if len(arg) == 0 {
		return Value{Type: "string", Str: "PONG"}
	}

	return Value{Type: "string", Str: arg[0].Bulk}
}

func set(args []Value) Value {

	if len(args) != 2 {
		return Value{Type: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key:= args[0].Bulk
	value := args[1].Bulk

	DB.lck.Lock()
	DB.Map[key] = value
	DB.lck.Unlock()

	return Value{Type: "string", Str: "OK"}
}


func get(args []Value) Value {
	if len(args) != 1 {
		return Value{Type: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key:= args[0].Bulk

	DB.lck.RLock()
	value, ok:= DB.Map[key]
	DB.lck.RUnlock()

	if !ok {
		return Value{Type: "null"}
	}

	
	return Value{Type: "bulk", Bulk: value}
}



