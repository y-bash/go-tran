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

func commandLangCodes(in string) {
	if in != "l" {
		in = in[2:]
	}
	if !langCodesToTerm(os.Stderr, in) {
		msg := cfg.ErrorColor.Apply("%q is not found\n")
		fmt.Fprintf(os.Stderr, msg, in)
	}
}

func brackets(s string) string {
	if s == "" {
		return ""
	}
	return "(" + s + ")"
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

func interact() {
	fmt.Fprintf(os.Stderr, "Welcome to the GO-TRAN! (Ver %s)\n", version)
	helpToTerm()
	source := cfg.DefaultSourceCode
	target := cfg.DefaultTargetCode
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

func read(f io.Reader) string {
	var sb strings.Builder
	sb.Grow(4096)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		sb.WriteString(sc.Text())
		sb.WriteString("\n")
	}
	return sb.String()
}

func readfiles(paths []string) (out []string, err error) {
	if len(paths) == 0 {
		out = []string{read(os.Stdin)}
		return
	}
	for _, path := range paths {
		var f *os.File
		f, err = os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()
		out = append(out, read(f))
	}
	return
}

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func main() {
	var help, lang, v bool
	var source, target string

	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&lang, "l", false, "Show language codes (ISO-639-1)")
	flag.StringVar(&source, "s", "", "Source language code")
	flag.StringVar(&target, "t", "", "Target language code")
	flag.BoolVar(&v, "v", false, "Show version")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}
	if lang {
		langCodesToNonTerm(os.Stdout)
		return
	}
	if v {
		fmt.Fprintf(os.Stderr, "GO-TRAN Version %s\n", version)
		return
	}
	var err error
	if cfg, err = config.Load(source, target); err != nil {
		fmt.Fprintf(os.Stderr, "GO-TRAN: %s\n", err)
		os.Exit(1)
	}
	if flag.NArg() == 0 && isTerminal(os.Stdin.Fd()) {
		interact()
		return
	}
	ss, err := readfiles(flag.Args())
	in := strings.Join(ss, "")
	out, err := cfg.APIEndpoint.Translate(in, source, target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GO-TRAN: %s\n", err)
		os.Exit(1)
	}
	fmt.Print(out)
}
