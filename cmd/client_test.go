package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	client := NewClient()

	res := client.ParseInput("WRITE k value with whitespace")
	assert.EqualValues(t, &instruction{Command: "WRITE", Args: []string{"k", "value with whitespace"}}, res)

	res = client.ParseInput("READ k value with whitespace")
	assert.EqualValues(t, &instruction{Command: "READ", Args: []string{"k", "value with whitespace"}}, res)

	res = client.ParseInput("DELETE k")
	assert.EqualValues(t, &instruction{Command: "DELETE", Args: []string{"k"}}, res)

	res = client.ParseInput("START")
	assert.EqualValues(t, &instruction{Command: "START", Args: []string(nil)}, res)

	res = client.ParseInput("COMMIT")
	assert.EqualValues(t, &instruction{Command: "COMMIT", Args: []string(nil)}, res)

	res = client.ParseInput("ABORT")
	assert.EqualValues(t, &instruction{Command: "ABORT", Args: []string(nil)}, res)

	res = client.ParseInput("QUIT")
	assert.EqualValues(t, &instruction{Command: "QUIT", Args: []string(nil)}, res)
}

func TestExecCommandSuccess(t *testing.T) {
	client := NewClient()

	writeCmd := &instruction{Command: "WRITE", Args: []string{"k", "value with whitespace"}}
	_, err := client.ExecCommand(*writeCmd)
	assert.NoError(t, err)

	readCmd := &instruction{Command: "READ", Args: []string{"k"}}
	_, err = client.ExecCommand(*readCmd)
	assert.NoError(t, err)

	deleteCmd := &instruction{Command: "DELETE", Args: []string{"k"}}
	_, err = client.ExecCommand(*deleteCmd)
	assert.NoError(t, err)

	startCmd := &instruction{Command: "START", Args: []string(nil)}
	_, err = client.ExecCommand(*startCmd)
	assert.NoError(t, err)

	commitCmd := &instruction{Command: "COMMIT", Args: []string(nil)}
	_, err = client.ExecCommand(*commitCmd)
	assert.NoError(t, err)
}

func TestExecCommandFail(t *testing.T) {
	client := NewClient()

	readCmd := &instruction{Command: "READ", Args: []string{"k"}}
	_, err := client.ExecCommand(*readCmd)
	assert.Error(t, err)

	deleteCmd := &instruction{Command: "DELETE", Args: []string{"k"}}
	_, err = client.ExecCommand(*deleteCmd)
	assert.Error(t, err)

	commitCmd := &instruction{Command: "COMMIT", Args: []string(nil)}
	_, err = client.ExecCommand(*commitCmd)
	assert.Error(t, err)
}
