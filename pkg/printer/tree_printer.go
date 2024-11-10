package printer

import (
	"errors"
	"fmt"
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/file"

	"github.com/jedib0t/go-pretty/v6/list"
)

type ListPrintable interface {
	// Pass a list.Writer and implementation should
	// append the info in the form of list to it
	// Use AppendItem() of list.Writer
	// Internal printing indentation is the responsibility
	// of the implementation
	Print(list.Writer)
}

/*
	THE INFO IN THE TREE NODE SHOULD BE LISTPRINTABLE
	WITH CUSTOM PRINT IMPLEMENTATION ACCORDINGLY

	FOR EXAMPLE: THERE CAN BE A PRINTER WHICH ONLY PRINTS NODE
	AND A PRINTER WHICH PRINTS TAGS AS WELL

*/

type TreePrinterManager struct {
	trMg *ds.TreeManager
	wr   list.Writer
}

func NewTreePrinterManager(trMg *ds.TreeManager) *TreePrinterManager {
	return &TreePrinterManager{
		trMg: trMg,
		wr:   list.NewWriter(),
	}
}

func getPrinter(ty string, info any) (func(list.Writer) error, error) {
	var resFunc func(list.Writer) error

	/*
		TODO: All the code inside each cases, can be
		encapsulated in a PrinterGenerator func and this
		can be stored for each "ty" in a map

		But should we do this for a cmd line tool?
	*/
	switch ty {
	case "node":
		printer, ok := info.(file.NodePrinter)
		if !ok {
			return nil, errors.New("info not convertible to NodePrinter")
		}
		resFunc = printer.PrintNode
	case "node-tags":
		printer, ok := info.(file.NodeTagsPrinter)
		if !ok {
			return nil, errors.New("info not convertible to NodeTagsPrinter")
		}
		resFunc = printer.PrintNodeTags
	default:
		return nil, errors.New("unimplemented")
	}

	return resFunc, nil
}

func (mg *TreePrinterManager) TrPrint(ty string) error {
	err := mg.trPrint(ty, mg.trMg.Root)
	if err != nil {
		return err
	}

	mg.wr.SetStyle(list.StyleConnectedLight)
	fmt.Println(mg.wr.Render())

	return nil
}

func (pr *TreePrinterManager) trPrint(ty string, curNode *ds.TreeNode) error {
	printFunc, err := getPrinter(ty, curNode.Info)
	if err != nil {
		return err
	}

	err = printFunc(pr.wr)
	if err != nil {
		return err
	}

	pr.wr.Indent()

	for _, child := range curNode.Children {
		err := pr.trPrint(ty, child)
		if err != nil {
			return err
		}
	}

	pr.wr.UnIndent()

	return nil
}