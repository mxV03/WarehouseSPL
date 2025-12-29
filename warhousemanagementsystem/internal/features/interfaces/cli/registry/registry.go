package registry

import (
	"context"
	"fmt"
	"os"
	"sort"
)

type Command struct {
	Name        string
	Usage       string
	Description string
	Run         func(ctx context.Context, args []string) error
}

var commands = map[string]Command{}

func Register(cmd Command) {
	if cmd.Name == "" || cmd.Run == nil {
		panic("cli command must have Name and Run")
	}
	if _, exists := commands[cmd.Name]; exists {
		panic("cli command already registered: " + cmd.Name)
	}
	commands[cmd.Name] = cmd
}

func Dispatch(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		PrintHelp()
		return nil
	}

	name := args[0]
	cmd, ok := commands[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", name)
		PrintHelp()
		return fmt.Errorf("unknown command: %s", name)
	}
	return cmd.Run(ctx, args[1:])
}

func PrintHelp() {
	fmt.Println("Warehouse CLI\n")
	fmt.Println("Commands:")
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		cmd := commands[name]
		fmt.Printf("  %-22s %s\n      %s\n", cmd.Name, cmd.Usage, cmd.Description)
	}
}
