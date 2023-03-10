package easy_ping

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

var (
	WarnColor = text.Colors{text.BgRed}
	LossCount = 0
)

func NewResult(f, t string) *Result {
	return &Result{
		File:  f,
		Type:  t,
		Table: &table.Table{},
	}
}

func (r *Result) GlobalConfigSet() error {
	warnTransformer := text.Transformer(func(val interface{}) string {
		if val.(float64) > 0 {
			// 统计丢包服务器总数
			return WarnColor.Sprintf("%v%%", val)
		}
		return fmt.Sprintf("%v%%", val)
	})

	config := []table.ColumnConfig{}
	for _, v := range []string{"ID", "IP", "Num", "PacketsRecv", "AvgRtt"} {
		config = append(config, table.ColumnConfig{
			Name:        v,
			Align:       text.AlignCenter,
			AlignFooter: text.AlignCenter,
			AlignHeader: text.AlignCenter,
		})
	}

	config = append(config, table.ColumnConfig{
		Name:         "PacketLoss",
		Align:        text.AlignCenter,
		AlignFooter:  text.AlignCenter,
		AlignHeader:  text.AlignCenter,
		Hidden:       false,
		Transformer:  warnTransformer,
		VAlign:       text.VAlignMiddle,
		VAlignFooter: text.VAlignTop,
		VAlignHeader: text.VAlignBottom,
		WidthMin:     6,
		WidthMax:     64,
	})

	r.Table.SetColumnConfigs(config)

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
		// Make Table
		r.Table.SetOutputMirror(os.Stdout)
		err := r.GlobalConfigSet()
		if err != nil {
			return err
		}
		r.Table.AppendHeader(table.Row{"ID", "IP", "Num", "PacketsRecv", "PacketLoss", "AvgRtt"})
		for k, v := range r.Output {
			if v.PacketLoss > 0 {
				LossCount += 1
			}
			r.Table.AppendRow([]interface{}{k + 1, v.IP, v.Num, v.PacketsRecv, v.PacketLoss, v.AvgRtt})
		}
		r.Table.AppendFooter(table.Row{"", "Total", "Total", "Total", "Total", LossCount}, table.RowConfig{AutoMerge: true})
		switch r.Type {

		case "csv":
			r.Table.RenderCSV()

		default:
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
