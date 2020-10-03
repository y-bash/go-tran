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
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		in := sc.Text()
		out, err := trans.Translate(in, source, target)
		if err != nil {
			return err
		}
		fmt.Println(" => " + out)
	}
	return nil
}

func read(f io.Reader) string {
	var sb strings.Builder
	sb.Grow(1024)
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
	var help bool
	var source, target string

	flag.BoolVar(&help, "h", false, "Show help")
	flag.StringVar(&source, "s", "", "Source language")
	flag.StringVar(&target, "t", "ja", "Target language")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if flag.NArg() == 0 && isTerminal(os.Stdin.Fd()) {
		fmt.Println("Please enter something")
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
