package parsables

import (
	"encoding/json"
	"fmt"
)

func ExampleDurationFromString() {
	s, err := DurationFromString("24m25s")
	if err != nil {
		panic(err)
	}

	fmt.Println(s.String())

	// Output:
	// 24m25s
}

func ExampleDuration_UnmarshalText() {
	var s Duration

	err := s.UnmarshalText([]byte("35ms"))
	if err != nil {
		panic(err)
	}

	fmt.Println(s)

	// Output:
	// 35ms
}

func ExampleDuration_json() {
	type Config struct {
		Interval Duration `json:"interval"`
	}

	var cfg Config
	err := json.Unmarshal([]byte(`{"interval": "5m"}`), &cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)

	// Output:
	// {Interval:5m0s}
}
