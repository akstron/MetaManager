package utils

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

/*
	Helps in creating dir structures for testing
*/

type DirTest struct {
	DirName string
	Files   []string
	Dirs    []*DirTest
}

/*
No need to save states, so keep it functional
Returns: <PathToDir>, <error>
*/
func CreateDirStructure(root *DirTest) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	randDirName := filepath.Join(wd, "testing_"+strconv.Itoa(rand.Int()))

	topLevelDir := DirTest{
		DirName: randDirName,
		Dirs:    []*DirTest{root},
	}

	queueDirs := []DirTest{topLevelDir}

	for i := 0; i < len(queueDirs); i++ {
		err := os.Mkdir(queueDirs[i].DirName, 0755)
		if err != nil {
			return "", err
		}

		for _, file := range queueDirs[i].Files {
			_, err = os.Create(file)
			if err != nil {
				return "", err
			}
		}

		for _, dir := range queueDirs[i].Dirs {

		}
	}
}
