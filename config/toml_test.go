package config

import (
	"io"
	"os"
	"testing"
)

type ExistsTest struct {
	path   string
	exists bool
}

var existstests = []ExistsTest{
	0: {"testdata/exists.txt", true},
	1: {"testdata/notexists.txt", false},
}

func TestExists(t *testing.T) {
	for i, tt := range existstests {
		have := exists(tt.path)
		if have != tt.exists {
			t.Errorf("#%d exists(%q) = %v, want: %v",
				i, tt.path, have, tt.exists)
		}
	}
}

func copyFile(dst, src string) error {
	srcf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcf.Close()

	dstf, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstf.Close()

	_, err = io.Copy(dstf, srcf)
	if err != nil {
		return err
	}
	return nil
}

func TestLoadTomlFrom(t *testing.T) {
	var initial1 Toml
	initial1.Default.Source = "1"
	initial1.Default.Target = "2"
	initial1.API.Endpoint = "3"
	initial1.API.MaxNumLines = 4
	initial1.Colors.Info = "#555555"
	initial1.Colors.State = "#666666"
	initial1.Colors.Error = "#777777"
	initial1.Colors.Result = "#888888"
	var initial2 Toml
	initial2.Default.Source = "5"
	initial2.Default.Target = "6"
	initial2.API.Endpoint = "7"
	initial2.API.MaxNumLines = 8
	initial2.Colors.Info = "#999999"
	initial2.Colors.State = "#AAAAAA"
	initial2.Colors.Error = "#BBBBBB"
	initial2.Colors.Result = "#CCCCCC"

	err := os.MkdirAll("../output", 0700)
	if err != nil {
		t.Errorf("testdata is failed: %s", err.Error())
		return
	}

	// Not Exists Test
	notExistsFile := "../output/load_notexists.toml"
	loaded, err := loadTomlFrom(notExistsFile, &initial1)
	if err != nil {
		t.Errorf("testdata is failed: %s", err.Error())
		return
	}
	if *loaded != initial1 {
		t.Errorf("loadTomlFrom(notexists, initial1) != initial1")
	}

	// Empty Toml Test
	emptyFile := "../output/load_empty.toml"
	err = copyFile(emptyFile, "testdata/load_empty.toml")
	if err != nil {
		t.Errorf("testdata is failed: %s", err.Error())
		return
	}
	loaded, err = loadTomlFrom(emptyFile, &initial1)
	if err != nil {
		t.Errorf("testdata is failed: %s", err.Error())
		return
	}
	if *loaded != initial1 {
		t.Errorf("loadTomlFrom(empty, initial1) != initial1")
	}

	// Filled Toml Test
	filledFile := emptyFile
	loaded, err = loadTomlFrom(filledFile, &initial2)
	if err != nil {
		t.Errorf("testdata is failed: %s", err.Error())
		return
	}
	if *loaded != initial1 {
		t.Errorf("loadTomlFrom(filled, initial2) != initial1")
	}
}
