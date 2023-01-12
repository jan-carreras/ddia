package config

import (
	"reflect"
	"testing"
)

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

	value, err = config.Integer("loglevel", 60)
	if err == nil {
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

	want := []string{"900 1", "300 10", "60 10000"}

	if !reflect.DeepEqual(values, want) {
		t.Fatalf("incorrect values: %v, want %v", values, want)
	}

}
