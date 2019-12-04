package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	prompt "github.com/c-bata/go-prompt"
	"github.com/go-redis/redis"
)

var db *redis.Client

var CmdSuggests = []prompt.Suggest{
	{"APPEND", " value append a value to a key"},
	{"BITCOUNT", "count set bits in a string"},
	{"SET", "set value in key"},
	{"SETNX", "set if not exist value in key"},
	{"SETRANGE", "value overwrite part of a string at key starting at the specified offset"},
	{"STRLEN", "get the length of the value stored in a key"},
	{"MSET", "set multiple keys to multiple values"},
	{"MSETNX", "set multiple keys to multiple values , only if none of the keys exist"},
	{"GET", "get value in key"},
	{"GETRANGE", "value get a substring value of a key and return its old value"},
	{"MGET", "get the values of all the given keys"},
	{"INCR", "increment value in key"},
	{"INCRBY", "increment increment the integer value of a key by the given amount"},
	{"INCRBYFLOAT", "increment increment the float value of a key by the given amount"},
	{"DECR", "decrement the integer value of key by one"},
	{"DECRBY", "decrement decrement the integer value of a key by the given number"},
	{"DEL", "delete key"},
	{"exit", "close the prompt"},
}

var StringSuggests []prompt.Suggest
var HashSuggests []prompt.Suggest
var ListSuggests []prompt.Suggest
var HashTableSuggests = make(map[string][]prompt.Suggest)

var sets = make(map[string]bool)

func executor(in string) {
	in = strings.TrimSpace(in)
	args := strings.Split(in, " ")
	if args[0] == "exit" {
		os.Exit(0)
	}
	var cmd *redis.StringCmd
	if len(args) < 2 {
		fmt.Println("too few arguments")
		return
	}
	switch args[0] {
	case "hget":
		cmd = redis.NewStringCmd(args[0], args[1], args[2])
		err := db.Process(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "select":
		d, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		currentDB = d
		db = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", DB: d})
		StringSuggests = []prompt.Suggest{}
		HashSuggests = []prompt.Suggest{}
		sets = make(map[string]bool)
	default:
		cmd = redis.NewStringCmd(args[0], args[1])
		err := db.Process(cmd)
		if err != nil {
			fmt.Printf("default: %v\n", err)
			return
		}
	}
	if cmd == nil {
		fmt.Println()
		return
	}
	res, err := cmd.Result()
	if err != nil {
		fmt.Printf("78: %v\n", err)
	}
	fmt.Println(res)
}

func completer(in prompt.Document) []prompt.Suggest {
	line := in.CurrentLineBeforeCursor()
	args := strings.Split(line, " ")
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(CmdSuggests, in.GetWordBeforeCursor(), true)
	}
	switch strings.ToLower(args[0]) {
	case "get":
		return prompt.FilterFuzzy(StringSuggests, in.GetWordBeforeCursor(), true)
	case "hget":
		if len(args) <= 2 {
			return prompt.FilterFuzzy(HashSuggests, in.GetWordBeforeCursor(), true)
		}
		return prompt.FilterFuzzy(HashTableKeys(args[1]), in.GetWordBeforeCursor(), true)
	}
	return prompt.FilterHasPrefix(CmdSuggests, in.GetWordBeforeCursor(), true)
}

func HashTableKeys(table string) []prompt.Suggest {
	if len(HashTableSuggests[table]) > 0 {
		return HashTableSuggests[table]
	}
	res, _ := db.HKeys(table).Result()
	var sugs []prompt.Suggest
	for _, v := range res {
		sugs = append(sugs, prompt.Suggest{Text: v})
	}
	if sugs == nil {
		return nil
	}
	HashTableSuggests[table] = sugs
	return sugs
}

func refresh() {
	for {
		time.Sleep(time.Second * 1)
		res, _ := db.Keys("*").Result()
		for _, v := range res {
			if !sets[v] {
				sets[v] = true
				t, _ := db.Type(v).Result()
				switch t {
				case "hash":
					HashSuggests = append(HashSuggests, prompt.Suggest{Text: v, Description: t})
				case "string":
					StringSuggests = append(StringSuggests, prompt.Suggest{Text: v, Description: t})
				}
			}
		}
	}
}

var currentDB = 0

func livePrefix() (string, bool) {
	return fmt.Sprintf("redis[%d]>", currentDB), true
}

func main() {
	db = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", DB: 0})
	go refresh()
	p := prompt.New(
		executor,
		completer,
		prompt.OptionLivePrefix(livePrefix),
		prompt.OptionPrefix("redis> "),
		prompt.OptionTitle("redis-prompt"),
	)
	p.Run()
}
