package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/spf13/cobra"
)

// defaultStore is used by cobra commands and package-level getters.
var defaultStore filesys.ContextRepository

func init() {
	defaultStore = filesys.NewContextRepositoryImpl(nil)
}

// contextCmd is the parent for context-related commands.
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Set or show storage context",
	Long:  `Context is a named environment (e.g. local, gdrive, work). The current context is read from ` + filesys.ContextEnvVar + ` if set, otherwise from a file next to the executable. All context names are in contexts.json.`,
}

// contextSetCmd sets the current context to the given name.
var contextSetCmd = &cobra.Command{
	Use:   "set <name>",
	Short: "Set the current context by name",
	Long:  `Persists the context name to the context file (used when ` + filesys.ContextEnvVar + ` is unset). For the current shell only, use: export ` + filesys.ContextEnvVar + `=<name>.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runContextSet,
}

// contextCreateCmd creates a named context with a type (local or gdrive) and adds it to contexts.json.
var contextCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a named context (local or gdrive) and save it in contexts.json",
	Long:  `Creates a new context with the given name and type (--type local or --type gdrive). Name must be unique. Use "context set <name>" to switch to it.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runContextCreate,
}

// contextGetCmd prints the current context name.
var contextGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current context name",
	Long:  `Prints the current context (from ` + filesys.ContextEnvVar + ` or the context file). Prints nothing if no context is set.`,
	Args:  cobra.NoArgs,
	RunE:  runContextGet,
}

// contextListCmd lists all contexts in table form.
var contextListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contexts in table form",
	Long:  `Lists every context from contexts.json with name, type, and which one is current.`,
	Args:  cobra.NoArgs,
	RunE:  runContextList,
}

var contextCreateType string

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextSetCmd)
	contextCmd.AddCommand(contextCreateCmd)
	contextCmd.AddCommand(contextGetCmd)
	contextCmd.AddCommand(contextListCmd)
	contextCreateCmd.Flags().StringVarP(&contextCreateType, "type", "t", "", "Context type: local or gdrive (required)")
	if err := contextCreateCmd.MarkFlagRequired("type"); err != nil {
		fmt.Fprintln(os.Stderr, "context create: mark flag required:", err)
		os.Exit(1)
	}
}

func runContextSet(cmd *cobra.Command, args []string) error {
	err := defaultStore.SetCurrent(args[0])
	if err != nil {
		return err
	}
	name := strings.ToLower(strings.TrimSpace(args[0]))
	fmt.Printf("Context set to %q\n", name)
	fmt.Printf("For this shell only: export %s=%s\n", filesys.ContextEnvVar, name)
	return nil
}

func runContextCreate(cmd *cobra.Command, args []string) error {
	contextType := strings.ToLower(strings.TrimSpace(contextCreateType))
	err := defaultStore.Create(args[0], contextType)
	if err != nil {
		return err
	}
	fmt.Printf("Context %q (%s) created\n", strings.ToLower(strings.TrimSpace(args[0])), contextType)
	return nil
}

func runContextGet(cmd *cobra.Command, args []string) error {
	name, err := defaultStore.GetContext()
	if err != nil {
		return err
	}
	if name != "" {
		fmt.Println(name)
	}
	return nil
}

func runContextList(cmd *cobra.Command, args []string) error {
	entries, err := defaultStore.LoadContexts()
	if err != nil {
		return err
	}
	current, err := defaultStore.GetContext()
	if err != nil {
		return err
	}
	current = strings.ToLower(strings.TrimSpace(current))

	const nameCol, typeCol, currentCol = 20, 12, 8
	header := fmt.Sprintf("%-*s %-*s %-*s", nameCol, "NAME", typeCol, "TYPE", currentCol, "CURRENT")
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", nameCol+typeCol+currentCol+2))
	for _, e := range entries {
		cur := ""
		if e.Name == current {
			cur = "*"
		}
		fmt.Printf("%-*s %-*s %-*s\n", nameCol, e.Name, typeCol, e.Type, currentCol, cur)
	}
	if len(entries) == 0 {
		fmt.Println("  (no contexts; use 'context create <name> --type local|gdrive')")
	}
	return nil
}

// GetContexts returns all context entries from the default store. Returns nil, nil if the file does not exist.
func GetContexts() ([]filesys.ContextEntry, error) {
	return defaultStore.LoadContexts()
}

// GetContext returns the current context name from the default store.
func GetContext() (string, error) {
	return defaultStore.GetContext()
}

// GetContextType returns the type for the given context name from the default store.
func GetContextType(name string) (string, error) {
	return defaultStore.GetContextType(name)
}
