package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/BurntSushi/toml"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

type cmd struct {
	out string
}

func (m *cmd) register(f *flag.FlagSet) {
	f.StringVar(&m.out, "out", "-", "Output file to write to")
}

func (m *cmd) run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("command takes 1 arg (dsn string), got %d args", len(args))
	}

	c, err := makeClientConfig(args[0])
	if err != nil {
		return fmt.Errorf("error making my.cnf configuration: %w", err)
	}

	f, err := m.output()
	if err != nil {
		return fmt.Errorf("error getting output: %w", err)
	}

	if err := toml.NewEncoder(f).Encode(&Config{Client: c}); err != nil {
		return fmt.Errorf("error encoding my.cnf configuration: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("error closing file: %w", err)

	}

	return nil
}

func (m *cmd) output() (*os.File, error) {
	if m.out == "-" {
		return os.Stdout, nil
	}

	f, err := os.Create("path/to/file.toml")
	if err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}

	return f, nil
}

func main() {
	m := new(cmd)

	m.register(flag.CommandLine)

	flag.Parse()

	if err := m.run(flag.Args()); err != nil {
		die(err)
	}
}

var (
	reFull = regexp.MustCompile(`(?P<user>[^:]+):(?P<password>[^@]+)@tcp\((?P<host>[^:]+):(?P<port>\d+)\)\/(?P<database>[\d\w-_]+).*`)
)

type Config struct {
	Client *ClientConfig `toml:"client"`
}

type ClientConfig struct {
	Host     string `toml:"host,omitempty" json:"host,omitempty"`
	Port     int    `toml:"port,omitempty" json:"port,omitempty,string"`
	User     string `toml:"user,omitempty" json:"user,omitempty"`
	Password string `toml:"password,omitempty" json:"password,omitempty"`
	Database string `toml:"database,omitempty" json:"database,omitempty"`
	SSLMode  string `toml:"ssl-mode,omitempty" json:"ssl-mode,omitempty"`
}

func makeClientConfig(dsn string) (*ClientConfig, error) {
	var (
		m = reFull.FindStringSubmatch(dsn)
		v = make(map[string]string)
	)

	for i, group := range reFull.SubexpNames()[0:] {
		if i != 0 && group != "" {
			v[group] = m[i]
		}
	}

	p, err := jsonMarshal(v)
	if err != nil {
		return nil, fmt.Errorf("error marshaling: %w", err)
	}

	c := &ClientConfig{
		Port:    3306,
		SSLMode: "PREFERRED",
	}

	if err := jsonUnmarshal(p, c); err != nil {
		return nil, fmt.Errorf("error unmarshaling: %w", err)
	}

	return c, nil
}

func jsonMarshal(v interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func jsonUnmarshal(p []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(p))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
