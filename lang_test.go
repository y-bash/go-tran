package tran

import (
	"os"
	"testing"
)

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
