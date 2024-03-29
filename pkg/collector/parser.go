package collector

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/EMnify/spu-exporter/pkg/transport"
)

var StringFind = regexp.MustCompile(`\s*([a-z-]+) "?([a-zA-Z0-9-\.]+)"?`)
var IntFind = regexp.MustCompile(`\s*([a-z-]+)\s(\d+)(\s{)*$`)

func parseLines(lines []string) ([]transport.Transport, error) {
	transportPattern := regexp.MustCompile(`transport (\d+)`)
	currentTransport := transport.Transport{}
	var trans []transport.Transport

	for _, line := range lines {
		tp := transportPattern.FindStringSubmatch(line)

		if tp != nil {
			if currentTransport.Number != nil {
				currentTransport.Peers = append(currentTransport.Peers, currentTransport.CurrentPeer)
				currentTransport.CurrentPeer = transport.Peer{}
				trans = append(trans, currentTransport)
			}
			i64, err := strconv.ParseInt(tp[1], 10, 64)
			if err != nil {
				return nil, err
			}
			currentTransport = transport.NewTransport(i64)

		} else {
			if currentTransport.Number != nil {
				ParseTransport(&currentTransport, line)
			} else {
				return nil, nil
			}
		}
	}
	currentTransport.Peers = append(currentTransport.Peers, currentTransport.CurrentPeer)
	currentTransport.CurrentPeer = transport.Peer{}
	trans = append(trans, currentTransport)

	return trans, nil
}

func ParseTransport(t *transport.Transport, line string) {
	n := IntFind.FindStringSubmatch(line)
	if n != nil {
		val, _ := strconv.ParseInt(n[2], 10, 64)
		switch n[1] {
		case "send-buffer":
			t.SendBuffer = val
		case "receive-buffer":
			t.ReceiveBuffer = val
		case "peer":
			if t.CurrentPeer.Number != nil {
				t.Peers = append(t.Peers, t.CurrentPeer)
			}
			t.CurrentPeer = transport.NewPeer(val)
		case "local-port":
			t.LocalPort = val
		}
		if t.CurrentPeer.Number != nil {
			ParsePeer(&t.CurrentPeer, line)
		}
		return
	}
	str := StringFind.FindStringSubmatch(line)
	if str != nil {

		switch str[1] {
		case "origin-host":
			t.OriginHost = str[2]
		case "origin-realm":
			t.OriginRealm = str[2]
		case "protocol":
			t.Protocol = str[2]
		case "local-ip":
			t.LocalIP = str[2]

		}
		if t.CurrentPeer.Number != nil {
			ParsePeer(&t.CurrentPeer, line)
		}

		t.LastKey = str[1]
		return
	}
	if strings.Contains(line, "client") {
		t.CurrentPeer = transport.NewPeer(0)
	}
	if strings.Contains(line, "{") {
		asdf := regexp.MustCompile("[a-z-]+")
		match := asdf.FindStringSubmatch(line)
		if match != nil {
			t.LastKey = match[0]
		}
	}
	if t.LastKey == "host-ips" {
		t.HostIps = append(t.HostIps, strings.TrimLeft(line, " "))
	}
	if t.LastKey == "applications" {
		t.Applications = append(t.Applications, strings.TrimLeft(line, " "))
	}
}

func ParsePeer(p *transport.Peer, line string) {

	//fmt.Println("parsing inside peer")
	n := IntFind.FindStringSubmatch(line)
	if n != nil {
		val, _ := strconv.ParseInt(n[2], 10, 64)
		switch n[1] {
		case "recv-cnt":
			p.Statistics.ReceiveCnt = val
		case "recv-max":
			p.Statistics.ReceiveMax = val
		case "recv-avg":
			p.Statistics.ReceiveAvg = val
		case "recv-oct":
			p.Statistics.ReceiveOct = val
		case "recv-dvi":
			p.Statistics.ReceiveDvi = val
		case "send-cnt":
			p.Statistics.SendCnt = val
		case "send-max":
			p.Statistics.SendMax = val
		case "send-avg":
			p.Statistics.SendAvg = val
		case "send-oct":
			p.Statistics.SendOct = val
		case "send-pend":
			p.Statistics.SendPending = val
		case "remote-port":
			p.RemotePort = val
		}

		return
	}
	str := StringFind.FindStringSubmatch(line)
	if str != nil {
		switch str[1] {
		case "destination-host":
			p.DestinationHost = str[2]
		case "destination-realm":
			p.DestinationRealm = str[2]
		case "remote-ip":
			p.RemoteIP = str[2]
		case "state":
			p.State.Name = str[2]
		}
		return
	}
}
