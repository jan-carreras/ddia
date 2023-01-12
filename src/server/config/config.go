package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var ErrInvalidFile = errors.New("invalid file")
var ErrInvalidType = errors.New("invalid type")

type Config struct {
	data map[string][]string
}

func New(configPath string) (Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("unable to open config file: %q: %w", configPath, err)
	}

	defer func() { _ = f.Close() }()
	return read(f)
}

// Integer returns the value of key as integer. If key does not exist, return def. If not integer, returns ErrInvalidType
func (c Config) Integer(key string, def int) (int, error) {
	value, ok := c.Get(key)
	if !ok {
		return def, nil
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%w: unable to convert %q to int: %v", ErrInvalidType, value, err)
	}

	return i, nil
}

// Get returns the first value for key found on the file
func (c Config) Get(key string) (value string, ok bool) {
	values, ok := c.data[key]
	if !ok || len(values) == 0 {
		return "", false
	}

	return values[0], true
}

// GetD gets the value for key, returns def if key not found
func (c Config) GetD(key, def string) string {
	value, ok := c.Get(key)
	if ok {
		return value
	}

	return def
}

// GetM returns all the values for the given key
func (c Config) GetM(key string) (values []string, ok bool) {
	values, ok = c.data[key]
	return values, ok
}

func read(input io.Reader) (Config, error) {
	config := Config{
		data: make(map[string][]string),
	}

	s := bufio.NewScanner(input)
	for i := 1; s.Scan(); i++ {
		line := s.Text()
		line = strings.Trim(line, " ")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("%w. line: %d: %s", ErrInvalidFile, i, line)
		}

		key, value := parts[0], parts[1]

		config.data[key] = append(config.data[key], value)
	}

	return config, nil
}
