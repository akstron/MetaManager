package data

/*
- Manages the nodes

DirPath
  - The parent root

root
  - Root of the scanned nodes. DirPath will be scanned
*/

// type TreeManager struct {
// 	DirPath string
// 	Root    *DirNode
// }

type TreeManager struct {
	Root *TreeNode
}

func NewTreeManager(root *TreeNode) *TreeManager {
	return &TreeManager{
		Root: root,
	}
}

/*
The tree node must contain
*/
type DirTreeManager struct {
	tgMg TreeManager
}

func (mg *TreeManager) FindNodeByAbsPath(path string) (NodeInformable, error) {
	ti := NewTreeIterator(mg)
	return mg.findNodeByAbsPathInternal(ti, path)
}

func (mg *TreeManager) findNodeByAbsPathInternal(it TreeIterator, path string) (NodeInformable, error) {
	for it.HasNext() {
		got, err := it.Next()
		if err != nil {
			return nil, err
		}

		if got.(NodeInformable).GetAbsPath() == path {
			return got.(NodeInformable), nil
		}
	}

	return nil, nil
}

/*
TODO: Remove this from TreeManager as this shouldn't be
aware of directory structure. Need to think about decoupling it.
// */
// func (mg *TreeManager) ScanDirectory() error {
// 	present, err := utils.IsFilePresent(mg.DirPath)
// 	if err != nil {
// 		return err
// 	}

// 	if !present {
// 		return &cmderror.InvalidPath{}
// 	}

// 	dirPathAbs, err := filepath.Abs(mg.DirPath)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println(dirPathAbs)

// 	root, err := os.Stat(dirPathAbs)
// 	if err != nil {
// 		return err
// 	}

// 	mg.Root = &DirNode{GeneralNode: GeneralNode{entry: root, absPath: dirPathAbs}}

// 	igMg, err := config.NewIgnoreManager()
// 	if err != nil {
// 		return err
// 	}

// 	igMg.Load()
// 	ignorer := &NodeAbsPathIgnorer{
// 		igMg: igMg,
// 	}

// 	/*
// 		TODO: Probably better to pass pointer
// 		We will use go routines to accelerate dfs
// 	*/
// 	handler := ScanHandler{
// 		ig: ignorer,
// 	}
// 	err = mg.Root.Scan(handler)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
