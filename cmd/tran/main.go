package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/mattn/go-isatty"
	"github.com/peterh/liner"
	"github.com/y-bash/go-tran"
	"github.com/y-bash/go-tran/config"
)

const version = "1.0.1"

var cfg *config.Config

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func apiScriptToNonTerm() {
	msg := `// You can register this script in your GAS (Google Apps
// Script) project, and set its URL in config.toml.
function doPost(e) {
    let body
    try {
        const p = e.parameter
        const s = LanguageApp.translate(p.text, p.source, p.target)
        body = {code: 200, text: s}
    } catch (e) {
        try {
            const msg = LanguageApp.translate(e.toString(), "", "en")
            body = {code: 400, message: msg}
        } catch (e) {
            body = {code: 500, message: e.toString()}
        }
    }

    let resp = ContentService.createTextOutput()
    resp.setMimeType(ContentService.MimeType.JSON)
    resp.setContent(JSON.stringify(body))

    return resp;
}`

	fmt.Fprintln(os.Stderr, msg)
}

func helpToNonTerm() {
	msg := `GO-TRAN (The language translator), version %s

Usage:  tran [option...] [file...]

Options:
    -a          show the script (Google Apps) for the API Server.
    -e          echo the source text.
    -h          show summary of options.
    -l          list the language codes(ISO639-1).
    -s CODE     specify the source language with CODE(ISO639-1).
    -t CODE     specify the target language with CODE(ISO639-1).
    -v          output version information.
`
	fmt.Fprintf(os.Stderr, msg, version)
}

func helpToTerm() {
	text := `┌──┬──────────┬────────┐
│Cmd │    Description     │    Examples    │
├──┼──────────┼──┬─────┤
│ h  │Show help           │h   │          │
│ l  │Show language codes │l en│l nor     │
│ s  │Source language code│s en│s french  │
│ t  │Target language code│t ja│t italian │
│ q  │Quit                │q   │          │
└──┴──────────┴──┴─────┘ `

	fmt.Fprintln(os.Stderr, cfg.InfoColor.Apply(text))
}

func langCodesToNonTerm(w io.Writer) {
	text := `Code Language name
---- -------------
{{range .}} {{.Code}}  {{.Name}}
{{end -}}
`
	a := tran.AllLangList()
	tmpl := template.Must(template.New("lang").Parse(text))
	tmpl.Execute(w, a)
}

func langCodesToTerm(w io.Writer, substr string) (ok bool) {
	text := `┌──┬──────────┐
│Code│Language name       │
├──┼──────────┤
{{range .}}│ {{.Code}} │{{printf "%-20s" .Name}}│
{{end -}}
└──┴──────────┘ 
`
	a := cfg.APIEndpoint.LangListContains(substr)
	if len(a) == 0 {
		return false
	}
	tmpl := template.Must(template.New("lang").Parse(text))
	var buf bytes.Buffer
	tmpl.Execute(&buf, a)
	fmt.Fprint(w, cfg.InfoColor.Apply(string(buf.Bytes())))
	return true
}

func brackets(s string) string {
	if s == "" {
		return ""
	}
	return "(" + s + ")"
}

func commandLangCodes(in string) {
	if in != "l" {
		in = in[2:]
	}
	if !langCodesToTerm(os.Stderr, in) {
		msg := cfg.ErrorColor.Apply("%q is not found\n")
		fmt.Fprintf(os.Stderr, msg, in)
	}
}

func commandSource(in, curr string) (source string, ok bool) {
	var code, name string
	if in == "s" {
		code = cfg.DefaultSourceCode
		name = cfg.DefaultSourceName
		ok = true
	} else {
		if strings.HasPrefix(in, "s ") {
			in = strings.TrimSpace(string([]rune(in)[2:]))
		}
		if code, name, ok = cfg.APIEndpoint.LookupLang(in); !ok {
			code, name, ok = tran.LookupPlang(in)
		}
		if !ok {
			msg := cfg.ErrorColor.Apply("%q is not found\n")
			fmt.Fprintf(os.Stderr, msg, in)
			return "", ok
		}
	}
	if curr != code {
		msg := cfg.StateColor.Apply("Srouce changed: %s %s\n")
		fmt.Fprintf(os.Stderr, msg, name, brackets(code))
	}
	return code, ok
}

func commandTarget(in, curr string) (target string, ok bool) {
	var code, name string
	if in == "t" {
		code = cfg.DefaultTargetCode
		name = cfg.DefaultTargetName
		ok = true
	} else {
		if strings.HasPrefix(in, "t ") {
			in = strings.TrimSpace(string([]rune(in)[2:]))
		}
		if code, name, ok = cfg.APIEndpoint.LookupLang(in); !ok {
			code, name, ok = tran.LookupPlang(in)
		}
		if !ok {
			msg := cfg.ErrorColor.Apply("%q is not found\n")
			fmt.Fprintf(os.Stderr, msg, in)
			return "", ok
		}
	}
	if curr != code {
		msg := cfg.StateColor.Apply("Target changed: %s %s\n")
		fmt.Fprintf(os.Stderr, msg, name, brackets(code))
	}
	return code, ok
}

func initialLang(source, target string) (string, string) {
	if source != "" {
		if code, ok := commandSource("s " + source, cfg.DefaultSourceCode); ok {
			source = code
		} else {
			source = cfg.DefaultSourceCode
		}
	}
	if target != "" {
		if code, ok := commandTarget(target, cfg.DefaultTargetCode); ok {
			target = code
		} else {
			target = cfg.DefaultTargetCode
		}
	}
	return source, target
}

func interact(source , target string) {
	fmt.Fprintf(os.Stderr, "Welcome to the GO-TRAN! (Ver %s)\n", version)
	helpToTerm()
	source, target = initialLang(source, target)

	line := liner.NewLiner()
	defer line.Close()
	for {
		pr := fmt.Sprintf("%s:%s> ", source, target)
		in, err := line.Prompt(pr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if len(in) <= 0 {
			continue
		}
		in = strings.TrimSpace(in)
		switch {
		case in == "q":
			fmt.Fprintln(os.Stderr, "Leaving GO-TRAN.")
			return

		case in == "h":
			helpToTerm()

		case in == "l" || strings.HasPrefix(in, "l "):
			commandLangCodes(in)

		case in == "s" || strings.HasPrefix(in, "s "):
			if code, ok := commandSource(in, source); ok {
				source = code
			}
		case len(in) <= 2 || strings.HasPrefix(in, "t "):
			if code, ok := commandTarget(in, target); ok {
				target = code
			}
		default:
			if out, ok := tran.Ptranslate(in, target); ok {
				fmt.Fprintln(os.Stderr, cfg.ResultColor.Apply(out))
			} else {
				out, err := cfg.APIEndpoint.Translate(in, source, target)
				if err != nil {
					fmt.Fprintln(os.Stderr, cfg.ErrorColor.Apply(err.Error()))
				} else {
					fmt.Fprintln(os.Stderr, cfg.ResultColor.Apply(out))
				}
			}
		}
		line.AppendHistory(in)
	}
}

func scanText(sc *bufio.Scanner, limit int) (out string, eof bool) {
	var sb strings.Builder
	sb.Grow(4096)
	for i := 0; i < limit; {
		if !sc.Scan() {
			break
		}
		s := sc.Text()
		sb.WriteString(s)
		sb.WriteString("\n")
		i += len([]rune(s)) + 1
	}
	out = sb.String()
	return out, len(out) == 0
}

func translate(w io.Writer, r io.Reader, srcEcho bool) error {
	source := cfg.DefaultSourceCode
	target := cfg.DefaultTargetCode
	tran := cfg.APIEndpoint.Translate
	limit := cfg.APILimitNChars
	sc := bufio.NewScanner(r)
	for {
		in, eof := scanText(sc, limit)
		if eof {
			break
		}
		out, err := tran(in, source, target)
		if err != nil {
			return err
		}
		if !srcEcho {
			fmt.Fprint(w, out)
			continue
		}
		iss := strings.Split(in, "\n")
		oss := strings.Split(out, "\n")
		for i, is := range iss {
			if i >= len(iss) - 1 && len(is) == 0 {
				continue
			}
			fmt.Fprintln(w, is)
			fmt.Fprintln(w, oss[i])
		}
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func batch(paths []string, srcEcho bool) {
	if len(paths) == 0 {
		translate(os.Stdout, os.Stdin, srcEcho)
		return
	}
	for _, path := range paths {
		if !exists(path) {
			fmt.Fprintf(os.Stderr, "GO-TRAN: %s:  No such file or directory\n", path)
			continue
		}
		if isDir(path) {
			fmt.Fprintf(os.Stderr, "GO-TRAN: %s: Is a directory\n", path)
			continue
		}
		var f *os.File
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "GO-TRAN: %s\n", err)
			continue
		}
		defer f.Close()
		translate(os.Stdout, f, srcEcho)
	}
	return
}

func main() {
	var api, srcEcho, help, lang, ver bool
	var source, target string

	flag.Usage	= helpToNonTerm
	flag.BoolVar(&api, "a", false, "show api (Google Apps Script)")
	flag.BoolVar(&srcEcho, "e", false, "echo the source text")
	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&lang, "l", false, "list the language codes (ISO-639-1)")
	flag.StringVar(&source, "s", "", "source language code")
	flag.StringVar(&target, "t", "", "target language code")
	flag.BoolVar(&ver, "v", false, "show version")
	flag.Parse()

	if api {
		apiScriptToNonTerm()
		return
	}
	if help {
		flag.Usage()
		return
	}
	if lang {
		langCodesToNonTerm(os.Stdout)
		return
	}
	if ver {
		fmt.Fprintf(os.Stderr, "GO-TRAN Version %s\n", version)
		return
	}

	var err error
	if cfg, err = config.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "GO-TRAN: %s\n", err)
		os.Exit(1)
	}
	if flag.NArg() == 0 && isTerminal(os.Stdin.Fd()) {
		interact(source, target)
		return
	}
	if err := cfg.ChangeDefault(source, target); err != nil {
		fmt.Fprintf(os.Stderr, "GO-TRAN: %s\n", err)
		return
	}
	batch(flag.Args(), srcEcho)
}
