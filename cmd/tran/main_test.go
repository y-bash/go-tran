package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

type ScanTextTest struct {
	in    string
	limit int
	out   []string
}

var scantexttests = []ScanTextTest{
	0: {"", 1, []string{}},
	1: {"123\n456\n789", 1, []string{"123\n", "456\n", "789\n"}},
	2: {"123\n456\n789", 2, []string{"123\n", "456\n", "789\n"}},
	3: {"123\n456\n789", 3, []string{"123\n", "456\n", "789\n"}},
	4: {"123\n456\n789", 4, []string{"123\n", "456\n", "789\n"}},
	5: {"123\n456\n789", 5, []string{"123\n456\n", "789\n"}},
	6: {"123\n456\n789", 6, []string{"123\n456\n", "789\n"}},
	7: {"123\n456\n789", 7, []string{"123\n456\n", "789\n"}},
	8: {"123\n456\n789", 8, []string{"123\n456\n", "789\n"}},
	9: {"123\n456\n789", 9, []string{"123\n456\n789\n"}},
	10: {"123\n456\n789", 20, []string{"123\n456\n789\n"}},
}

func format(in []string) string {
	out := []string{}
	for i, s := range in {
		out = append(out, fmt.Sprintf("%d: %q", i, s))
	}
	return "[" + strings.Join(out, ", ") + "]"
}

func TestScanText(t *testing.T) {
	for i, tt := range scantexttests {
		r := strings.NewReader(tt.in)
		sc := bufio.NewScanner(r)
		out := []string{}
		for {
			s, eof := scanText(sc, tt.limit)
			if eof {
				break
			}
			out = append(out, s)
		}
		if len(out) != len(tt.out) {
			t.Errorf("#%d len(out) = \nhave: %d %s, \nwant: %d %s",
				i, len(out), format(out), len(tt.out), format(tt.out))
			continue
		}
		for j, s := range out {
			if s != tt.out[j] {
				t.Errorf("#%d [%d] have: %s, want: %s", i, j, s, tt.out[j])
			}
		}
	}
}
