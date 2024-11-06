package utils

import (
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/file"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
)

type ListPrintable interface {
	// Pass a list.Writer and implementation should
	// append the info in the form of list to it
	// Use AppendItem() of list.Writer
	Print(list.Writer)
}

type ListNodeAndTagsPrinter struct {
}

func (ListNodeAndTagsPrinter) Print(list.Writer) {

}

/*
	THE INFO IN THE TREE NODE SHOULD BE LISTPRINTABLE
	WITH CUSTOM PRINT IMPLEMENTATION ACCORDINGLY

	FOR EXAMPLE: THERE CAN BE A PRINTER WHICH ONLY PRINTS NODE
	AND A PRINTER WHICH PRINTS TAGS AS WELL

*/

type TreePrinterManager struct {
	*ds.TreeManager
}

func (mg *TreePrinterManager) TrPrint() error {
	return mg.trPrint(mg.Root)
}

func (*TreePrinterManager) trPrint(curNode *ds.TreeNode) error {

}

func ConstructTreeWriter(curNode *ds.TreeNode, cutPrefix string, wr list.Writer) error {
	info := curNode.Info.(file.NodeInformable)
	insPath, _ := strings.CutPrefix(info.GetAbsPath(), cutPrefix+"/")
	wr.AppendItem(insPath)

	wr.Indent()

	for _, child := range curNode.Children {
		err := ConstructTreeWriter(child, info.GetAbsPath(), wr)
		if err != nil {
			return err
		}
	}

	wr.UnIndent()
	return nil
}
