package data

import "os"

/*
Tag related functionalities are implemented here
*/
type TagManager struct {
	dataFilePath string
	trMg         *TreeManager
}

func NewTagManager(dataFilePath string) (*TagManager, error) {
	tgMg := &TagManager{
		dataFilePath: dataFilePath,
	}

	// Read data in bytes from dataFilePath and construct TreeManager
	content, err := os.ReadFile(dataFilePath)
	if err != nil {
		return nil, err
	}

	/*
		WARNING: This does not initialize Root member of TreeManager
		There can be consequences
	*/
	tgMg.trMg = &TreeManager{}

	err = tgMg.trMg.Load(content)
	if err != nil {
		return nil, err
	}

	return tgMg, nil
}
