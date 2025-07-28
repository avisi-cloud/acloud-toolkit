package table

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func Print(header []string, body [][]string) {
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Borders: tw.BorderNone,
			Settings: tw.Settings{
				Separators: tw.SeparatorsNone,
				Lines:      tw.LinesNone,
			},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Padding:    tw.CellPadding{Global: tw.Padding{Right: "    "}},
				Formatting: tw.CellFormatting{Alignment: tw.AlignLeft},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapNone},
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
				Padding:    tw.CellPadding{Global: tw.Padding{Right: tw.PaddingDefault.Right}},
			},
		}),
	)

	table.Header(header)
	table.Bulk(body)
	table.Render()
}
