// Package config reads a configuration file and exposes methods to get information from it
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// options describes all the different options/directives supported by the configuration system
// it's a function because we cannot declare slices as `const` and I don't like to have global state.
func options() []option {
	return []option{
		{name: "port", flags: singleFlag},
		{name: "databases", flags: singleFlag},
		{name: "requirepass", flags: singleFlag},
		{name: "save", flags: multipleFlag},
		{name: "include", flags: multipleFlag},
		{name: "appendonly", flags: singleFlag},
		{name: "appendfsync", flags: singleFlag},
		{name: "appenddirname", flags: singleFlag},
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

// NewEmpty returns an empty configuration file. Useful when you don't have a file to read from
func NewEmpty() Config {
	return Config{
		data: make(map[string][]string),
	}
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
	// the callback function receives 3 parameters:
	// 		lineNo: line number in the file (eg: line 42)
	//      key: the specific key (eg: "databases")
	//      value: is the value for the specific key (eg: "16")
	return readFileByLine(config, filePath, func(lineNo int, key, value string) error {
		command, err := checkIsSupported(key)
		if err != nil {
			return err
		}

		if key == "include" {
			// the file must be relative to the original parent
			filename := path.Join(path.Dir(filePath), value)
			if err := includeFile(config, filename, lineNo); err != nil {
				return err
			}
		}

		if command.hasFlag(multipleFlag) {
			// We might have multiple values for the same key, so we append them
			config.data[key] = append(config.data[key], value)
		} else { // singleFlag
			// We can only have a single value for that key, so if multiples keys are found
			// we only keep the last one
			config.data[key] = []string{value}
		}

		return nil
	})
}

func includeFile(config *Config, filename string, lineNo int) error {
	// expand filenames like "*.config" or "redis-?.config". If no expansion is being done,
	// the same file is returned.
	filenames, err := filepath.Glob(filename)
	if err != nil {
		return fmt.Errorf("unable to expand path: %w: %v", ErrInvalidType, err)
	}
	for _, f := range filenames {
		if err := parseConfig(config, f); err != nil {
			return fmt.Errorf("%s:%d : %w", filename, lineNo, err)
		}
	}

	return nil
}

func readFileByLine(config *Config, filename string, processLine func(lineNumber int, key, value string) error) error {
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

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return fmt.Errorf("%w. line: %d: %s", ErrInvalidFile, i, line)
		}

		key, value := parts[0], parts[1]

		err := processLine(i, key, value)
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
