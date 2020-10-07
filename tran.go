package trans

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const gURL = "https://script.google.com/macros/s/" +
	"AKfycbxCmg2CHtCFEzF5mYDwHX2iJS_BlVvEz52F-zcAMPu1jYhqS_E/exec"

type TransData struct {
	Code    int    `json:"code"`
	Text    string `json:"text"`
	Message string `json:"message"`
}

func Translate(text, source, target string) (string, error) {
	v := url.Values{}
	v.Add("text", text)
	v.Add("srouce", source)
	v.Add("target", target)

	resp, err := http.PostForm(gURL, v)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
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
