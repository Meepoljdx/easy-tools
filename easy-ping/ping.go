package main

import (
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/go-ping/ping"
)

type Ping struct {
	IP                    string `json:"IP"`
	Size                  int
	Num                   int     `json:"Num"`
	PacketsRecv           int     `json:"PacketsRecv"`
	PacketsSent           int     `json:"PacketsSent"`
	PacketsRecvDuplicates int     `json:"PacketsRecvDuplicates"`
	PacketLoss            float64 `json:"PacketLoss"`
	Addr                  string
	MinRtt                time.Duration `json:"MinRtt"`
	MaxRtt                time.Duration `json:"MaxRtt"`
	AvgRtt                time.Duration `json:"AvgRtt"`
	StdDevRtt             time.Duration `json:"StdDevRtt"`
}

func (p *Ping) sentIcmp() (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(p.IP)
	if err != nil {
		return nil, err
	}
	pinger.Count = p.Num
	pinger.Interval = time.Millisecond * 100
	pinger.Timeout = time.Second * 5
	pinger.SetPrivileged(true) // 使用icmp

	if err := pinger.Run(); err != nil {
		return nil, err
	}

	return pinger.Statistics(), nil
}

func NewPing(ip string, num int, size int) *Ping {
	return &Ping{
		IP:   ip,
		Num:  num,
		Size: size,
	}
}

func (p *Ping) Run() error {
	s, err := p.sentIcmp()
	if err != nil {
		return err
	}
	err = copyToResult(*s, p)
	if err != nil {
		return err
	}

	return nil
}

func copyToResult(s ping.Statistics, p *Ping) error {
	if err := convertor.CopyProperties(p, &s); err != nil {
		return err
	}

	return nil
}
