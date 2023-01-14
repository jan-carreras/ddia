// Package config reads a configuration file and exposes methods to get information from it
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

// options describes all the different options/directives supported by the configuration system
// it's a function because we cannot declare slices as `const` and I don't like to have global state.
func options() []option {
	return []option{
		{name: "databases", flags: singleFlag},
		{name: "requirepass", flags: singleFlag},
		{name: "save", flags: multipleFlag},
		{name: "include", flags: multipleFlag},
	}
}

// ErrInvalidFile is returned when there is something wrong with the file (denied permission, it does not exists..)
var ErrInvalidFile = errors.New("invalid file")

// ErrInvalidType when parsing the value of an option to a given type, and cannot be parsed
var ErrInvalidType = errors.New("invalid type")

// ErrUnknownOption returned if the configuration file has an option unknown or not supported
var ErrUnknownOption = errors.New("unknown option")

// ErrCyclicImports is returned if configuration files have an import cycle
var ErrCyclicImports = errors.New("cyclic imports")

// Config object describes all the parameters defined in a redis.config file
// More: https://redis.io/docs/management/config/
type Config struct {
	data          map[string][]string
	filesImported []string
}

// New reads the configuration file and returns a Config object
func New(configPath string) (Config, error) {
	config := &Config{
		data: make(map[string][]string),
	}

	err := parseConfig(config, configPath)
	if err != nil {
		return Config{}, err
	}
	return *config, nil
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

func parseConfig(config *Config, filePath string) error {
	return readFileByLine(config, filePath, func(lineNo int, line string) error {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return fmt.Errorf("%w. line: %d: %s", ErrInvalidFile, lineNo, line)
		}

		key, value := parts[0], parts[1]

		command, err := checkIsSupported(key)
		if err != nil {
			return err
		}

		if key == "include" {
			filename := path.Join(path.Dir(filePath), value)
			if err := parseConfig(config, filename); err != nil {
				return fmt.Errorf("%s:%d : %w", filename, lineNo, err)
			}
		}

		if command.hasFlag(multipleFlag) {
			config.data[key] = append(config.data[key], value)
		} else { // singleFlag
			config.data[key] = []string{value}
		}

		return nil
	})

}

func readFileByLine(config *Config, filename string, processLine func(lineNumber int, line string) error) error {
	if err := fileImportedAlready(config, filename); err != nil {
		return err
	}

	config.filesImported = append(config.filesImported, filename)

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidFile, err)
	}

	defer func() { _ = f.Close() }()

	s := bufio.NewScanner(f)
	for i := 1; s.Scan(); i++ {
		line := s.Text()
		line = strings.Trim(line, " ")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		err := processLine(i, line)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkIsSupported(key string) (option, error) {
	// I know — we could create a hashmap in advance and do a lookup to the table. I
	// would rather keep it simple for the moment, and configuration is loaded once
	// at startup so performance is not a concern.
	for _, o := range options() {
		if o.name == key {
			return o, nil
		}
	}

	return option{}, fmt.Errorf("unknown or unupported option %q: %w", key, ErrUnknownOption)
}

func fileImportedAlready(config *Config, filepath string) error {
	for _, f := range config.filesImported {
		if filepath == f {
			return fmt.Errorf("%w: %s tried to be imported twice", ErrCyclicImports, filepath)
		}
	}

	return nil
}
