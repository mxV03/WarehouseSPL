package registry

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Command struct {
	Name        string
	Usage       string
	Description string
	Group       string
	Aliases     []string
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

	for _, alias := range cmd.Aliases {
		if alias == "" {
			continue
		}
		if _, exists := commands[alias]; exists {
			panic("cli command alias already registered: " + alias)
		}
		aliasCmd := cmd
		aliasCmd.Name = alias
		commands[alias] = aliasCmd
	}
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
	type row struct {
		name        string
		usage       string
		description string
		aliases     []string
	}

	type sig struct {
		run   string
		usage string
		desc  string
		group string
	}
	seen := map[sig]*row{}

	for name, cmd := range commands {
		if cmd.Run == nil {
			continue
		}
		s := sig{
			run:   fmt.Sprintf("%p", cmd.Run),
			usage: cmd.Usage,
			desc:  cmd.Description,
			group: cmd.Group,
		}

		if _, ok := seen[s]; !ok {
			r := &row{
				name:        name,
				usage:       cmd.Usage,
				description: cmd.Description,
				aliases:     []string{},
			}
			seen[s] = r
		} else {
			seen[s].aliases = append(seen[s].aliases, name)
		}

	}

	groups := map[string][]row{}
	for _, pr := range seen {
		group := "Commands"
		g := commands[pr.name].Group
		if g != "" {
			group = g
		}

		cleanAliases := make([]string, 0, len(pr.aliases))
		for _, a := range pr.aliases {
			if a != pr.name {
				cleanAliases = append(cleanAliases, a)
			}
		}
		sort.Strings(cleanAliases)
		pr.aliases = cleanAliases

		groups[group] = append(groups[group], *pr)
	}

	groupNames := make([]string, 0, len(groups))
	for g := range groups {
		groupNames = append(groupNames, g)
	}
	sort.Strings(groupNames)

	maxName := 0
	maxUsage := 0

	for _, g := range groupNames {
		for _, r := range groups[g] {
			n := r.name
			if len(r.aliases) > 0 {
				n = n + " (" + strings.Join(r.aliases, ", ") + ")"
			}
			if len(n) > maxName {
				maxName = len(n)
			}
			if len(r.usage) > maxUsage {
				maxUsage = len(r.usage)
			}
		}
	}

	if maxName < 12 {
		maxName = 12
	}
	if maxUsage < 28 {
		maxUsage = 28
	}

	fmt.Println("Warehouse CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  app <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")

	for _, g := range groupNames {
		fmt.Println(g)
		fmt.Println(strings.Repeat("-", len(g)))

		rows := groups[g]
		sort.Slice(rows, func(i, j int) bool { return rows[i].name < rows[j].name })

		for _, r := range rows {
			name := r.name
			if len(r.aliases) > 0 {
				name = name + " (" + strings.Join(r.aliases, ", ") + ")"
			}
			fmt.Printf("  %-*s  %-*s  %s\n", maxName, name, maxUsage, r.usage, r.description)
		}
		fmt.Println()
	}

	fmt.Println("Help:")
	fmt.Println("  app --help")
}
