package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Default struct {
	Source string `toml:"source"`
	Target string `toml:"target"`
}

type API struct {
	Endpoint    string `toml:"endpoint"`
	MaxNumLines uint   `toml:"max_num_lines"`
}

type Colors struct {
	Info   string `toml:"info"`
	State  string `toml:"state"`
	Error  string `toml:"error"`
	Result string `toml:"result"`
}

type Toml struct {
	Default Default `toml:"default"`
	API     API     `toml:"api"`
	Colors  Colors  `toml:"colors"`
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (t *Toml) complete(initial *Toml) (overwritten bool) {
	if t.Default.Source == "" {
		t.Default.Source = initial.Default.Source
		overwritten = true
	}
	if t.Default.Target == "" {
		t.Default.Target = initial.Default.Target
		overwritten = true
	}
	if t.API.Endpoint == "" {
		t.API.Endpoint = initial.API.Endpoint
		overwritten = true
	}
	if t.API.MaxNumLines == 0 {
		t.API.MaxNumLines = initial.API.MaxNumLines
		overwritten = true
	}
	if t.Colors.Info == "" {
		t.Colors.Info = initial.Colors.Info
		overwritten = true
	}
	if t.Colors.State == "" {
		t.Colors.State = initial.Colors.State
		overwritten = true
	}
	if t.Colors.Error == "" {
		t.Colors.Error = initial.Colors.Error
		overwritten = true
	}
	if t.Colors.Result == "" {
		t.Colors.Result = initial.Colors.Result
		overwritten = true
	}
	return
}

func getTomlPath() (path string, err error) {
	var cfgdir string
	if runtime.GOOS == "windows" {
		appdir := os.Getenv("APPDATA")
		if appdir == "" {
			home := os.Getenv("USERPROFILE")
			appdir = filepath.Join(home, "Application Data")
		}
		cfgdir = filepath.Join(appdir, "y-bash", "tran")
	} else {
		home := os.Getenv("HOME")
		cfgdir = filepath.Join(home, ".config", "y-bash", "tran")
	}
	err = os.MkdirAll(cfgdir, 0700)
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgdir, "config.toml"), nil
}

func loadTomlFrom(path string, initial *Toml) (*Toml, error) {
	if exists(path) {
		var loaded Toml
		_, err := toml.DecodeFile(path, &loaded)
		if err != nil {
			return nil, err
		}
		if loaded.complete(initial) {
			f, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			toml.NewEncoder(f).Encode(loaded)
		}
		return &loaded, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	toml.NewEncoder(f).Encode(initial)
	return initial, nil
}

func loadToml(initial *Toml) (*Toml, error) {
	path, err := getTomlPath()
	if err != nil {
		return nil, err
	}
	return loadTomlFrom(path, initial)
}
