package main

import (
	"bytes"
	"testing"
)

func TestConvert(t *testing.T) {
	origStr := `{"Time":1490801541833151,"TimeDelta":0,"Data":"\u001b[?1034h"}
{"Time":1490801541833483,"TimeDelta":332,"Data":"bash-3.2$ "}
{"Time":1490801542900585,"TimeDelta":1067102,"Data":"exit\r\n"}
`
	infp := bytes.NewBufferString(origStr)
	var outfp bytes.Buffer
	convertFromJSON(infp, &outfp)
	infp2 := bytes.NewBufferString(outfp.String())
	var outfp2 bytes.Buffer
	convertToJSON(infp2, &outfp2)
	if outfp2.String() != origStr {
		t.Error("mismatch", outfp2.String(), origStr)
	}
}

func TestFilter(t *testing.T) {
	origStr := `{"Time":1490801541833483,"Data":"bash-3.2$ ","Filter":"tr a-z A-Z"}`
	expected := `{"Time":1490801541833483,"TimeDelta":0,"Data":"BASH-3.2$ "}
`
	infp := bytes.NewBufferString(origStr)
	var outfp bytes.Buffer
	convertFromJSON(infp, &outfp)
	infp2 := bytes.NewBufferString(outfp.String())
	var outfp2 bytes.Buffer
	convertToJSON(infp2, &outfp2)
	if outfp2.String() != expected {
		t.Error("mismatch", outfp2.String(), expected)
	}
}
