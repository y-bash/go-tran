package config

import (
	"fmt"

	"github.com/morikuni/aec"
	"github.com/y-bash/go-tran"
)

const (
	cInfo   = "#80a0d0" // Blue  ( 96, 192,  96)
	cState  = "#60c060" // Green ( 96, 192,  96) - State changed
	cError  = "#d04040" // Red   (208,  64,  64)
	cResult = "#ffc864" // Yellow(255, 200, 100) - Translation result
)

func hex2ansi(hex string) (aec.ANSI, error) {
	var r, g, b uint8
	if _, err := fmt.Sscanf(hex, "#%2x%2x%2x", &r, &g, &b); err != nil {
		return nil, fmt.Errorf("invalid: %s", hex)
	}
	return aec.FullColorF(r, g, b), nil
}

/*
func hex2ansiMust(hex, alt string) aec.ANSI {
	if ansi, err := hex2ansi(hex); err == nil {
		return ansi
	}
	ansi, err := hex2ansi(alt)
	if err != nil {
		panic(err) // xxx
	}
	return ansi
}
*/

type Config struct {
	DefaultSource  string
	DefaultTarget  string
	APIEndpoint    tran.Endpoint
	APIMaxNumLines uint
	InfoColor      aec.ANSI
	StateColor     aec.ANSI
	ErrorColor     aec.ANSI
	ResultColor    aec.ANSI
}

func Load() (*Config, error) {
	var initial Toml
	initial.Default.Source, _ = tran.CurrentLang()
	initial.Default.Target = ""
	initial.API.Endpoint = string(tran.DefaultAPI())
	initial.API.MaxNumLines = 30 // xxx
	initial.Colors.Info = cInfo
	initial.Colors.State = cState
	initial.Colors.Error = cError
	initial.Colors.Result = cResult

	var err error
	loaded, err := loadToml(&initial)
	if err != nil {
		return nil, err
	}

	var config Config
	config.DefaultSource = loaded.Default.Source
	config.DefaultTarget = loaded.Default.Target
	config.APIEndpoint = tran.Endpoint(loaded.API.Endpoint)
	config.APIMaxNumLines = loaded.API.MaxNumLines
	config.InfoColor, err = hex2ansi(loaded.Colors.Info)
	if err != nil {
		return nil, fmt.Errorf("config.toml;[colors];info is %s", err.Error())
	}
	config.StateColor, err = hex2ansi(loaded.Colors.State)
	if err != nil {
		return nil, fmt.Errorf("config.toml;[color];state is %s", err.Error())
	}
	config.ErrorColor, err = hex2ansi(loaded.Colors.Error)
	if err != nil {
		return nil, fmt.Errorf("config.toml;[color];error is %s", err.Error())
	}
	config.ResultColor, err = hex2ansi(loaded.Colors.Result)
	if err != nil {
		return nil, fmt.Errorf("config.toml;[color];result is %s", err.Error())
	}

	return &config, nil
}
