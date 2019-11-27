package util

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type Printer struct {
}

type PrinterInterface interface {
	PrintWithColumns([][]string, []string)
}

func (pr *Printer) PrintWithColumns(out [][]string, columns []string) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(columns)
	table.SetBorders(tablewriter.Border{
		Left:   false,
		Top:    false,
		Right:  false,
		Bottom: false,
	})
	table.SetCenterSeparator(" ")

	table.AppendBulk(out)
	table.Render()
}
