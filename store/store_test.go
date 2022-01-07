package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var KEY = "fruit"
var VALUE = "apple"

func TestWrite(t *testing.T) {
	store := NewStore()

	store.Write(KEY, VALUE)
}

func TestRead(t *testing.T) {
	store := NewStore()

	// Invalid key
	_, err := store.Read("null")
	assert.Error(t, err)

	// Valid key
	store.Write(KEY, VALUE)
	res, err := store.Read(KEY)
	assert.NoError(t, err)
	assert.Equal(t, VALUE, res)
}

func TestDelete(t *testing.T) {
	store := NewStore()

	// Invalid key
	err := store.Delete("null")
	assert.Error(t, err)

	// Valid key
	store.Write(KEY, VALUE)
	err = store.Delete(KEY)
	assert.NoError(t, err)
}

func TestTransaction(t *testing.T) {
	store := NewStore()

	// Valid key
	store.Write("a", "hello")

	res, err := store.Read("a")
	assert.NoError(t, err)
	assert.Equal(t, "hello", res)

	// Parent Transaction
	store.StartTransaction()
	store.Write("a", "hello-again")

	res, err = store.Read("a")
	assert.NoError(t, err)
	assert.Equal(t, "hello-again", res)

	// Child Transaction
	store.StartTransaction()
	store.Delete("a")

	res, err = store.Read("a")
	assert.Error(t, err)

	// Commit Child Transaction
	err = store.CommitTransaction()
	assert.NoError(t, err)

	res, err = store.Read("a")
	assert.Error(t, err)

	store.Write("a", "once-more")

	res, err = store.Read("a")
	assert.NoError(t, err)
	assert.Equal(t, "once-more", res)

	// Abort Parent Transaction
	err = store.AbortTransaction()
	assert.NoError(t, err)

	res, err = store.Read("a")
	assert.NoError(t, err)
	assert.Equal(t, "hello", res)
}
