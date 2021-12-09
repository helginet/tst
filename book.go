package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var data = "./data.json"
var companies = map[string]bool{
	"coke":  true,
	"pepsi": true,
}
var actions = map[string]bool{
	"book":   true,
	"cancel": true,
	"list":   true,
}

func main() {
	gogo()
}

func JSONString(any interface{}) string {
	content, err := json.Marshal(any)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
	return string(content)
}

func storeit(what string) {
	if f, err := os.OpenFile(data, os.O_WRONLY|os.O_CREATE, 0644); err == nil {
		defer f.Close()
		f.WriteString(what)
	} else {
		panic(fmt.Sprintf("%s", err))
	}
}

func PrettyFormat(v interface{}) string {
	result, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
	return string(result)
}

func gogo() {
	var company, action string
	var roomN, hours int
	flag.StringVar(&company, "company", "", "company name: coke or pepsi")
	flag.IntVar(&roomN, "room", 0, "room number to book: [1,10]")
	flag.IntVar(&hours, "hours", 0, "hours to book: [1, 2^63-1]")
	flag.StringVar(&action, "action", "", "action: book, cancel or list")
	flag.Parse()
	if !companies[company] {
		fmt.Println("Wrong company. See the script help.")
		return
	}
	if !actions[action] {
		fmt.Println("Wrong action. See the script help.")
		return
	}
	if action != "list" {
		if roomN < 1 || roomN > 10 {
			fmt.Println("Wrong room number.")
			return
		}
		if action == "book" && hours < 1 {
			fmt.Println("Wrong hours to book.")
			return
		}
	}
	rooms := make(map[string][]int)
	fileContent, err := ioutil.ReadFile(data)
	if err != nil {
		for companyName := range companies {
			for i := 0; i < 10; i++ {
				rooms[companyName] = append(rooms[companyName], 0)
			}
		}
		storeit(JSONString(rooms))
	} else {
		if err = json.Unmarshal(fileContent, &rooms); err != nil {
			panic(err)
		}
	}
	switch action {
	case "book":
		if rooms[company][roomN-1] == 0 {
			rooms[company][roomN-1] = hours
			fmt.Println("Room booked successfully.")
		} else {
			fmt.Println("This room booked already.")
		}
	case "cancel":
		rooms[company][roomN-1] = 0
		fmt.Println("Room cancelled successfully.")
	}
	if action == "list" {
		for i := 0; i < 10; i++ {
			fmt.Printf("Room: %2d, booked hours: %d\n", i+1, rooms[company][i])
		}
	}
	storeit(JSONString(rooms))
}
