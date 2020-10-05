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
)

func interact(source, target string) error {
	fmt.Fprintln(os.Stderr, "Please enter something.")
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		in := sc.Text()
		out, err := trans.Translate(in, source, target)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, " => " + out)
	}
	return nil
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

func printLangCodes() {
	langs := trans.LangList()
	fmt.Println("ISO639-1 - Codes for the representation of names of languages.")
	fmt.Println("(https://en.wikipedia.org/wiki/ISO_639-1)")
	fmt.Println("---- -------------")
	fmt.Println("Code Language name")
	fmt.Println("---- -------------")
	for _,lang := range langs {
		fmt.Printf(" %s  %s\n", lang.Code, lang.Name)
	}
}

func main() {
	var help, lang bool
	var source, target string

	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&lang, "l", false, "Show ISO-639-1 Language codes")
	flag.StringVar(&source, "s", "", "Source language (ISO-639-1 code, Optional)")
	flag.StringVar(&target, "t", "ja", "Target language (ISO-639-1 code, Required)")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if lang {
		printLangCodes()
		return
	}

	if flag.NArg() == 0 && isTerminal(os.Stdin.Fd()) {
		err := interact(source, target)
		if err != nil {
			log.Fatal(err)
		}
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
