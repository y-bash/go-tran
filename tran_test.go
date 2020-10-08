package tran

import (
	"strings"
	"testing"
)

type TranslateTest struct {
	in     string
	source string
	target string
	out    string
	err    string
}

var translatetests = []TranslateTest{
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

func TestTranslate(t *testing.T) {
	for i, tt := range translatetests {
		out, err := Translate(tt.in, tt.source, tt.target)
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
