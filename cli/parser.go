package cli

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)
type Options struct {
	Config string `short:"c" long:"config" description:"Config Path" required:"true"`
}

func Parse() (result * Options) {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
	if !fileExists(opts.Config) {
		_, _ = fmt.Fprintln(os.Stderr, "specified file not exists")
		os.Exit(1)
	}
	return &opts
}
func fileExists(filename string) bool{
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}