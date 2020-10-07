package tran

import (
	"strings"
	"testing"
)

type LookupLangTest struct {
	in   string
	code string
	name string
	err  string
}

var findlangtests = []LookupLangTest{
	0: {"ja", "ja", "Japanese", ""},
	1: {"JA", "ja", "Japanese", ""},
	2: {"jap", "ja", "Japanese", ""},
	3: {"japan", "ja", "Japanese", ""},
	4: {"jAPANESE", "ja", "Japanese", ""},
	5: {"en", "en", "English", ""},
	6: {"en", "en", "English", ""},
	7: {"frisia", "fy", "Western Frisian", ""},
	8: {"zz", "", "", "not found"},
}

func TestLookupLang(t *testing.T) {
	for i, tt := range findlangtests {
		code, name, err := LookupLang(tt.in)
		if err != nil {
			if tt.err == "" {
				t.Errorf("#%d have error: %s, want error: none", i, err.Error())
				continue
			}
			if !strings.Contains(err.Error(), tt.err) {
				t.Errorf("#%d have error: %s, want error: %s", i, err.Error(), tt.err)
			}
			continue
		}
		if tt.err != "" {
			t.Errorf("#%d have error: none, want error: %s", i, tt.err)
			continue
		}
		if code != tt.code || name != tt.name {
			t.Errorf("#%d LookupLang(%q) = (%q, %q, nil)  want: (%q, %q, nil)",
				i, tt.in, code, name, tt.code, tt.name)
		}
	}
}
