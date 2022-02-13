package main

import (
	"flag"
	"fmt"
	"os"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

type cmd struct {
	out string
}

func (m *cmd) register(f *flag.FlagSet) {
	f.StringVar(&m.out, "out", "-", "Output file to write to")
}

func (m *cmd) run(args []string) error {
	return nil
}

func main() {
	m := new(cmd)

	m.register(flag.CommandLine)

	if err := m.run(flag.Args()); err != nil {
		die(err)
	}
}
