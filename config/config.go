package config

import (
	"fmt"

	"github.com/morikuni/aec"
	"github.com/y-bash/go-tran"
)

const (
	cInfo   = "#80a0d0" // Blue
	cState  = "#60c060" // Green  - State changed
	cError  = "#d04040" // Red
	cResult = "#ffc864" // Yellow - Translation result
)

func hex2ansi(hex string) (aec.ANSI, error) {
	var r, g, b uint8
	if _, err := fmt.Sscanf(hex, "#%2x%2x%2x", &r, &g, &b); err != nil {
		return nil, fmt.Errorf("invalid: %s", hex)
	}
	return aec.FullColorF(r, g, b), nil
}

type Config struct {
	DefaultSourceCode string
	DefaultSourceName string
	DefaultTargetCode string
	DefaultTargetName string
	APIEndpoint       tran.Endpoint
	APIMaxNumLines    uint
	InfoColor         aec.ANSI
	StateColor        aec.ANSI
	ErrorColor        aec.ANSI
	ResultColor       aec.ANSI
}

func initialToml() *Toml {
	var initial Toml
	initial.Default.Source = ""
	initial.Default.Target, _ = tran.CurrentLang()
	initial.API.Endpoint = string(tran.DefaultAPI())
	initial.API.MaxNumLines = 0
	initial.Colors.Info = cInfo
	initial.Colors.State = cState
	initial.Colors.Error = cError
	initial.Colors.Result = cResult
	return &initial
}

func tomlToConfig(toml *Toml) (*Config, error) {
	var config Config

	if toml.Default.Source == "" {
		config.DefaultSourceCode = ""
		config.DefaultSourceName = "Auto"
	} else {
		code, name, ok := tran.LookupLangCode(toml.Default.Source)
		if !ok {
			return nil, fmt.Errorf(
				"config.toml;[default];source is invalid: %s", code)
		}
		config.DefaultSourceCode = code
		config.DefaultSourceName = name
	}

	code, name, ok := tran.LookupLangCode(toml.Default.Target)
	if !ok {
		return nil, fmt.Errorf(
			"config.toml;[default];target is invalid: %s", code)
	}
	config.DefaultTargetCode = code
	config.DefaultTargetName = name

	config.APIEndpoint = tran.Endpoint(toml.API.Endpoint)
	config.APIMaxNumLines = toml.API.MaxNumLines

	var err error
	config.InfoColor, err = hex2ansi(toml.Colors.Info)
	if err != nil {
		return nil, fmt.Errorf("config.toml;[colors];info is %s", err.Error())
	}
	config.StateColor, err = hex2ansi(toml.Colors.State)
	if err != nil {
		return nil, fmt.Errorf("config.toml;[color];state is %s", err.Error())
	}
	config.ErrorColor, err = hex2ansi(toml.Colors.Error)
	if err != nil {
		return nil, fmt.Errorf("config.toml;[color];error is %s", err.Error())
	}
	config.ResultColor, err = hex2ansi(toml.Colors.Result)
	if err != nil {
		return nil, fmt.Errorf("config.toml;[color];result is %s", err.Error())
	}

	return &config, nil
}

func Load() (*Config, error) {
	loaded, err := loadToml(initialToml())
	if err != nil {
		return nil, err
	}
	return tomlToConfig(loaded)
}
