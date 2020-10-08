package tran

import (
	"os"
	"testing"
)

type LookupLangTest struct {
	in   string
	code string
	name string
	ok   bool
}

var findlangtests = []LookupLangTest{
	0: {"ja", "ja", "Japanese", true},
	1: {"JA", "ja", "Japanese", true},
	2: {"jap", "ja", "Japanese", true},
	3: {"japan", "ja", "Japanese", true},
	4: {"jAPANESE", "ja", "Japanese", true},
	5: {"en", "en", "English", true},
	6: {"en", "en", "English", true},
	7: {"frisia", "fy", "Western Frisian", true},
	8: {"英語", "en", "English", true},
	9: {"zz", "", "", false},
}

func TestLookupLang(t *testing.T) {
	for i, tt := range findlangtests {
		code, name, ok := LookupLang(tt.in)
		if !ok {
			if tt.ok {
				t.Errorf("#%d have ok: %v, want ok: %v", i, ok, tt.ok)
			}
			continue
		}
		if !tt.ok {
			t.Errorf("#%d have ok: %v, want ok: %v", i, ok, tt.ok)
			continue
		}
		if code != tt.code || name != tt.name {
			t.Errorf("#%d LookupLang(%q) = (%q, %q, nil)  want: (%q, %q, nil)",
				i, tt.in, code, name, tt.code, tt.name)
		}
	}
}

type LangListContainsTest struct {
	in string
	a  string
	n  int
}

var langlistcontainstests = []LangListContainsTest{
	0: {"pan", "[ja:Japanese es:Spanish]", 2},
	1: {"", "", 184},
	2: {"xyz", "[]", 0},
}

func TestLangListContains(t *testing.T) {
	for i, tt := range langlistcontainstests {
		a := LangListContains(tt.in)
		if tt.in == "" {
			if len(a) != tt.n {
				t.Errorf("#%d len(LangListContains(%q)) = %d, want: %d",
					i, tt.in, len(a), tt.n)
			}
			continue
		}
		if a.String() != tt.a {
			t.Errorf("#%d LangListContains(%q) = \nhave:\t%v, \nwant:\t%v",
				i, tt.in, a, tt.a)
		}
	}
}

type CurrentLangTest struct {
	lang string
	code string
	name string
}

var currentlangtests = []CurrentLangTest{
	0: {"C.UTF8", "en", "English"},
	1: {"en_US.UTF-8", "en", "English"},
	2: {"ja_JP.UTF8", "ja", "Japanese"},
}

func TestCurrentLang(t *testing.T) {
	lang := os.Getenv("LANG")
	defer func() {
		os.Setenv("LANG", lang)
	}()
	for i, tt := range currentlangtests {
		os.Setenv("LANG", tt.lang)
		code, name := CurrentLang()
		if code != tt.code || name != tt.name {
			t.Errorf("#%d LANG=%q CurrentLang() = (%v, %v), want: (%v, %v)",
				i, tt.lang, code, name, tt.code, tt.name)
		}
	}
}
