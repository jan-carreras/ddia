package storage

import (
	"container/list"
	"ddia/src/server"
	"errors"
	"fmt"
)

// LRange returns the specified elements of the list stored at key.
// More: https://redis.io/commands/lrange/
//
// TODO: Optimisation, we don't necessarily need to start from the Start of the list. We can traverse it
// starting from the end, if the requested range is "nearer" the end.
func (m *InMemory) LRange(key string, start, stop int) (values []string, err error) {
	l, err := m.listGetKeyOrNew(key)
	if err != nil {
		return nil, err
	}

	toPositive := func(n int) int {
		if n >= 0 {
			return n
		}
		return l.Len() + n
	}

	start, stop = toPositive(start), toPositive(stop)

	// If start is larger than the end of the list, an empty list is returned
	if start > l.Len() {
		return values, nil
	}

	// If stop is larger than the actual end of the list, Redis will treat it like
	// the last element of the list
	if stop >= l.Len() {
		stop = l.Len() - 1
	}

	n := l.Front()
	for i := 0; i < l.Len() && i <= stop; i++ {
		if i >= start && i <= stop {
			v, ok := n.Value.(string)
			if !ok {
				return nil, ErrTypeCorruption
			}
			values = append(values, v)
		}

		n = n.Next()
	}

	return values, nil
}

// LRem removes the first count occurrences of elements equal to element from the list stored at key
//
// TODO: Optimisation, we don't necessarily need to start from the Start of the list. We can traverse it
// starting from the end, if the requested range is "nearer" the end.
func (m *InMemory) LRem(key string, count int, element string) (removed int, err error) {
	l, err := m.listGetKeyOrNew(key)
	if err != nil {
		return 0, err
	}

	// We want to remove all the occurrences of element, thus we set count to the list length
	if count == 0 {
		count = l.Len()
	}

	moveFront := func(l *list.Element) *list.Element { return l.Next() }
	movePrev := func(l *list.Element) *list.Element { return l.Prev() }

	move, n := moveFront, l.Front()
	if count < 0 {
		move, n = movePrev, l.Back()
		count = -count
	}

	for i := 0; i < l.Len() && count > 0; i++ {
		next := move(n)

		if v, ok := n.Value.(string); !ok {
			return 0, ErrTypeCorruption
		} else if v == element {
			l.Remove(n)
			count--
			removed++
		}

		n = next
	}

	m.saveList(key, l)

	return removed, nil
}

// LIndex returns the element at "index" index in the list stored at key.
//
// TODO: Optimisation, we don't necessarily need to start from the Start of the list. We can traverse it
// starting from the end, if the requested range is "nearer" the end.
func (m *InMemory) LIndex(key string, index int) (string, error) {
	l, err := m.listGetKey(key)
	if err != nil {
		return "", err
	}

	n, err := m.listFindIndex(l, index)
	if err != nil {
		return "", err
	}

	value, ok := n.Value.(string)
	if !ok {
		return "", ErrTypeCorruption
	}

	return value, nil
}

// LSet sets the list element at index to element.
//
// TODO: Optimisation, we don't necessarily need to start from the Start of the list. We can traverse it
// starting from the end, if the requested range is "nearer" the end.
func (m *InMemory) LSet(key string, index int, value string) error {
	l, err := m.listGetKey(key)
	if err != nil {
		return err
	}

	// Check the index bounds
	if index >= l.Len() || (index < -l.Len()) {
		return server.ErrIndexOurOfRange
	}

	var n *list.Element
	if index >= 0 {
		n = l.Front()
		for i := 0; i < index; i++ {
			n = n.Next()
		}

	} else { // Negative numbers means to read from the end
		// Transform the index to a positive scale
		// [-3, -2, -1] is [0, 1, 2] in a positive scale. Hence, the -1.
		index = (-index) - 1

		n = l.Back()
		for i := 0; i < index; i++ {
			n = n.Prev()
		}
	}

	n.Value = value

	return nil
}

// LPop removes and returns the first elements of the list stored at key.
func (m *InMemory) LPop(key string) (string, error) {
	l, err := m.listGetKey(key)
	if err != nil {
		return "", err
	}

	first, err := m.listGetFirst(l)
	if err != nil {
		return "", err
	}

	v, err := m.listReadValue(first)
	if err != nil {
		return "", err
	}

	l.Remove(first)
	m.saveList(key, l)

	return v, nil
}

// RPop removes and returns the last elements of the list stored at key.
func (m *InMemory) RPop(key string) (string, error) {
	l, err := m.listGetKey(key)
	if err != nil {
		return "", err
	}

	last, err := m.listGetLast(l)
	if err != nil {
		return "", err
	}

	v, err := m.listReadValue(last)
	if err != nil {
		return "", err
	}

	l.Remove(last)
	m.saveList(key, l)

	return v, nil
}

// LPush insert all the specified values at the head of the list stored at key.
func (m *InMemory) LPush(key string, values []string) (int, error) {
	l, err := m.listGetKeyOrNew(key)
	if err != nil {
		return 0, err
	}

	for _, value := range values {
		l.PushFront(value)
	}

	m.saveList(key, l)

	return len(values), nil
}

// RPush insert all the specified values at the tail of the list stored at key.
func (m *InMemory) RPush(key string, values []string) (int, error) {
	l, err := m.listGetKeyOrNew(key)
	if err != nil {
		return 0, err
	}

	for _, value := range values {
		l.PushBack(value)
	}

	m.saveList(key, l)

	return len(values), nil
}

// LLen returns the length of the list stored at key
func (m *InMemory) LLen(key string) (int, error) {
	if err := m.assertType(key, listKind); err != nil {
		return 0, err
	}

	a, ok := m.records[key]
	if !ok {
		return 0, nil // List does not exist? Return it  has 0 elements.
	}

	l, err := a.List()
	if err != nil {
		return 0, err
	}

	return l.Len(), nil
}

func (m *InMemory) listFindIndex(l *list.List, index int) (*list.Element, error) {
	// Check the index bounds
	if index >= l.Len() || (index < -l.Len()) {
		return nil, server.ErrIndexOurOfRange
	}

	var n *list.Element
	if index >= 0 {
		n = l.Front()
		for i := 0; i < index; i++ {
			n = n.Next()
		}

	} else { // Negative numbers mean to read from the end
		// Transform the index to a positive scale
		// [-3, -2, -1] is [0, 1, 2] in a positive scale. Hence, the -1.
		index = (-index) - 1

		n = l.Back()
		for i := 0; i < index; i++ {
			n = n.Prev()
		}
	}

	return n, nil
}

func (m *InMemory) listGetKeyOrNew(key string) (*list.List, error) {
	l, err := m.listGetKey(key)
	if errors.Is(err, server.ErrNotFound) {
		return list.New(), nil
	}
	return l, err
}

func (m *InMemory) listGetKey(key string) (*list.List, error) {
	if err := m.assertType(key, listKind); err != nil {
		return nil, err
	}

	a, ok := m.records[key]
	if !ok {
		return nil, server.ErrNotFound
	}

	return a.List()
}

func (m *InMemory) listReadValue(e *list.Element) (string, error) {
	v, ok := e.Value.(string)
	if !ok {
		return "", fmt.Errorf("%w: expecting list element to be string", ErrTypeCorruption)
	}

	return v, nil
}

func (m *InMemory) listGetFirst(l *list.List) (*list.Element, error) {
	first := l.Front()
	if first == nil {
		return nil, fmt.Errorf("%w: expecting list to be non-empty", ErrTypeCorruption)
	}

	return first, nil
}

func (m *InMemory) listGetLast(l *list.List) (*list.Element, error) {
	first := l.Back()
	if first == nil {
		return nil, fmt.Errorf("%w: expecting list to be non-empty", ErrTypeCorruption)
	}

	return first, nil
}

// saveList must be called after all the operations that remove elements from the list. If the list becomes empty
// we need to remove the key from the storage.
func (m *InMemory) saveList(key string, l *list.List) {
	if l.Len() == 0 {
		delete(m.records, key)
	}

	m.records[key] = atom{kind: listKind, value: l}
}
