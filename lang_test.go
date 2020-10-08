package tran

import (
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
	8: {"zz", "", "", false},
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
