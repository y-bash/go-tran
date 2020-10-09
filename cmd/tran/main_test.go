package main

import (
	"io/ioutil"
	"testing"
	"strings"
)

type ReadfilesTest struct{
	in  []string
	out string
}

var readfilestests = []ReadfilesTest{
	0: {[]string{"testdata/00in01.txt","testdata/00in02.txt","testdata/00in03.txt"},
		"testdata/00out.txt"},
	1: {[]string{"testdata/00in01.txt"}, "testdata/00in01.txt"},
}

func TestReadfiles(t *testing.T) {
	for i, tt := range readfilestests {
		out, err := readfiles(tt.in)
		if err != nil {
			t.Errorf("#%d testdata is failed: %s: ", i, err.Error())
			continue
		}
		have := strings.Join(out, "")
		buf, err := ioutil.ReadFile(tt.out)
		if err != nil {
			t.Errorf("#%d testdata is failed: %s: ", i, err.Error())
			continue
		}
		want := string(buf)
		if have != want {
			t.Errorf("#%d readfiles(%v) = \nhave:\n%v, \nwant:\n%v", i, tt.in, have, want)
		}
	}
}
