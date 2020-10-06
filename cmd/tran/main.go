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
	"github.com/y-bash/go-trans"
)

func printHelp() {
	blue := aec.FullColorF(128, 160, 208)
	text := `--- --------------------------------  ---------
Cmd Description                       Examples
--- --------------------------------  ---------
 h  Show help                         :ja> h
 l  Show ISO-639-1 Language codes     :ja> l en
 s  Source language code (ISO-639-1)  :ja> s en
 t  Target language code (ISO-639-1)  :ja> t ja
 q  Quit                              :ja> q`
	fmt.Fprintln(os.Stderr, blue.Apply(text))
}

func printLangCodes(w io.Writer, substr string) {
	text := `---- -------------
Code Language name
---- -------------
{{range .}} {{.Code}}  {{.Name}}
{{end -}}
`
	a := trans.ContainsLangList(substr)
	if len(a) == 0 {
		return
	}
	tmpl := template.Must(template.New("lang").Parse(text))
	var buf bytes.Buffer
	tmpl.Execute(&buf, a)
	blue := aec.FullColorF(128, 160, 208)
	fmt.Fprint(w, blue.Apply(string(buf.Bytes())))
}

func commandSource(in, curr string) (source string, ok bool) {
	var arg string
	var code, name string
	switch {
	case len(in) == 1: // in is "s"
		if curr != "" {
			green := aec.FullColorF(96, 192, 96)
			msg := green.Apply("Source changed: Auto")
			fmt.Fprintln(os.Stderr, msg)
		}
		return "", true
	case len(in) >= 2: // in contains "s "
		arg = strings.TrimSpace(string([]rune(in)[2:]))
		code, name, ok = trans.LookupLang(arg)
	default:
		ok = false
	}
	if !ok {
		red := aec.FullColorF(208, 64, 64)
		msg := red.Apply("%s is not found\n")
		fmt.Fprintf(os.Stderr, msg, arg)
		return "", false
	}
	if curr != code {
		green := aec.FullColorF(96, 192, 96)
		msg := green.Apply("Source changed: %s (%s)\n")
		fmt.Fprintf(os.Stderr, msg, name, code)
	}
	return code, true
}

func commandTarget(in, curr string) (target string, ok bool) {
	var arg string
	var code, name string
	switch {
	case len(in) == 1: // in is "t"
		code, name = trans.CurrentLang()
		ok = true
	case len(in) >= 2: // in contains "t "
		arg = strings.TrimSpace(string([]rune(in)[2:]))
		code, name, ok = trans.LookupLang(arg)
		if !ok {
			code, name, ok = trans.LookupPlang(arg)
		}
	default:
		ok = false
	}
	if !ok {
		red := aec.FullColorF(208, 64, 64)
		msg := red.Apply("%q is not found\n")
		fmt.Fprintf(os.Stderr, msg, arg)
		return "", ok
	}
	if curr != code {
		green := aec.FullColorF(96, 192, 96)
		msg := green.Apply("Target changed: %s (%s)\n")
		fmt.Fprintf(os.Stderr, msg, name, code)
	}
	return code, ok
}

func interact(source, target string) {
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
		switch {
		case in == "q":
			fmt.Fprintln(os.Stderr, "Leaving GO-TRAN.")
			return

		case in == "h":
			printHelp()

		case in == "l" || strings.HasPrefix(in, "l "):
			var substr string
			if in != "l" {
				substr = in[2:]
			}
			printLangCodes(os.Stderr, substr)

		case in == "s" || strings.HasPrefix(in, "s "):
			if code, ok := commandSource(in, source); ok {
				source = code
			}
		case len(in) <= 2 || in == "t" || strings.HasPrefix(in, "t "):
			if in != "t" && len(in) <= 2 {
				in = "t " + in
			}
			if code, ok := commandTarget(in, target); ok {
				target = code
			}
		default:
			if out, ok := trans.Ptranslate(in, target); ok {
				fmt.Fprintln(os.Stderr, out)
			} else {
				out, err := trans.Translate(in, source, target)
				if err != nil {
					red := aec.FullColorF(208, 64, 64)
					fmt.Fprintln(os.Stderr, red.Apply(err.Error()))
				} else {
					yellow := aec.FullColorF(255, 200, 100)
					fmt.Fprintln(os.Stderr, yellow.Apply(out))
				}
			}
		}
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
	curr, _ := trans.CurrentLang()
	var help, lang bool
	var source, target string

	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&lang, "l", false, "Show ISO-639-1 Language codes")
	flag.StringVar(&source, "s", "", "Source language (ISO-639-1 code, Optional)")
	flag.StringVar(&target, "t", curr, "Target language (ISO-639-1 code, Required)")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}
	if lang {
		printLangCodes(os.Stdout, "")
		return
	}
	if flag.NArg() == 0 && isTerminal(os.Stdin.Fd()) {
		fmt.Fprintln(os.Stderr, "Welcome to the GO-TRAN!")
		printHelp()
		interact(source, target)
		return
	}

	ss, err := readfiles(flag.Args())
	in := strings.Join(ss, "\n")

	out, err := trans.Translate(in, source, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
}
