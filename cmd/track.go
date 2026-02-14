/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/filesys"
	"github.com/heroku/self/MetaManager/internal/printer"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/storage"
	"github.com/heroku/self/MetaManager/internal/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func trackInternal(ctxName, pathExp string) error {
	logrus.Debugf("[track] trackInternal start ctx=%q pathExp=%q", ctxName, pathExp)

	rw, err := storage.GetRW(ctxName)
	if err != nil {
		logrus.Debugf("[track] GetRW error: %v", err)
		return err
	}

	root, err := rw.Read()
	if err != nil {
		logrus.Debugf("[track] Read root error: %v", err)
		return err
	}

	if root == nil {
		return fmt.Errorf("root is nil")
	}
	info, ok := root.Info.(file.NodeInformable)
	if !ok {
		return fmt.Errorf("root info is not a NodeInformable")
	}
	logrus.Debugf("[track] current root path: %q", info.GetAbsPath())

	// Resolve "." and relative paths before any logic.
	if pathExp == "." || pathExp == "" {
		if isTrackGDriveByContext(pathExp) {
			cwd, _ := defaultStore.GetGDriveCwd()
			logrus.Debugf("[track] gdrive cwd: %q", cwd)
			pathExp = contextrepo.ResolveGDrivePath(cwd, ".")
			logrus.Debugf("[track] resolved . to gdrive cwd: %q", pathExp)
		} else {
			abs, err := filepath.Abs(".")
			if err != nil {
				return fmt.Errorf("resolve current directory: %w", err)
			}
			pathExp = abs
			logrus.Debugf("[track] resolved . to local path: %q", pathExp)
		}
	} else if isTrackGDriveByContext(pathExp) {
		// Other relative paths in gdrive context: resolve against Drive cwd.
		cwd, err := defaultStore.GetGDriveCwd()
		if err != nil {
			return fmt.Errorf("get gdrive cwd: %w", err)
		}
		pathExp = contextrepo.ResolveGDrivePath(cwd, pathExp)
		logrus.Debugf("[track] resolved relative gdrive path to: %q", pathExp)
	}

	tracker := filesys.NewContextAwareTracker(defaultStore)
	subTree, err := tracker.Track(pathExp)
	if err != nil {
		logrus.Debugf("[track] track (gdrive/local) error: %v", err)
		return err
	}

	logrus.Debugf("[track] merge subtree into root")
	drMg := data.NewDirTreeManager(ds.NewTreeManager(root))

	err = drMg.MergeNode(subTree)
	if err != nil {
		logrus.Debugf("[track] MergeNode error: %v", err)
		return err
	}

	err = rw.Write(drMg.Root)
	if err != nil {
		logrus.Debugf("[track] Write root error: %v", err)
		return err
	}

	logrus.Debugf("[track] trackInternal done")
	return nil
}

// trackShowInternal lists tracked nodes from the current directory (local cwd or gdrive cwd) in a tree structure.
func trackShowInternal(ctxName string, tagFlag, idFlag bool) error {
	rw, err := storage.GetRW(ctxName)
	if err != nil {
		return err
	}
	root, err := rw.Read()
	if err != nil {
		return err
	}
	var dirPath string
	ctxType, err := GetContextType(ctxName)
	if err == nil && ctxType == contextrepo.TypeGDrive {
		cwd, _ := defaultStore.GetGDriveCwd()
		resolved := contextrepo.ResolveGDrivePath(cwd, ".")
		if resolved == "/" {
			dirPath = file.GDrivePathPrefix
		} else {
			dirPath = file.GDrivePathPrefix + strings.TrimPrefix(resolved, "/")
		}
	} else {
		dirPath, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	drMg := data.NewDirTreeManager(ds.NewTreeManager(root))
	requiredNode, err := drMg.FindTreeNodeByAbsPath(dirPath)
	if err != nil {
		return err
	}
	pr := printer.NewTreePrinterManager(ds.NewTreeManager(requiredNode))
	typesOfPrinting := []string{"node"}
	if idFlag {
		typesOfPrinting = append(typesOfPrinting, "id")
	}
	if tagFlag {
		typesOfPrinting = append(typesOfPrinting, "tags")
	}
	return pr.TrPrint(typesOfPrinting)
}

// isTrackGDriveByContext returns true when current context is gdrive and path looks like a Drive path (starts with / or is a single segment).
func isTrackGDriveByContext(pathExp string) bool {
	name, err := GetContext()
	if err != nil || name == "" {
		return false
	}
	typ, err := GetContextType(name)
	if err != nil || typ != contextrepo.TypeGDrive {
		return false
	}
	// Path like "/Folder", "Folder", "/Folder/Sub", "Folder*"
	return len(pathExp) > 0 && (pathExp[0] == '/' || (pathExp[0] != '*' && pathExp != ""))
}

func trackGDrive(pathExp string) (*ds.TreeNode, error) {
	ctx := context.Background()
	svc, err := GetGDriveService(ctx)
	if err != nil {
		return nil, err
	}
	scanner := filesys.NewGDriveScanner(svc)
	drivePath, recursive := scanner.NormalizeTrackPath(pathExp)
	return scanner.TrackGDrive(ctx, drivePath, recursive)
}

func track(cmd *cobra.Command, args []string) {
	var err error
	var ctxName string

	logrus.Debugf("[track] track command args=%v", args)
	if len(args) != 1 {
		err = &cmderror.InvalidNumberOfArguments{}
		goto finally
	}

	ctxName, err = getContextRequired()
	if err != nil {
		logrus.Debugf("[track] getContextRequired error: %v", err)
		goto finally
	}
	logrus.Debugf("[track] context=%q", ctxName)
	_, err = utils.CommonAlreadyInitializedChecks(ctxName)
	if err != nil {
		logrus.Debugf("[track] CommonAlreadyInitializedChecks error: %v", err)
		goto finally
	}

	err = trackInternal(ctxName, args[0])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
		// Print stack trace in case of error
		debug.PrintStack()
	} else {
		fmt.Println("Location tracked successfully")
	}
}

// trackShowCmd lists tracked files/dirs from the current (local or gdrive) directory in a tree structure.
var trackShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show tracked files/dirs from current directory in a tree structure",
	Long:  "Lists all tracked files/dirs from the current root (local cwd or gdrive cwd) in a tree structure. Works in both local and gdrive contexts.",
	Run:   runTrackShow,
}

func runTrackShow(cmd *cobra.Command, args []string) {
	var err error
	var tagFlag, idFlag bool
	var ctxName string
	ctxName, err = getContextRequired()
	if err != nil {
		goto finally
	}
	_, err = utils.CommonAlreadyInitializedChecks(ctxName)
	if err != nil {
		goto finally
	}
	tagFlag, err = cmd.Flags().GetBool("tag")
	if err != nil {
		goto finally
	}
	idFlag, err = cmd.Flags().GetBool("id")
	if err != nil {
		goto finally
	}
	err = trackShowInternal(ctxName, tagFlag, idFlag)
	if err != nil {
		goto finally
	}
	return
finally:
	if err != nil {
		fmt.Println(err)
	}
}

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Starts tracking a path (local or Google Drive)",
	Long: `Starts tracking a file/folder or whole tree.

Local paths:
  track "/home/dev/project"
  track "/home/dev/project*"   (recursive; use quotes)

Google Drive (run "PathTracer login" first):
  track "gdrive:/"             (root of My Drive)
  track "gdrive:/Folder"
  track "gdrive:/Folder/Sub*"  (recursive)

With context set to gdrive, you can also use:
  track "/Folder"   or   track "/Folder*"

After "gdrive cd /SomeFolder", relative paths use that directory:
  track .   track SubFolder   track SubFolder*

Subcommands:
  track show   show tracked nodes from current directory (local or gdrive cwd)`,
	Run: track,
}

func init() {
	rootCmd.AddCommand(trackCmd)
	trackCmd.AddCommand(trackShowCmd)
	trackShowCmd.Flags().BoolP("tag", "t", false, "include tags for each node")
	trackShowCmd.Flags().BoolP("id", "i", false, "include id for each node")
}
