package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/urfave/cli"
)

// Packet ...
type Packet struct {
	Second    uint32 `json:"-"`
	Usecond   uint32 `json:"-"`
	Time      uint64
	TimeDelta uint64
	Length    uint32 `json:"-"`
	Data      string
	Filter    string `json:",omitempty"`
}

func filter(cmdname, s string) string {
	cmd := exec.Command("sh", "-c", cmdname)
	cmd.Stdin = bytes.NewBufferString(s)
	b, err := cmd.Output()
	if err != nil {
		log.Print("exec error", err)
		return s
	}
	return string(b[:])
}

func read1(infp io.Reader) (pkt Packet, err error) {
	if err = binary.Read(infp, binary.LittleEndian, &pkt.Second); err != nil {
		if err != io.EOF {
			log.Print("error reading header1", err)
		}
		return
	}
	if err = binary.Read(infp, binary.LittleEndian, &pkt.Usecond); err != nil {
		log.Print("error reading header2", err)
		return
	}
	pkt.Time = uint64(pkt.Second)*1000000 + uint64(pkt.Usecond)
	if err = binary.Read(infp, binary.LittleEndian, &pkt.Length); err != nil {
		log.Print("error reading header3", err)
		return
	}
	data := make([]byte, pkt.Length, pkt.Length)
	var rdlen int
	if rdlen, err = infp.Read(data); err != nil || rdlen != int(pkt.Length) {
		log.Print("error reading body", err, rdlen)
		return
	}
	pkt.Data = string(data[:])
	return
}

func write1(outfp io.Writer, pkt Packet) (err error) {
	if err = binary.Write(outfp, binary.LittleEndian, pkt.Second); err != nil {
		log.Print("error writing header1", err)
		return
	}
	if err = binary.Write(outfp, binary.LittleEndian, pkt.Usecond); err != nil {
		log.Print("error writing header2", err)
		return
	}
	if pkt.Filter != "" {
		pkt.Data = filter(pkt.Filter, pkt.Data)
		// log.Print("filtered", pkt.Data)
	}
	data := []byte(pkt.Data)
	var dataLen = uint32(len(data))
	if err = binary.Write(outfp, binary.LittleEndian, dataLen); err != nil {
		log.Print("error writing header3", err)
		return
	}
	var n int
	if n, err = outfp.Write(data); n != int(dataLen) || err != nil {
		log.Print("error writing body", err, n)
	}
	return
}

func convertToJSON(infp io.Reader, outfp io.Writer) error {
	var pretime uint64
	encodeTo := json.NewEncoder(outfp)
	for {
		if pkt, err := read1(infp); err != nil {
			if err != io.EOF {
				log.Print("error reading", err)
			}
			break
		} else {
			if pretime == 0 {
				pkt.TimeDelta = 0
			} else {
				pkt.TimeDelta = pkt.Time - pretime
			}
			pretime = pkt.Time
			encodeTo.Encode(pkt)
		}
	}
	return nil
}

func convertFromJSON(infp io.Reader, outfp io.Writer) error {
	decodeFrom := json.NewDecoder(infp)
	var cur uint64
	for {
		var pkt Packet
		if err := decodeFrom.Decode(&pkt); err != nil {
			if err != io.EOF {
				log.Print("read error", err)
			}
			break
		}
		if cur != 0 {
			pkt.Time = cur + pkt.TimeDelta
		}
		cur = pkt.Time
		pkt.Second = uint32(pkt.Time / 1000000)
		pkt.Usecond = uint32(pkt.Time % 1000000)
		pkt.Length = uint32(len([]byte(pkt.Data)))
		if err := write1(outfp, pkt); err != nil {
			log.Print("write error", err)
			break
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	name := path.Base(os.Args[0])
	app.Name = name
	app.Version = "0.0.1"
	app.Usage = "ttyrec <-> json converter"
	app.UsageText = fmt.Sprintf("%s [tojson|fromjson] < input > output", name)
	app.Commands = []cli.Command{
		{
			Name:  "tojson",
			Usage: "convert to json from ttyrec",
			Action: func(c *cli.Context) error {
				return convertToJSON(os.Stdin, os.Stdout)
			},
		}, {
			Name:  "fromjson",
			Usage: "convert from json to ttyrec",
			Action: func(c *cli.Context) error {
				return convertFromJSON(os.Stdin, os.Stdout)
			},
		},
	}
	app.Run(os.Args)
}
