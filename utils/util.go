package utils

import (
	"bufio"
	"io"
	"net"
	"os"
)

func ReadIPFromFile(filename string) ([]string, error) {
	var ipList []string
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(f)
	for {
		str, _, err := buf.ReadLine()
		if err != nil || string(str) == "" {
			if err == io.EOF {
				break
			}
			continue
		}
		ipList = append(ipList, string(str))
	}
	return ipList, err
}

func FileExisted(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetBondIP() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	address, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}
