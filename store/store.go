package store

import (
	"fmt"
	"sync"
)

// Store implements an interface to the store
type Store interface {
	StartTransaction()
	AbortTransaction() error
	CommitTransaction() error
	Read(key string) (string, error)
	Delete(key string) error
	Write(key, value string)
	ActiveTransactions() int
}

// NewStore returns an interface to the store
func NewStore() Store {
	return &store{
		data:         make(map[string]string),
		lock:         &sync.RWMutex{},
		transactions: &transactions{},
	}
}

type store struct {
	data         map[string]string
	lock         *sync.RWMutex
	transactions *transactions
}

type transactions struct {
	top   *transaction
	count int
}

type transaction struct {
	data map[string]string
	lock *sync.RWMutex
	next *transaction
}

// StartTransaction starts a transaction and set's it as the top of the stack
func (store *store) StartTransaction() {
	transaction := transaction{data: make(map[string]string), lock: &sync.RWMutex{}}
	transaction.next = store.transactions.top

	// Copy parent data to child transaction
	if store.transactions.top != nil {
		transaction.data = store.transactions.top.data
	}
	store.transactions.top = &transaction
	store.transactions.count++
}

// AbortTransaction aborts all changes made in the current transaction
func (store *store) AbortTransaction() error {
	currentTransaction := store.transactions.top
	if currentTransaction == nil {
		return fmt.Errorf("No active transaction")
	}

	store.transactions.top = store.transactions.top.next
	store.transactions.count--
	return nil
}

// CommitTransaction commits all changes to the transaction and any parent but not to grandparents+
func (store *store) CommitTransaction() error {
	currentTransaction := store.transactions.top
	if currentTransaction == nil {
		return fmt.Errorf("No active transaction")
	}

	for key, value := range currentTransaction.data {
		if currentTransaction.next != nil {
			currentTransaction.next.data[key] = value
		} else {
			store.transactions.top = nil
			store.Write(key, value)
		}
	}
	if store.transactions.top != nil {
		store.transactions.top = store.transactions.top.next
	}

	return nil
}

// ActiveTransactions returns the count of uncommitted active transactions
func (store *store) ActiveTransactions() int {
	return store.transactions.count
}

// NoSuchKeyError returns an error when the specified key is missing
type NoSuchKeyError struct {
	key string
}

func (e *NoSuchKeyError) Error() string {
	return fmt.Sprintf("Key not found: %s", e.key)
}

// Read returns a value for the specified key if it exists
func (store *store) Read(k string) (string, error) {
	currentTransaction := store.transactions.top

	if currentTransaction == nil {
		store.lock.RLock()
		defer store.lock.RUnlock()

		if val, ok := store.data[k]; ok {
			return val, nil
		}
		return "", &NoSuchKeyError{key: k}
	}

	currentTransaction.lock.RLock()
	defer currentTransaction.lock.RUnlock()

	if val, ok := currentTransaction.data[k]; ok {
		return val, nil
	}
	return "", &NoSuchKeyError{key: k}
}

// Delete removes a key:value pair if it exists
func (store *store) Delete(k string) error {
	currentTransaction := store.transactions.top

	if currentTransaction == nil {
		if _, ok := store.data[k]; ok {
			delete(store.data, k)
			return nil
		}
		return &NoSuchKeyError{key: k}
	}

	currentTransaction.lock.RLock()
	defer currentTransaction.lock.RUnlock()

	if _, ok := currentTransaction.data[k]; ok {
		delete(currentTransaction.data, k)
		return nil
	}
	return &NoSuchKeyError{key: k}
}

// Write upserts a value to the specified key
// Existing key:value pairs will be overwritten
func (store *store) Write(k string, v string) {
	currentTransaction := store.transactions.top
	if currentTransaction == nil {
		store.lock.Lock()
		defer store.lock.Unlock()

		store.data[k] = v
		return
	}

	currentTransaction.lock.Lock()
	defer currentTransaction.lock.Unlock()
	currentTransaction.data[k] = v
}
