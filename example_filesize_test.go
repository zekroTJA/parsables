package parsables

import (
	"encoding/json"
	"fmt"
)

func ExampleFileSizeFromString() {
	s, err := FileSizeFromString("420 Gigabyte")
	if err != nil {
		panic(err)
	}

	fmt.Println(s)

	// Output:
	// 420000000000
}

func ExampleFileSize_UnmarshalText() {
	var s FileSize

	err := s.UnmarshalText([]byte("12.34 GiB"))
	if err != nil {
		panic(err)
	}

	fmt.Println(s)

	// Output:
	// 13249974108
}

func ExampleFileSize_json() {
	type Config struct {
		MaxSize FileSize `json:"max_size"`
	}

	var cfg Config
	err := json.Unmarshal([]byte(`{"max_size": "512gb"}`), &cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)

	// Output:
	// {MaxSize:512000000000}
}
