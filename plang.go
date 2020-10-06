package trans

import (
	"bytes"
	"strings"
	"text/template"
)

type plang struct {
	name string
	tmpl string
}

var plangmap = map[string]plang{
	"c": {"C (programming language)",
		`#include <stdio.h>
int main() {
    printf("{{.}}");
}`},
	"c+": {"C++ (programming language)",
		`#include <iostream>
using namespace std;
int main() {
    cout << "{{.}}" << endl;
}`},
	"j": {"Java (programming language)",
		`public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("{{.}}");
    }
}`},
	"go": {"Go (programming language)",
		`package main

import "fmt"

func main() {
	fmt.Println("{{.}}")
}`},
	"rb": {"Ruby (programming language)",
		`puts "{{.}}"`},
	"py": {"Python (programming language)",
		`print("{{.}}")`},
	"js": {"JavaScript (programming language)",
		`console.log("{{.}}")`},
	"tp": {"TypeScript (programming language)",
		`const s: string = "{{.}}"
console.log(s)`},
	"hs": {"Haskell (programming language)",
		`main = putStrLn "{{.}}"`},
	"rs": {"Rust (programming language)",
		`fn main() {
    println!("hello, world");
}`},
	"v": {"Vim script (scripting language)",
		`echo "{{.}}"`},
	"em": {"Emacs Lisp (scripting language)",
		`(princ "{{.}}")`},
}

func LookupPlang(lang string) (code, name string, ok bool) {
	k := strings.ToLower(lang)
	pl, found := plangmap[k]
	if !found {
		return "", "", false
	}
	return k, pl.name, true
}

func lookupTmpl(lang string) (tmpl string, ok bool) {
	k := strings.ToLower(lang)
	pl, found := plangmap[k]
	if !found {
		return "", false
	}
	return pl.tmpl, true
}

func Ptranslate(text, target string) (translated string, ok bool) {
	var buf bytes.Buffer
	tt, found := lookupTmpl(strings.ToLower(target))
	if !found {
		return "", false
	}
	tmpl := template.Must(template.New("plang").Parse(tt))
	tmpl.Execute(&buf, text)
	return string(buf.Bytes()), true
}
