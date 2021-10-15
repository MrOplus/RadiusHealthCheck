package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Parse(path string) (result *Config) {
	content , err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr,"Unable to read the file")
		os.Exit(1)
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		fmt.Fprintln(os.Stderr,"Unable to parse config file")
		os.Exit(1)
	}
	return result
}
