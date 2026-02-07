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

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/filesys"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/services"
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
	if root != nil && root.Info != nil {
		if info, ok := root.Info.(file.NodeInformable); ok {
			logrus.Debugf("[track] current root path: %q", info.GetAbsPath())
		} else {
			logrus.Debugf("[track] current root (info type %T, no path)", root.Info)
		}
	} else {
		logrus.Debugf("[track] current root: nil")
	}

	// Resolve "." and relative paths before any logic.
	if pathExp == "." || pathExp == "" {
		if isTrackGDriveByContext(pathExp) {
			cwd, _ := defaultStore.GetGDriveCwd()
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
	} else if isTrackGDriveByContext(pathExp) && !file.IsGDrivePath(pathExp) && pathExp[0] != '/' {
		// Other relative paths in gdrive context: resolve against Drive cwd.
		cwd, _ := defaultStore.GetGDriveCwd()
		pathExp = contextrepo.ResolveGDrivePath(cwd, pathExp)
		logrus.Debugf("[track] resolved relative gdrive path to: %q", pathExp)
	}

	var subTree *ds.TreeNode
	if file.IsGDrivePath(pathExp) || isTrackGDriveByContext(pathExp) {
		logrus.Debugf("[track] using gdrive track path=%q", pathExp)
		subTree, err = trackGDrive(pathExp)
	} else {
		logrus.Debugf("[track] using local track path=%q", pathExp)
		subTree, err = filesys.Track(pathExp)
	}
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
	if len(embeddedCredentials) == 0 {
		return nil, fmt.Errorf("no embedded credentials; rebuild with credentials.json for Drive tracking")
	}
	tokenPath, err := resolveTokenPath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(tokenPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("token not found; run \"PathTracer login\" first for Drive")
		}
		return nil, err
	}
	ctx := context.Background()
	svc, err := services.NewGDriveServiceFromTokenPath(ctx, tokenPath, embeddedCredentials)
	if err != nil {
		return nil, err
	}
	drivePath, recursive := filesys.NormalizeGDriveTrackPath(pathExp)
	return filesys.TrackGDrive(ctx, drivePath, recursive, svc)
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
  track .   track SubFolder   track SubFolder*`,
	Run: track,
}

func init() {
	rootCmd.AddCommand(trackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
