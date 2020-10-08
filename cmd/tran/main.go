package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/mattn/go-isatty"
	"github.com/morikuni/aec"
	"github.com/peterh/liner"
	"github.com/y-bash/go-tran"
)

var (
	cInfo   = aec.FullColorF(128, 160, 208) // Blue
	cState  = aec.FullColorF(96, 192, 96)   // Green - State changed
	cError  = aec.FullColorF(208, 64, 64)   // Red
	cResult = aec.FullColorF(255, 200, 100) // Yellow - Translation result
)

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

	fmt.Fprintln(os.Stderr, cInfo.Apply(text))
}

func langCodesToTerm(w io.Writer, substr string) (ok bool) {
	text := `┌──┬──────────┐
│Code│Language name       │
├──┼──────────┤
{{range .}}│ {{.Code}} │{{printf "%-20s" .Name}}│
{{end -}}
└──┴──────────┘ 
`
	a := tran.LangListContains(substr)
	if len(a) == 0 {
		return false
	}
	tmpl := template.Must(template.New("lang").Parse(text))
	var buf bytes.Buffer
	tmpl.Execute(&buf, a)
	fmt.Fprint(w, cInfo.Apply(string(buf.Bytes())))
	return true
}

func langCodesToNonTerm(w io.Writer, substr string) {
	text := `Code Language name
---- -------------
{{range .}} {{.Code}}  {{.Name}}
{{end -}}
`
	a := tran.LangListContains(substr)
	if len(a) == 0 {
		return
	}
	tmpl := template.Must(template.New("lang").Parse(text))
	tmpl.Execute(w, a)
}

func commandLangCodes(in string) {
	if in != "l" {
		in = in[2:]
	}
	if !langCodesToTerm(os.Stderr, in) {
		msg := cError.Apply("%q is not found\n")
		fmt.Fprintf(os.Stderr, msg, in)
	}
}

func commandSource(in, curr string) (source string, ok bool) {
	var arg string
	var code, name string
	switch {
	case len(in) == 1: // in is "s"
		if curr != "" {
			msg := cState.Apply("Source changed: Auto")
			fmt.Fprintln(os.Stderr, msg)
		}
		return "", true
	case len(in) >= 2: // in contains "s "
		arg = strings.TrimSpace(string([]rune(in)[2:]))
		code, name, ok = tran.LookupLang(arg)
	default:
		ok = false
	}
	if !ok {
		msg := cError.Apply("%s is not found\n")
		fmt.Fprintf(os.Stderr, msg, arg)
		return "", false
	}
	if curr != code {
		msg := cState.Apply("Source changed: %s (%s)\n")
		fmt.Fprintf(os.Stderr, msg, name, code)
	}
	return code, true
}

func commandTarget(in, curr string) (target string, ok bool) {
	var code, name string
	if in == "t" {
		code, name = tran.CurrentLang()
		msg := cState.Apply("Target changed: %s (%s)\n")
		fmt.Fprintf(os.Stderr, msg, name, code)
		return code, true
	}

	if strings.HasPrefix(in, "t ") {
		in = strings.TrimSpace(string([]rune(in)[2:]))
	}
	if code, name, ok = tran.LookupLang(in); !ok {
		code, name, ok = tran.LookupPlang(in)
	}
	if !ok {
		msg := cError.Apply("%q is not found\n")
		fmt.Fprintf(os.Stderr, msg, in)
		return "", ok
	}
	if curr != code {
		msg := cState.Apply("Target changed: %s (%s)\n")
		fmt.Fprintf(os.Stderr, msg, name, code)
	}
	return code, ok
}

func interact(source, target string) {
	fmt.Fprintln(os.Stderr, "Welcome to the GO-TRAN!")
	helpToTerm()
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
				fmt.Fprintln(os.Stderr, cResult.Apply(out))
			} else {
				out, err := tran.Translate(in, source, target)
				if err != nil {
					fmt.Fprintln(os.Stderr, cError.Apply(err.Error()))
				} else {
					fmt.Fprintln(os.Stderr, cResult.Apply(out))
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
	curr, _ := tran.CurrentLang()
	var help, lang bool
	var source, target string

	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&lang, "l", false, "Show language codes (ISO-639-1)")
	flag.StringVar(&source, "s", "", "Source language code (optional)")
	flag.StringVar(&target, "t", curr, "Target language code")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}
	if lang {
		langCodesToNonTerm(os.Stdout, "")
		return
	}
	if flag.NArg() == 0 && isTerminal(os.Stdin.Fd()) {
		interact(source, target)
		return
	}

	ss, err := readfiles(flag.Args())
	in := strings.Join(ss, "\n")

	out, err := tran.Translate(in, source, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
}
