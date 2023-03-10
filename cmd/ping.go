package cmd

import (
	"fmt"
	"os"

	easy_ping "github.com/Meepoljdx/easy-tools/easy-ping"
	"github.com/Meepoljdx/easy-tools/utils"
	"github.com/spf13/cobra"
)

var (
	ip       []string
	packNum  int
	packSize int
	file     string
	output   string
	channl   int

	pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "Easy to ping ip",
		Long: `The easy-ping can send icmp packets to multiple server and print the result to users.
Users can choose json, csv, excel as the export format`,
		Run: func(cmd *cobra.Command, args []string) {
			ping()
		},
	}
)

func init() {
	pingCmd.Flags().StringSliceVarP(&ip, "ip", "i", []string{"127.0.0.1"}, "If set ip, you can use ip1,ip2,ip3 to specify the server on which the ping test is to be performed.")
	pingCmd.Flags().IntVarP(&packNum, "packet-num", "n", 10, "The num of packets will be send to remote server.")
	pingCmd.Flags().IntVarP(&packSize, "packet-size", "s", 64, "The size of packets will be send to remote server.")
	pingCmd.Flags().StringVarP(&file, "file", "f", "", "The ip flag will be ignored if file has been set, you can write all ip which the ping test is to be performed in a file.")
	pingCmd.Flags().StringVarP(&output, "output", "o", "stdout", "The output of the result, you can set stdout, csv, json, excel.")
	pingCmd.Flags().IntVarP(&channl, "channle", "c", 30, "Number of pings performed at the same time")
	pingCmd.MarkFlagsMutuallyExclusive("ip", "file")
	// 添加到根命令下
	rootCmd.AddCommand(pingCmd)
}

func ping() {
	if file != "" && utils.FileExisted(file) {
		var err error
		ip, err = utils.ReadIPFromFile(file)
		if err != nil {
			fmt.Println("Read ip from file failed.")
		}
	}
	o := easy_ping.ServerPing(ip, output, channl, packNum, packSize)
	if err := o.ResultOutPut(); err != nil {
		fmt.Fprint(os.Stderr, err)
	}

}
