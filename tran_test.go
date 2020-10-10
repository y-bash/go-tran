package tran

import (
	"strings"
	"testing"
)

type Endpoint_TranslateTest struct {
	in     string
	source string
	target string
	out    string
	err    string
}

var endpoint_translatetests = []Endpoint_TranslateTest{
	0: {"猫", "", "de", "Katze", ""},     // Deutsch
	1: {"猫", "", "en", "Cat", ""},       // English
	2: {"猫", "", "es", "Gato", ""},      // Spanish
	3: {"猫", "", "fr", "Chat", ""},      // French
	4: {"猫", "", "it", "Gatto", ""},     // Italian
	5: {"Cat", "", "ja", "ネコ", ""},      // Japanese
	6: {"Cat", "", "ko", "고양이", ""},     // Korean
	7: {"猫", "", "pt", "Gato", ""},      // Portuguese
	8: {"Cat", "", "zh", "猫", ""},       // Chinese
	9: {"Cat", "", "xx", "", "Invalid"}, // Invalid
}

func TestEndpoint_translate(t *testing.T) {
	ep := DefaultAPI()
	for i, tt := range endpoint_translatetests {
		out, err := ep.Translate(tt.in, tt.source, tt.target)
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
		if strings.ToLower(out) != strings.ToLower(tt.out) {
			t.Errorf("#%d Translate(%q, %q, %q) = %q, want: %q",
				i, tt.in, tt.source, tt.target, out, tt.out)
		}
	}
}

type Endpoint_LookupLangTest struct {
	in   string
	code string
	name string
	ok   bool
}

var endpoint_findlangtests = []Endpoint_LookupLangTest{
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

func TestEndpoint_LookupLang(t *testing.T) {
	ep := DefaultAPI()
	for i, tt := range endpoint_findlangtests {
		code, name, ok := ep.LookupLang(tt.in)
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

type Endpoint_LangListContainsTest struct {
	in string
	a  string
	n  int
}

var endpoint_langlistcontainstests = []Endpoint_LangListContainsTest{
	0: {"pan", "[ja:Japanese es:Spanish]", 2},
	1: {"", "", 184},
	2: {"xyz", "[]", 0},
}

func TestEndpoint_LangListContains(t *testing.T) {
	ep := DefaultAPI()
	for i, tt := range endpoint_langlistcontainstests {
		a := ep.LangListContains(tt.in)
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
