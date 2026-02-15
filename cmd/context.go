package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heroku/self/MetaManager/internal/config"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/context"
	filesys "github.com/heroku/self/MetaManager/internal/repository/context"
	"github.com/heroku/self/MetaManager/internal/repository/tree"
	"github.com/heroku/self/MetaManager/internal/utils"
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

// contextDeleteCmd removes a context from contexts.json.
var contextDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a context by name, or all contexts with --all",
	Long:  `Removes the context from contexts.json. If it was the current context, the current context is cleared. Use --all to delete every context.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runContextDelete,
}

var contextCreateType string

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextSetCmd)
	contextCmd.AddCommand(contextCreateCmd)
	contextCmd.AddCommand(contextGetCmd)
	contextCmd.AddCommand(contextListCmd)
	contextCmd.AddCommand(contextDeleteCmd)
	contextDeleteCmd.Flags().BoolP("all", "a", false, "Delete all contexts")
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
	name := strings.ToLower(strings.TrimSpace(args[0]))
	if err := EnsureAppDataDir(name); err != nil {
		return fmt.Errorf("ensure app data dir for context: %w", err)
	}
	fmt.Printf("Context %q (%s) created\n", name, contextType)
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

func runContextDelete(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	if all {
		if len(args) > 0 {
			return fmt.Errorf("cannot pass a name when using --all")
		}
		// Get all contexts before deleting so we can delete their .mm directories
		contexts, err := defaultStore.LoadContexts()
		if err != nil {
			return fmt.Errorf("load contexts: %w", err)
		}
		if err := defaultStore.DeleteAll(); err != nil {
			return err
		}
		// Delete .mm directories for all contexts
		for _, ctx := range contexts {
			if err := deleteContextMMDir(ctx.Name); err != nil {
				// Log error but continue deleting other contexts
				fmt.Fprintf(os.Stderr, "Warning: failed to delete .mm directory for context %q: %v\n", ctx.Name, err)
			}
		}
		fmt.Println("All contexts deleted")
		return nil
	}
	if len(args) != 1 {
		return fmt.Errorf("context name required (or use --all to delete all contexts)")
	}
	name := strings.ToLower(strings.TrimSpace(args[0]))
	if err := defaultStore.Delete(name); err != nil {
		return err
	}
	// Delete the .mm directory for this context
	if err := deleteContextMMDir(name); err != nil {
		return fmt.Errorf("delete .mm directory: %w", err)
	}
	fmt.Printf("Context %q deleted\n", name)
	return nil
}

// deleteContextMMDir deletes the .mm/<contextName>/ directory for the given context.
func deleteContextMMDir(contextName string) error {
	appDir, err := utils.GetAppDataDirForContext(contextName)
	if err != nil {
		return err
	}
	exists, err := utils.IsFilePresent(appDir)
	if err != nil {
		return err
	}
	if !exists {
		// Directory doesn't exist, nothing to delete
		return nil
	}
	return os.RemoveAll(appDir)
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

// getContextRequired returns the current context name, or error if unset (for commands that need a context).
func getContextRequired() (string, error) {
	name, err := GetContext()
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", fmt.Errorf("no context set; use 'context set <name>' first")
	}
	return name, nil
}

// EnsureAppDataDir creates the .mm/<contextName> directory for the given context (next to the executable
// or under MM_TEST_CONTEXT_DIR) with config.json and data.json if it does not exist. Idempotent.
func EnsureAppDataDir(contextName string) error {
	if contextName == "" {
		return fmt.Errorf("context name cannot be empty")
	}
	parentDir, err := utils.GetAppDataDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return err
	}
	appDir, err := utils.GetAppDataDirForContext(contextName)
	if err != nil {
		return err
	}
	exists, err := utils.IsFilePresent(appDir)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if err := os.Mkdir(appDir, 0755); err != nil {
		return err
	}
	baseDir, err := utils.GetBaseDir()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(appDir, utils.ConfigFileName)
	cfg := config.Config{RootPath: baseDir}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.WriteFile(configFilePath, data, 0666); err != nil {
		return err
	}

	absPath := "/"
	contextType, err := GetContextType(contextName)
	if err == nil && contextType == contextrepo.TypeGDrive {
		absPath = file.GDrivePathPrefix
	} else if os.Getenv("MM_TEST_CONTEXT_DIR") != "" {
		// In tests, root must match the track root so merge creates the right number of nodes.
		absPath = baseDir
	}
	emptyRoot := &ds.TreeNode{
		// root is a special path that is used to represent the root of the tree
		Info:     &file.FileNode{GeneralNode: file.GeneralNode{AbsPath: absPath}},
		Children: nil,
	}
	dataFilePath := filepath.Join(appDir, utils.DataFileName)
	rw, err := tree.NewFileStorageRW(dataFilePath)
	if err != nil {
		return err
	}
	return rw.Write(emptyRoot)
}
