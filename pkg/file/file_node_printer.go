package file

import (
	"github/akstron/MetaManager/pkg/utils"

	"github.com/jedib0t/go-pretty/v6/list"
)

type GetPrintStringFunc func(info any) (string, error)

type NodeExtraInfoPrinter struct {
	GetPrintString GetPrintStringFunc
}

func NewNodeExtraInfoPrinter(f GetPrintStringFunc) *NodeExtraInfoPrinter {
	return &NodeExtraInfoPrinter{
		GetPrintString: f,
	}
}

func (np *NodeExtraInfoPrinter) GetPrinter(info any) (PrinterFunc, error) {
	return func(wr list.Writer) error {
		str, err := np.GetPrintString(info)
		if err != nil {
			return err
		}

		wr.AppendItem(str)

		return nil
	}, nil
}

type TagsPrinter interface {
	PrintTags(list.Writer) error
}

type NodePrinter interface {
	PrintNode(list.Writer) error
}

type IdPrinter interface {
	PrintId(list.Writer) error
}

type PrinterFunc func(list.Writer) error

type NodePrinterBuilder struct {
	printerFuncs []PrinterFunc
}

func (nb *NodePrinterBuilder) AppendPrinter(printerFunc func(list.Writer) error) {
	nb.printerFuncs = append(nb.printerFuncs, printerFunc)
}

func (nb *NodePrinterBuilder) Build() func(list.Writer) error {
	return func(wr list.Writer) error {
		for _, prFunc := range nb.printerFuncs {
			err := prFunc(wr)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (gn *GeneralNode) PrintNode(wr list.Writer) error {
	curNodeName, err := utils.GetCurNodeFromAbsPath(gn.AbsPath)
	if err != nil {
		return err
	}

	wr.AppendItem(curNodeName)

	return nil
}

func (gn *GeneralNode) PrintTags(wr list.Writer) error {
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

func (gn *GeneralNode) PrintId(wr list.Writer) error {
	if gn.Id != "" {
		wr.Indent()
		wr.AppendItem("id: " + gn.Id)
		wr.UnIndent()
	}

	return nil
}
