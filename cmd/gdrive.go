package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/filesys"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/context"
	"github.com/heroku/self/MetaManager/internal/services"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var gdriveCmd = &cobra.Command{
	Use:   "gdrive",
	Short: "Google Drive commands (requires login)",
	Long:  `Commands that use your Google Drive. Run "PathTracer login" first to authenticate.`,
}

var gdriveListCmd = &cobra.Command{
	Use:   "list [path or folder-id]",
	Short: "List directory structure of Google Drive",
	Long:  `Lists files and folders at the given path or folder ID. Path is like local: "/" for root, "/FolderName", "/FolderName/SubFolder". You can also pass a Drive folder ID directly. Without arguments, lists root.`,
	RunE:  runGDriveList,
}

var gdrivePwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "Print current Google Drive working directory",
	Long:  `Prints the current Drive directory used for shell-style navigation and relative paths in "track".`,
	RunE:  runGDrivePwd,
}

var gdriveCdCmd = &cobra.Command{
	Use:   "cd [path]",
	Short: "Change current Google Drive working directory",
	Long:  `Set the current Drive directory for relative paths. With no argument, prints current directory. Use "/" for root, "/Folder" or "Folder" for a subfolder, ".." to go up.`,
	RunE:  runGDriveCd,
}

var gdriveLsCmd = &cobra.Command{
	Use:   "ls [path]",
	Short: "List current or given Drive directory",
	Long:  `Lists files and folders at the current Drive directory (see "gdrive pwd") or at the given path. Path is relative to current directory unless it starts with "/". With no argument, lists current directory.`,
	RunE:  runGDriveLs,
}

var gdriveGetLinkCmd = &cobra.Command{
	Use:   "link [file-path]",
	Short: "Get the Google Drive web link for a file",
	Long:  `Returns the web view link for a file at the given path. Path can be a full path like "/Folder/file.txt" or a file ID.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGDriveGetLink,
}

// contextLsCmd, contextCdCmd, contextPwdCmd run gdrive ls/cd/pwd when current context is gdrive.
var contextLsCmd = &cobra.Command{
	Use:   "ls [path]",
	Short: "List current or given directory (uses current context)",
	Long:  `When current context is gdrive, lists the current Drive directory or the given path. Use "context set <name>" to switch context.`,
	RunE:  runContextLs,
}

var contextPwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "Print current working directory (uses current context)",
	Long:  `When current context is gdrive, prints the current Drive directory.`,
	RunE:  runContextPwd,
}

var contextCdCmd = &cobra.Command{
	Use:   "cd [path]",
	Short: "Change current working directory (uses current context)",
	Long:  `When current context is gdrive, changes the current Drive directory. Use ".." or a path relative to current directory.`,
	RunE:  runContextCd,
}

func init() {
	RootCmd.AddCommand(gdriveCmd)
	gdriveCmd.AddCommand(gdriveListCmd)
	gdriveCmd.AddCommand(gdrivePwdCmd)
	gdriveCmd.AddCommand(gdriveCdCmd)
	gdriveCmd.AddCommand(gdriveLsCmd)

	// Create "get" subcommand
	var gdriveGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get information about Google Drive files",
	}
	gdriveGetCmd.AddCommand(gdriveGetLinkCmd)
	gdriveCmd.AddCommand(gdriveGetCmd)

	RootCmd.AddCommand(contextLsCmd)
	RootCmd.AddCommand(contextPwdCmd)
	RootCmd.AddCommand(contextCdCmd)
}

// isGDriveContext returns true if the current context type is gdrive.
func isGDriveContext() (bool, error) {
	name, err := GetContext()
	if err != nil || name == "" {
		return false, err
	}
	typ, err := GetContextType(name)
	if err != nil || typ != contextrepo.TypeGDrive {
		return false, err
	}
	return true, nil
}

func requireGDriveContext(cmdName string) error {
	ok, err := isGDriveContext()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%s requires gdrive context; use \"context set <name>\" (with a gdrive context) or \"gdrive %s\"", cmdName, cmdName)
	}
	return nil
}

func runContextLs(cmd *cobra.Command, args []string) error {
	if err := requireGDriveContext("ls"); err != nil {
		return err
	}
	return runGDriveLs(cmd, args)
}

func runContextPwd(cmd *cobra.Command, args []string) error {
	if err := requireGDriveContext("pwd"); err != nil {
		return err
	}
	return runGDrivePwd(cmd, args)
}

func runContextCd(cmd *cobra.Command, args []string) error {
	if err := requireGDriveContext("cd"); err != nil {
		return err
	}
	return runGDriveCd(cmd, args)
}

func runGDriveList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	svc, err := GetGDriveService(ctx)
	if err != nil {
		return err
	}

	// Default: list root. Arg can be path (contains "/" or empty) or folder ID.
	target := "/"
	if len(args) > 0 {
		target = args[0]
	}

	var entries []services.RootEntry
	var displayPath string
	if target == "" || target == "/" || strings.Contains(target, "/") {
		displayPath = target
		if displayPath == "" {
			displayPath = "/"
		}
		entries, err = svc.ListAtPath(ctx, strings.TrimRight(displayPath, "/"))
		if err != nil {
			return fmt.Errorf("list %q: %w", displayPath, err)
		}
	} else {
		// Treat as folder ID
		displayPath = target
		entries, err = svc.ListFolder(ctx, target)
		if err != nil {
			return fmt.Errorf("list folder %q: %w", target, err)
		}
	}

	fmt.Println(displayPath)
	fmt.Println("---")
	for _, e := range entries {
		if e.IsFolder {
			fmt.Printf("  %s/\n", e.Name)
		} else {
			fmt.Printf("  %s\n", e.Name)
		}
	}
	if len(entries) == 0 {
		fmt.Println("  (empty)")
	}
	return nil
}

func runGDrivePwd(cmd *cobra.Command, args []string) error {
	cwd, err := defaultStore.GetGDriveCwd()
	if err != nil {
		return err
	}
	if cwd == "" {
		cwd = "/"
	}
	fmt.Println(cwd)
	return nil
}

func runGDriveCd(cmd *cobra.Command, args []string) error {
	cwd, err := defaultStore.GetGDriveCwd()
	if err != nil {
		return err
	}
	if len(args) == 0 {
		if cwd == "" {
			cwd = "/"
		}
		fmt.Println(cwd)
		return nil
	}
	resolver := filesys.NewBasicResolver(defaultStore)
	resolvedPath, err := resolver.Resolve(args[0])
	if err != nil {
		return err
	}
	// Strip gdrive:/ prefix for API call
	resolved := strings.TrimPrefix(resolvedPath, file.GDrivePathPrefix)
	if resolved == "" {
		resolved = "/"
	}
	if err := defaultStore.SetGDriveCwd(resolved); err != nil {
		return err
	}
	fmt.Println(resolved)
	return nil
}

func runGDriveLs(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	svc, err := GetGDriveService(ctx)
	if err != nil {
		return err
	}

	cwd, err := defaultStore.GetGDriveCwd()
	if err != nil {
		return err
	}
	if cwd == "" {
		cwd = "/"
	}

	target := cwd
	if len(args) > 0 {
		resolver := filesys.NewBasicResolver(defaultStore)
		resolvedPath, err := resolver.Resolve(args[0])
		if err != nil {
			return err
		}
		// Strip gdrive:/ prefix for API call
		target = strings.TrimPrefix(resolvedPath, file.GDrivePathPrefix)
		if target == "" {
			target = "/"
		}
	}

	displayPath := strings.TrimRight(target, "/")
	if displayPath == "" {
		displayPath = "/"
	}
	entries, err := svc.ListAtPath(ctx, displayPath)
	if err != nil {
		return fmt.Errorf("list %q: %w", displayPath, err)
	}

	fmt.Println(displayPath)
	fmt.Println("---")
	for _, e := range entries {
		if e.IsFolder {
			fmt.Printf("  %s/\n", e.Name)
		} else {
			fmt.Printf("  %s\n", e.Name)
		}
	}
	if len(entries) == 0 {
		fmt.Println("  (empty)")
	}
	return nil
}

func runGDriveGetLink(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	svc, err := GetGDriveService(ctx)
	if err != nil {
		return err
	}

	resolver := filesys.NewBasicResolver(defaultStore)
	resolvedPath, err := resolver.Resolve(args[0])
	if err != nil {
		return err
	}

	// Strip gdrive:/ prefix for API call
	filePath := strings.TrimPrefix(resolvedPath, file.GDrivePathPrefix)
	if filePath == "" {
		return fmt.Errorf("invalid path: %q", args[0])
	}

	link, err := svc.GetFileLink(ctx, filePath)
	if err != nil {
		return fmt.Errorf("get link for %q: %w", filePath, err)
	}

	fmt.Println(link)

	// Open the link in the default browser
	if err := browser.OpenURL(link); err != nil {
		return fmt.Errorf("open browser: %w", err)
	}

	return nil
}
