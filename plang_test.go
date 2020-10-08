package tran

import (
	"testing"
)

type LookupPlangTest struct {
	in   string
	code string
	name string
	ok   bool
}

var lookupplangtests = []LookupPlangTest{
	0: {"c", "c", "C (programming language)", true},
	1: {"d", "", "", false},
}

func TestLookupPlang(t *testing.T) {
	for i, tt := range lookupplangtests {
		code, name, ok := LookupPlang(tt.in)
		if !tt.ok {
			if ok != tt.ok {
				t.Errorf("#%d LookupPlang(%q) = %v, want: %v",
					i, tt.in, ok, tt.ok)
			}
			continue
		}
		if code != tt.code || name != tt.name || ok != tt.ok {
			t.Errorf("#%d LookupPlang(%q) = (%q, %q), want: (%q, %q)",
				i, tt.in, code, name, tt.code, tt.name)
		}

	}
}

type PtranslateTest struct {
	text       string
	target     string
	translated string
	ok         bool
}

var ptranslatetest = []PtranslateTest{
	0: {"hello, world", "c", `#include <stdio.h>
int main() {
    printf("hello, world");
}`, true},
	1: {"hello, world", "d", "", false},
}

func TestPtranslate(t *testing.T) {
	for i, tt := range ptranslatetest {
		translated, ok := Ptranslate(tt.text, tt.target)
		if !tt.ok {
			if ok != tt.ok {
				t.Errorf("#%d Ptranslate(%q, %q) = %v, want: %v",
					i, tt.text, tt.target, ok, tt.ok)
			}
			continue
		}
		if translated != tt.translated {
			if translated != tt.translated || ok != tt.ok {
				t.Errorf("#%d Ptranslated(%q, %q) = (%q, %v), want: (%q, %v)",
					i, tt.text, tt.target, translated, ok, tt.translated, tt.ok)
			}
		}
	}
}
