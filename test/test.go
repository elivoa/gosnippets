package main

import (
	"encoding/json"
	"fmt"
)

type Container struct {
	Guest1 `json:,inline`
	Guest2 `json:"g2"`
}

type Guest1 struct {
	Name string
}

type Guest2 struct {
	Name     string
	LastName string
}

func main() {
	a := Container{
		Guest1: Guest1{
			Name: "Tit",
		},
		Guest2: Guest2{
			LastName: "Petric",
		},
	}
	b, _ := json.MarshalIndent(a, "", "  ")

	fmt.Println(string(b))

	var testContainer Container
	err := json.Unmarshal([]byte(`{"Name": "Tit"}`), &testContainer)
	fmt.Println(
		testContainer.Guest1,
		testContainer.Guest2,
		err,
	)
}
