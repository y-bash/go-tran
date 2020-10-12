package config

import (
	"strings"
	"testing"

	"github.com/morikuni/aec"
	"github.com/y-bash/go-tran"
)

type TomlToConfigTest struct {
	toml   Toml
	config Config
	err    string
}

var tomltoconfigtests = []TomlToConfigTest{
	0: {
		Toml{
			Default{"", "ja"}, API{"url", 3},
			Colors{"#000000", "#000000", "#000000", "#000000"},
		},
		Config{
			"", "Auto", "ja", "Japanese", tran.Endpoint("url"), 3,
			aec.FullColorF(0x0, 0x0, 0x0), aec.FullColorF(0x0, 0x0, 0x0),
			aec.FullColorF(0x0, 0x0, 0x0), aec.FullColorF(0x0, 0x0, 0x0),
		},
		"",
	},
	1: {
		Toml{
			Default{"ja", "en"}, API{"uri", 4},
			Colors{"#ffeedd", "#ccbbaa", "#998877", "#665544"},
		},
		Config{
			"ja", "Japanese", "en", "English", tran.Endpoint("uri"), 4,
			aec.FullColorF(0xff, 0xee, 0xdd), aec.FullColorF(0xcc, 0xbb, 0xaa),
			aec.FullColorF(0x99, 0x88, 0x77), aec.FullColorF(0x66, 0x55, 0x44),
		},
		"",
	},
	2: {Toml{Default{"zz", ""}, API{}, Colors{}},
		Config{}, "source is invalid"},
	3: {Toml{Default{"", "zz"}, API{}, Colors{}},
		Config{}, "target is invalid"},
	4: {Toml{Default{"", "ja"}, API{"", 1}, Colors{}},
		Config{}, "endpoint is invalid"},
	5: {Toml{Default{"", "ja"}, API{"url", 0}, Colors{}},
		Config{}, "limit_n_chars is invalid"},
	6: {Toml{Default{"", "ja"}, API{"url", 1}, Colors{"#Z", "", "", ""}},
		Config{}, "info is invalid"},
	7: {Toml{Default{"", "ja"}, API{"url", 1}, Colors{"#000000", "#Z", "", ""}},
		Config{}, "state is invalid"},
	8: {Toml{Default{"", "ja"}, API{"url", 1}, Colors{"#000000", "#000000", "#Z", ""}},
		Config{}, "error is invalid"},
	9: {Toml{Default{"", "ja"}, API{"url", 1}, Colors{"#000000", "#000000", "#000000", "#Z"}},
		Config{}, "result is invalid"},
}

func TestTomlToConfig(t *testing.T) {
	for i, tt := range tomltoconfigtests {
		config, err := tomlToConfig(&tt.toml)
		if err != nil {
			if tt.err == "" {
				t.Errorf("#%d\n\thave error: %s,\n\twant error: nil", i, err)
				continue
			}
			if !strings.Contains(err.Error(), tt.err) {
				t.Errorf("#%d\n\thave error: %s,\n\twant error: %q", i, err, tt.err)
			}
			continue
		}
		if tt.err != "" {
			panic("yyy")
		}
		if config.DefaultSourceCode != tt.config.DefaultSourceCode {
			t.Errorf("#%d have: config.DefaultSourceCode = %s, want: %s",
				i, config.DefaultSourceCode, tt.config.DefaultSourceCode)
		}
		if config.DefaultSourceName != tt.config.DefaultSourceName {
			t.Errorf("#%d have: config.DefaultSourceName = %s, want: %s",
				i, config.DefaultSourceName, tt.config.DefaultSourceName)
		}
		if config.DefaultTargetCode != tt.config.DefaultTargetCode {
			t.Errorf("#%d have: config.DefaultTargetCode = %s, want: %s",
				i, config.DefaultTargetCode, tt.config.DefaultTargetCode)
		}
		if config.DefaultTargetName != tt.config.DefaultTargetName {
			t.Errorf("#%d have: config.DefaultTargetName = %s, want: %s",
				i, config.DefaultTargetName, tt.config.DefaultTargetName)
		}
		if config.APIEndpoint != tt.config.APIEndpoint {
			t.Errorf("#%d have: config.APIEndpoint = %s, want: %s",
				i, config.APIEndpoint, tt.config.APIEndpoint)
		}
		if config.APILimitNChars != tt.config.APILimitNChars {
			t.Errorf("#%d have: config.APILimitNChars = %d, want: %d",
				i, config.APILimitNChars, tt.config.APILimitNChars)
		}
		if config.InfoColor.String() != tt.config.InfoColor.String() {
			t.Errorf("#%d have: config.InfoColor = %s, want: %s",
				i, config.InfoColor.String(), tt.config.InfoColor.String())
		}
		if config.StateColor.String() != tt.config.StateColor.String() {
			t.Errorf("#%d have: config.StateColor = %s, want: %s",
				i, config.StateColor.String(), tt.config.StateColor.String())
		}
		if config.ErrorColor.String() != tt.config.ErrorColor.String() {
			t.Errorf("#%d have: config.ErrorColor = %s, want: %s",
				i, config.ErrorColor.String(), tt.config.ErrorColor.String())
		}
		if config.ResultColor.String() != tt.config.ResultColor.String() {
			t.Errorf("#%d have: config.ResultColor = %s, want: %s",
				i, config.ResultColor.String(), tt.config.ResultColor.String())
		}
	}
}
