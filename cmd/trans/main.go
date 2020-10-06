package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/y-bash/go-trans"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

func printHelp() {
	text := `--- --------------------------------  --------
Cmd Description                       Examples
--- --------------------------------  --------
 h  Show help                         h
 l  Show ISO-639-1 Language codes     l en
 s  Source language code (ISO-639-1)  s en
 t  Target language code (ISO-639-1)  t ja
 q  Quit                              q`
	fmt.Fprintln(os.Stderr, text)
}

func printLangCodes(w io.Writer, substr string) {
	text := `ISO639-1 Codes for the representation of names of languages.
---- -------------
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
	tmpl.Execute(w, a)
}

func commandSource(in, curr string) (source string, ok bool) {
	var arg string
	var code, name string
	switch {
	case len(in) == 1: // in is "s"
		if curr != "" {
			fmt.Fprintln(os.Stderr, "Source changed: Auto")
		}
		return "", true
	case len(in) >= 2: // in contains "s "
		arg = strings.TrimSpace(string([]rune(in)[2:]))
		code, name, ok = trans.LookupLang(arg)
	default:
		ok = false
	}
	if !ok {
		fmt.Fprintf(os.Stderr, "%s is not found\n", arg)
		return "", false
	}
	if curr != code {
		fmt.Fprintf(os.Stderr, "Source changed: %s (%s)\n", name, code)
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
		fmt.Fprintf(os.Stderr, "%q is not found\n", arg)
		return "", ok
	}
	if curr != code {
		fmt.Fprintf(os.Stderr, "Target changed: %s (%s)\n", name, code)
	}
	return code, ok
}

func interact(source, target string) {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "%s:%s> ", source, target)
		if !sc.Scan() {
			break
		}
		in := strings.TrimSpace(sc.Text())
		if len(in) <= 0 {
			continue
		}
		switch {
		case in == "q":
			fmt.Fprintln(os.Stderr, "Leaving TRANS.")
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
					fmt.Fprintln(os.Stderr, err)
				} else {
					fmt.Fprintln(os.Stderr, out)
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
		fmt.Fprintln(os.Stderr, "Welcome to the GO-TRANS!")
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
