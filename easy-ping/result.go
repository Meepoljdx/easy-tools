package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/xuri/excelize/v2"
)

type Result struct {
	File   string
	Type   string
	Output []Ping
	Lock   sync.Mutex
	Table  *table.Table
}

func NewResult(f, t string) *Result {
	return &Result{
		File:  f,
		Type:  t,
		Table: &table.Table{},
	}
}

func (r *Result) GlobalConfigSet() error {
	r.Table.SetColumnConfigs([]table.ColumnConfig{
		{
			Align:        text.AlignLeft,
			AlignFooter:  text.AlignLeft,
			AlignHeader:  text.AlignLeft,
			Colors:       text.Colors{text.BgBlack, text.FgRed},
			ColorsHeader: text.Colors{text.BgRed, text.FgBlack, text.Bold},
			ColorsFooter: text.Colors{text.BgRed, text.FgBlack},
			Hidden:       false,
			VAlign:       text.VAlignMiddle,
			VAlignFooter: text.VAlignTop,
			VAlignHeader: text.VAlignBottom,
			WidthMin:     6,
			WidthMax:     64,
		}})
	return nil
}

func (r *Result) ResultOutPut() error {
	fmt.Println("Output:")
	if r.Type == "json" {

	} else if r.Type == "excel" {
		now := time.Now().UnixMilli()
		excel := fmt.Sprintf("ping_%v.xlsx", now)
		return writeToExcel(excel, r.Output)
	} else {
		// 生成表格
		r.Table.SetOutputMirror(os.Stdout)

		r.Table.AppendHeader(table.Row{"ID", "IP", "Num", "PacketsRecv", "PacketLoss", "AvgRtt"})
		for k, v := range r.Output {
			// t.AppendRow([]interface{}{fmt.Sprintf("%v", k+1), v.IP, fmt.Sprintf("%v", v.Num), fmt.Sprintf("%v", v.PacketsRecv), fmt.Sprintf("%v%%", v.PacketLoss), fmt.Sprintf("%v", v.AvgRtt)})
			r.Table.AppendRow([]interface{}{k + 1, v.IP, v.Num, v.PacketsRecv, fmt.Sprintf("%v%%", v.PacketLoss), v.AvgRtt})
		}

		switch r.Type {

		case "csv":
			r.Table.RenderCSV()

		default:
			fmt.Println("Unknown output format or stdout. Will use stdout as output type.")
			r.Table.Render()

		}

	}
	return nil
}

func writeToExcel(excelFile string, data []Ping) error {
	f := excelize.NewFile()

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	_, err := f.NewSheet("Sheet1")
	if err != nil {
		return err
	}

	f.SetSheetRow("Sheet1", "A1", &[]interface{}{"IP", "数据包数量", "接收数据报数量", "丢包率", "平均响应时长"})
	rowLen := len(data)
	for i := 0; i < rowLen; i++ {
		row := data[i]
		if err := f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), &[]interface{}{
			row.IP,
			row.Num,
			row.PacketsRecv,
			fmt.Sprintf("%v%%", row.PacketLoss),
			row.AvgRtt.String(),
		}); err != nil {
			return err
		}
	}

	if err := f.SaveAs(excelFile); err != nil {
		fmt.Println("Failed to save excele file.")
	}

	abs, err := filepath.Abs(excelFile)
	if err != nil {
		return err
	}
	fmt.Printf("Create excel file success: %s", abs)
	return nil
}
