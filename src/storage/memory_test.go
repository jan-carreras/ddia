package storage_test

import (
	"ddia/src/server"
	"ddia/src/storage"
	"errors"
	"testing"
)

func TestInMemory_GetSet(t *testing.T) {
	k, v := "hello", "world"

	store := storage.NewInMemory()
	err := store.Set(k, v)
	if err != nil {
		t.Fatalf("error returned: %v wanted no error", err)
	}

	readValue, err := store.Get(k)
	if err != nil {
		t.Fatalf("error returned: %v wanted no error", err)
	}

	if want := readValue; v != want {
		t.Fatalf("invalid read value: %q, want %q", v, want)
	}
}

func TestInMemory_Set_ValueOverwrite(t *testing.T) {
	k, v, v2 := "hello", "world", "cruel world"

	store := storage.NewInMemory()
	err := store.Set(k, v)
	if err != nil {
		t.Fatalf("error returned: %v wanted no error", err)
	}

	readValue, err := store.Get(k)
	if err != nil {
		t.Fatalf("error returned: %v wanted no error", err)
	}
	if want := readValue; v != want {
		t.Fatalf("invalid read value: %q, want %q", v, want)
	}

	err = store.Set(k, v2)
	if err != nil {
		t.Fatalf("error returned: %v wanted no error", err)
	}

	readValue, err = store.Get(k)
	if err != nil {
		t.Fatalf("error returned: %v wanted no error", err)
	}
	if want := v2; readValue != want {
		t.Fatalf("invalid read value: %q, want %q", v, want)
	}
}

func TestIncrementBy(t *testing.T) {
	store := storage.NewInMemory()
	v, err := store.IncrementBy("visits", 10)
	if err != nil {
		t.Fatalf("expeced no error: %q", err.Error())
	}

	if want := "10"; v != want {
		t.Fatalf("invalid count: %s want %s", v, want)
	}

	_, _ = store.IncrementBy("visits", 20) // 30
	v, _ = store.IncrementBy("visits", -5) // 25

	if want := "25"; v != want {
		t.Fatalf("invalid count: %s want %s", v, want)
	}
}

func TestInMemory_Get_NonExisting(t *testing.T) {
	store := storage.NewInMemory()
	v, err := store.Get("non-existing-key")

	if want := ""; v != want {
		t.Fatalf("expecting nothing: %q returned", v)
	}

	if !errors.Is(err, server.ErrNotFound) {
		t.Fatalf("incorrect error returned: %q, want %q", err.Error(), server.ErrNotFound.Error())
	}
}
