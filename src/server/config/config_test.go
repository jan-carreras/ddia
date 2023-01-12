package config

import (
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	config, err := New("testdata/redis.conf")
	if err != nil {
		t.Fatalf("expected no error: %q", err.Error())
	}

	t.Run("single key", func(t *testing.T) {
		value, ok := config.Get("databases")
		if !ok {
			t.Fatalf("key expected to be found")
		}

		if want := "16"; value != want {
			t.Fatalf("invalid databases value: %q, want %q", value, want)
		}
	})

	t.Run("single key, passing default", func(t *testing.T) {
		value := config.GetD("databases", "whatever")
		if want := "16"; value != want {
			t.Fatalf("invalid databases value: %q, want %q", value, want)
		}
	})

	t.Run("non existing key", func(t *testing.T) {
		value, ok := config.Get("non-existent-key-1234")
		if ok {
			t.Fatalf("key expected not to be found")
		}
		if want := ""; value != want {
			t.Fatalf("expecting key to be empty: %q, want %q", value, want)
		}
	})

	t.Run("non existing key with default", func(t *testing.T) {
		value := config.GetD("non-existing-key", "default")
		if want := "default"; value != want {
			t.Fatalf("invalid default value returned: %q, want %q", value, want)
		}

	})

	t.Run("multi key", func(t *testing.T) {
		values, ok := config.GetM("save")
		if !ok {
			t.Fatalf("key expected to be found")
		}

		want := []string{"900 1", "300 10", "60 10000"}

		if !reflect.DeepEqual(values, want) {
			t.Fatalf("incorrect values: %v, want %v", values, want)
		}
	})

}
