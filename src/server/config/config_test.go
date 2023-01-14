package config

import (
	"errors"
	"reflect"
	"testing"
)

func TestCommandConsistency(t *testing.T) {
	t.Run("check for duplicates", func(t *testing.T) {
		set := make(map[string]bool)
		for _, o := range options() {
			if _, ok := set[o.name]; ok {
				t.Fatalf("duplicate found: %s", o.name)
			}
			set[o.name] = true
		}
	})

	t.Run("check for empty options", func(t *testing.T) {
		for _, o := range options() {
			if o.name == "" {
				t.Fatalf("option with empty name: %v. All options should have a name", o)
			}
		}
	})
}

func TestConfig_Get(t *testing.T) {
	config, err := New("testdata/redis.conf")
	if err != nil {
		t.Fatalf("expected no error: %q", err.Error())
	}

	value, ok := config.Get("databases")
	if !ok {
		t.Fatalf("key expected to be found")
	}

	if want := "16"; value != want {
		t.Fatalf("invalid databases value: %q, want %q", value, want)
	}

	value, ok = config.Get("requirepass")
	if !ok {
		t.Fatalf("key expected to be found")
	}

	if want := "hello-there-2"; value != want {
		t.Fatalf("invalid databases value: %q, want %q", value, want)
	}
}

func TestConfig_Integer(t *testing.T) {
	config, err := New("testdata/redis.conf")
	if err != nil {
		t.Fatalf("expected no error: %q", err.Error())
	}

	value, err := config.Integer("databases", 0)
	if err != nil {
		t.Fatalf("expecting no error: %q", err.Error())
	}

	if want := 16; value != want {
		t.Fatalf("invalid databases: %d, want %d", value, want)
	}

	value, err = config.Integer("databases-non-existing", 60)
	if err != nil {
		t.Fatalf("expecting no error: %q", err.Error())
	}

	if want := 60; value != want {
		t.Fatalf("invalid databases: %d, want %d", value, want)
	}

	// requirepass should be parsed as a string, not integer. Should always return error
	value, err = config.Integer("requirepass", 60)
	if !errors.Is(err, ErrInvalidType) {
		t.Fatalf("expecting error: %q, want %q", err.Error(), ErrInvalidType)
	}

	if want := 0; value != want {
		t.Fatalf("invalid databases: %d, want %d", value, want)
	}
}

func TestConfig_GetD(t *testing.T) {
	config, err := New("testdata/redis.conf")
	if err != nil {
		t.Fatalf("expected no error: %q", err.Error())
	}

	value := config.GetD("databases", "whatever")
	if want := "16"; value != want {
		t.Fatalf("invalid databases value: %q, want %q", value, want)
	}

	value = config.GetD("non-existing-key", "default")
	if want := "default"; value != want {
		t.Fatalf("invalid default value returned: %q, want %q", value, want)
	}
}

func TestConfig_NonExistingKey(t *testing.T) {
	config, err := New("testdata/redis.conf")
	if err != nil {
		t.Fatalf("expected no error: %q", err.Error())
	}

	value, ok := config.Get("non-existent-key-1234")
	if ok {
		t.Fatalf("key expected not to be found")
	}
	if want := ""; value != want {
		t.Fatalf("expecting key to be empty: %q, want %q", value, want)
	}
}

func TestConfig_GetM(t *testing.T) {
	config, err := New("testdata/redis.conf")
	if err != nil {
		t.Fatalf("expected no error: %q", err.Error())
	}

	values, ok := config.GetM("save")
	if !ok {
		t.Fatalf("key expected to be found")
	}

	want := []string{"900 1", "300 10", "60 10000", "30 100000"}

	if !reflect.DeepEqual(values, want) {
		t.Fatalf("incorrect values: %v, want %v", values, want)
	}

}

func TestNew_InvalidFile(t *testing.T) {
	_, err := New("testdata/non-existing-file.conf")
	if !errors.Is(err, ErrInvalidFile) {
		t.Fatalf("expected error: %q, want %q", err, ErrInvalidFile)
	}
}
