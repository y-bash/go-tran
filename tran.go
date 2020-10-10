package tran

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Endpoint string

func DefaultAPI() Endpoint {
	return "https://script.google.com/macros/s/" +
		"AKfycbxCmg2CHtCFEzF5mYDwHX2iJS_BlVvEz52F-zcAMPu1jYhqS_E/exec"
}

func NewAPI(url string) Endpoint {
	return Endpoint(url)
}

type TransData struct {
	Code    int    `json:"code"`
	Text    string `json:"text"`
	Message string `json:"message"`
}

func (ep Endpoint) Translate(text, source, target string) (string, error) {
	v := url.Values{}
	v.Add("text", text)
	v.Add("srouce", source)
	v.Add("target", target)

	resp, err := http.PostForm(string(ep), v)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var td TransData
	if err := json.Unmarshal(buf, &td); err != nil {
		return "", err
	}
	if td.Code != 200 {
		msg := td.Message
		prefix := "exception:"
		if strings.HasPrefix(strings.ToLower(msg), prefix) {
			msg = string(msg[len(prefix):])
			msg = strings.TrimSpace(msg)
		}
		return "", errors.New(msg)
	}
	return td.Text, nil
}

func (ep Endpoint) LookupLang(s string) (code, name string, ok bool) {
	switch {
	case len(s) == 2:
		if code, name, ok = lookupLangCode(s); ok {
			return
		}
	case len(s) >= 3:
		if code, name, ok = lookupLangName(s); ok {
			return
		}
		if en, err := ep.Translate(s, "", "en"); err == nil {
			if code, name, ok = lookupLangName(en); ok {
				return
			}
		}
	default:
		// Do nothing
	}
	return "", "", false
}

func (ep Endpoint) LangListContains(substr string) ISO639List {
	if a := langListContains(substr); len(a) > 0 {
		return a
	}
	if en, err := ep.Translate(substr, "", "en"); err == nil {
		return langListContains(en)
	}
	return []*ISO639{}
}

func (ep Endpoint) CurrentLang() (code, name string) {
	var lang string
	if s, ok := os.LookupEnv("LANG"); ok {
		lang = strings.ToLower(s)
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "Get-Culture | Select-Object -exp Name")
		if bs, err := cmd.Output(); err == nil {
			lang = strings.ToLower(string(bs))
		}
	}
	if len(lang) >= 2 {
		code, name, ok := ep.LookupLang(string(lang[:2]))
		if ok {
			return code, name
		}
	}
	return "en", "English"
}
