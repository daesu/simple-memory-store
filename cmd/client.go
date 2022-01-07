package cmd

import (
	"bufio"
	"fmt"
	"os"
	"simple-memory-store/store"
	"strings"
)

// Client implements an interface to the client
type Client interface {
	Run()
	ParseInput(input string) *instruction
	ExecCommand(cmd instruction) (string, error)
}

// NewClient returns an interface to the store
func NewClient() Client {
	return &client{
		scanner: bufio.NewScanner(os.Stdin),
		store:   store.NewStore(),
	}
}

type client struct {
	scanner *bufio.Scanner
	store   store.Store
}

// Run reads user input, parses it, and attempts to execute it
func (client *client) Run() {
	for {
		fmt.Print("> ")
		// reads user input until \n by default
		client.scanner.Scan()
		// Holds the string that was scanned
		input := client.scanner.Text()
		if len(input) < 1 {
			continue
		}

		cmd := client.ParseInput(input)
		if cmd.Command == "QUIT" {
			if client.store.ActiveTransactions() > 0 {
				fmt.Printf("%d active transactions aborted\n", client.store.ActiveTransactions())
			}
			os.Exit(0)
		}

		res, err := client.ExecCommand(*cmd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		fmt.Print(res)
	}
}

// parseInput takes the user input and parses it to an instruction
func (client *client) ParseInput(input string) *instruction {
	// Parse the input, splitting the command, args, and any value(s).
	args := strings.Split(input, " ")

	// Set args based on length of received arguments
	var cmdArgs []string
	if len(args) > 2 {
		// Presume value can contain any valid ASCII including white space
		cmdArgs = []string{args[1], strings.Join(args[2:], " ")}
	} else if len(args) > 1 {
		cmdArgs = []string{args[1]}
	}

	return &instruction{
		Command: strings.ToUpper(args[0]),
		Args:    cmdArgs,
	}
}

// ArgumentError returns an error when the incorrect number of args are passed
type ArgumentError struct {
	commandName       string
	expectedArgsCount int
	actualArgsCount   int
}

func (e *ArgumentError) Error() string {
	return fmt.Sprintf("%s expects %d argument but received %d", e.commandName, e.expectedArgsCount, e.actualArgsCount)
}

func (client *client) ExecCommand(cmd instruction) (string, error) {
	switch cmd.Command {
	case "READ":
		if len(cmd.Args) != 1 {
			return "", &ArgumentError{commandName: cmd.Command, expectedArgsCount: 1, actualArgsCount: len(cmd.Args)}
		}
		res, err := client.store.Read(cmd.Args[0])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("'%s'\n", res), nil
	case "DELETE":
		if len(cmd.Args) != 1 {
			return "", &ArgumentError{commandName: cmd.Command, expectedArgsCount: 1, actualArgsCount: len(cmd.Args)}
		}
		if err := client.store.Delete(cmd.Args[0]); err != nil {
			return "", err
		}
	case "WRITE":
		if len(cmd.Args) != 2 {
			return "", &ArgumentError{commandName: cmd.Command, expectedArgsCount: 2, actualArgsCount: len(cmd.Args)}
		}
		client.store.Write(cmd.Args[0], cmd.Args[1])
	case "START":
		client.store.StartTransaction()
	case "COMMIT":
		if err := client.store.CommitTransaction(); err != nil {
			return "", err
		}
	case "ABORT":
		if err := client.store.AbortTransaction(); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("command '%s' not found", cmd.Command)
	}

	return fmt.Sprintf("OK\n"), nil
}
