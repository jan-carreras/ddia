package storage_test

import (
	"ddia/src/storage"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestInMemory_GetSet(t *testing.T) {
	k, v := "hello", "world"

	store := storage.NewInMemory()
	err := store.Set(k, v)
	require.NoError(t, err)

	readValue, err := store.Get(k)
	require.NoError(t, err)
	require.Equal(t, v, readValue)
}

func TestInMemory_Set_ValueOverwrite(t *testing.T) {
	k, v, v2 := "hello", "world", "cruel world"

	store := storage.NewInMemory()
	err := store.Set(k, v)
	require.NoError(t, err)

	readValue, err := store.Get(k)
	require.NoError(t, err)
	require.Equal(t, v, readValue)

	err = store.Set(k, v2)
	require.NoError(t, err)

	readValue, err = store.Get(k)
	require.NoError(t, err)
	require.Equal(t, v2, readValue)
}

func TestInMemory_GetSet_Concurrent(t *testing.T) {
	k, v := "hello", "world"

	store := storage.NewInMemory()
	err := store.Set(k, v)
	require.NoError(t, err)

	wg := sync.WaitGroup{}

	write := func() {
		defer wg.Done()
		err := store.Set(k, v)
		require.NoError(t, err)
	}

	read := func() {
		defer wg.Done()
		readValue, err := store.Get(k)
		require.NoError(t, err)
		require.Equal(t, v, readValue)
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
	require.Empty(t, v)

	require.ErrorIs(t, err, storage.ErrNotFound)

}
