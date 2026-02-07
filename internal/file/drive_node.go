package file

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/heroku/self/MetaManager/internal/ds"
)

// driveNodeJSON is the persisted shape for Drive nodes.
type driveNodeJSON struct {
	Parent  string   `json:"Parent"`
	DriveId string   `json:"DriveId"`
	Tags    []string `json:"Tags"`
	Id      string   `json:"Id"`
}

// GDrivePathPrefix is the virtual path prefix for Google Drive nodes (e.g. "gdrive:/Folder/file").
const GDrivePathPrefix = "gdrive:/"

// DriveNodeBase holds fields common to Drive file and dir nodes (implements NodeInformable).
type DriveNodeBase struct {
	AbsPath string   // Virtual path, e.g. "gdrive:/My Folder/doc.pdf"
	DriveId string   // Google Drive file ID for API calls
	Tags    []string
	Id      string   // User-facing id
}

func (d *DriveNodeBase) GetAbsPath() string { return d.AbsPath }
func (d *DriveNodeBase) GetTags() []string   { return d.Tags }
func (d *DriveNodeBase) SetId(id string)     { d.Id = id }
func (d *DriveNodeBase) GetId() string       { return d.Id }

func (d *DriveNodeBase) AddTag(tag string) {
	for _, t := range d.Tags {
		if t == tag {
			return
		}
	}
	d.Tags = append(d.Tags, tag)
}

func (d *DriveNodeBase) DeleteTag(tag string) {
	var out []string
	for _, t := range d.Tags {
		if t != tag {
			out = append(out, t)
		}
	}
	d.Tags = out
}

// DriveDirNode represents a folder in Google Drive.
type DriveDirNode struct {
	DriveNodeBase
}

func (d *DriveDirNode) GetInfoProvider() NodeInformable { return d }
func (d *DriveDirNode) Name() string                    { return "GDRIVE_DIR" }

func (d *DriveDirNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(driveNodeJSON{Parent: d.AbsPath, DriveId: d.DriveId, Tags: d.Tags, Id: d.Id})
}
func (d *DriveDirNode) UnmarshalJSON(data []byte) error {
	var o driveNodeJSON
	if err := json.Unmarshal(data, &o); err != nil {
		return err
	}
	d.AbsPath = o.Parent
	d.DriveId = o.DriveId
	d.Tags = o.Tags
	d.Id = o.Id
	return nil
}

// DriveFileNode represents a file in Google Drive.
type DriveFileNode struct {
	DriveNodeBase
}

func (d *DriveFileNode) GetInfoProvider() NodeInformable { return d }
func (d *DriveFileNode) Name() string                    { return "GDRIVE_FILE" }

func (d *DriveFileNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(driveNodeJSON{Parent: d.AbsPath, DriveId: d.DriveId, Tags: d.Tags, Id: d.Id})
}
func (d *DriveFileNode) UnmarshalJSON(data []byte) error {
	var o driveNodeJSON
	if err := json.Unmarshal(data, &o); err != nil {
		return err
	}
	d.AbsPath = o.Parent
	d.DriveId = o.DriveId
	d.Tags = o.Tags
	d.Id = o.Id
	return nil
}

// IsGDrivePath returns true if path is a Drive virtual path.
func IsGDrivePath(path string) bool {
	return strings.HasPrefix(path, GDrivePathPrefix)
}

// NewDriveDirNode creates a tree node for a Drive folder.
func NewDriveDirNode(virtualPath, driveId string) *ds.TreeNode {
	return ds.NewTreeNode(&DriveDirNode{
		DriveNodeBase: DriveNodeBase{AbsPath: virtualPath, DriveId: driveId},
	})
}

// NewDriveFileNode creates a tree node for a Drive file.
func NewDriveFileNode(virtualPath, driveId string) *ds.TreeNode {
	return ds.NewTreeNode(&DriveFileNode{
		DriveNodeBase: DriveNodeBase{AbsPath: virtualPath, DriveId: driveId},
	})
}

// CreateTreeNodeFromPath creates a node from a path. For paths starting with GDrivePathPrefix
// it creates a Drive node without calling os.Stat (defaults to dir for path-segment creation).
func CreateTreeNodeFromPath(path string) (*ds.TreeNode, error) {
	if IsGDrivePath(path) {
		return ds.NewTreeNode(&DriveDirNode{
			DriveNodeBase: DriveNodeBase{AbsPath: path},
		}), nil
	}
	entry, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return CreateTreeNodeFromPathAndType(path, entry.IsDir())
}
