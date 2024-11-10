package file

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
)

type NodeTagsPrinter interface {
	PrintNodeTags(list.Writer) error
}

type NodePrinter interface {
	PrintNode(list.Writer) error
}

func getCurNodeFromAbsPath(absPath string) (string, error) {
	temp := strings.Split(absPath, "/")
	if len(temp) == 0 {
		return "", fmt.Errorf("can't print path: %s", absPath)
	}
	return temp[len(temp)-1], nil
}

func (gn *GeneralNode) PrintNode(wr list.Writer) error {
	curNodeName, err := getCurNodeFromAbsPath(gn.AbsPath)
	if err != nil {
		return err
	}

	wr.AppendItem(curNodeName)

	return nil
}

func (gn *GeneralNode) PrintNodeTags(wr list.Writer) error {
	err := gn.PrintNode(wr)
	if err != nil {
		return err
	}

	wr.Indent()
	if len(gn.Tags) > 0 {
		wr.AppendItem("<tags>")
		wr.Indent()
	}

	for _, tag := range gn.Tags {
		wr.AppendItem(tag)
	}

	if len(gn.Tags) > 0 {
		wr.UnIndent()
	}
	wr.UnIndent()

	return nil
}
