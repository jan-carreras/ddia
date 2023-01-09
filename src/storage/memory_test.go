package storage_test

import (
	"ddia/src/server"
	"ddia/src/storage"
	"errors"
	"sync"
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

func TestInMemory_GetSet_Concurrent(t *testing.T) {
	k, v := "hello", "world"

	store := storage.NewInMemory()
	err := store.Set(k, v)
	if err != nil {
		t.Fatalf("error returned: %v wanted no error", err)
	}

	wg := sync.WaitGroup{}

	write := func() {
		defer wg.Done()
		err := store.Set(k, v)
		if err != nil {
			t.Fatalf("error returned: %v wanted no error", err)
		}
	}

	read := func() {
		defer wg.Done()
		readValue, err := store.Get(k)
		if err != nil {
			t.Fatalf("error returned: %v wanted no error", err)
		}
		if want := readValue; v != want {
			t.Fatalf("invalid read value: %q, want %q", v, want)
		}
	}

	for i := 0; i < 30; i++ {
		wg.Add(2)
		go write()
		go read()

	}
	wg.Wait()
}

func TestInMemory_Get_NonExisting(t *testing.T) {
	store := storage.NewInMemory()
	v, err := store.Get("non-existing-key")

	if want := ""; v != want {
		t.Fatalf("expecting nothing: %q returned", v)
	}

	if !errors.Is(err, server.ErrNotFound) {
		t.Fatalf("incorrect error returned: %v, want %v", err, server.ErrNotFound)
	}
}
