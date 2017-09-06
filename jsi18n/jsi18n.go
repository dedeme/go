// Copyright 23-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

/*
Internationalization management.

jsi18n Use

From command line use:
  jsi18n <languages> <roots>
For example:
  jsi18n "en:es_ES" "src:/deme/lib/src"
Note that different languages and roots are separated by ':'

Operations with jsi18n

jsi18n do next operations:

  1. Create directory "i18n" in the current directory and inside it creates
     so many files .txt as indicated <languages>.
  2. Read all files ".js" in <roots> and all its subdirectories, extracting
     those strings which have keys to translate.
  3. Update .txt files.
  4. Create-Update file "src/i18n.js"

Process of Use

jsi18n requires next use cycle:

  1. Run jsi18n in the inmediate superior directory of 'src' directory
     (e.g.: jsi18n "en:es" "src:/deme/lib/src")
  2. Modify .txt files translating keys
  3. Rerun jsi18n to update "src/i18n.js"
  4. Recompile sources (which will include "i18n.js")

Annotation of keys

You have to mark keys to translate using '_("key")'.

For example:
  fmt.Println(_("A message"))
The expression _(" must be written as is, without inner spaces or tabulations.

Key Syntax

Keys are ordinary strings, but they must not include the character '='.

Translation Syntax

You have to put the translation after the first character '='.

Translations must be written like strings in code.

For example:
  That is a \"problem\".\nDon't forget it.

Translations marked as "ORPHAN" have nonexistent keys. If there are no error,
they should be deleted.

Keys marked as "TO DO" are pending translation.

Example of txt File

Example of file "es.txt":
  # File generated by jsi18n.

  # src/hello.js: 25
  Hello = Hola

  # src/hello.js: 48
  problem = Esto es un \"problema\".\nNo olvidarlo.

  # src/hello.go: 5
  start = comienzo
Lines commented with '#' and keys are automatically written by jsi18n. User
should write only after "=".

Use of File i18n

i18n.js can be used by program following next process:
  1. Set variable 'I18n.lang' with one of available dictionaries.
     For example: I18n.lang = I18n.es
  2. Use functions '_' and '_args' to write keys

If a key is not in dictionary, '_' and '_args' will yield such key without
modification; otherwise it will be translated.

Example of use:
  function main() {
    I18n.lang = I18n.es;
    ...
    console.log(_("Hello"));
    ...
    let day = "3";
    let hour = "14";
    console.log(_args(_("Day %0, at %1 p.m."), day, hour));
    ...
  }
  =>
  Hola
  Dia 3, a las 14 p.m.
'_args' allows arguments '0' to '9'. Each argument matches a variable. '0'
matches the first variable, '1' the second and so on.

Each argument can appear more than one time.
*/
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

// dmio functions ------------------------------------------

func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func IsDirectory(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}

func Wopen(path string) *os.File {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func Ropen(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func Write(file *os.File, text string) {
	_, err := file.WriteString(text)
	if err != nil {
		log.Fatal(err)
	}
}

func Lines(path string, f func(s string)) {
	file := Ropen(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// END dmio functions --------------------------------------

func help(msg string) {
	fmt.Println(msg + "\n" +
		"Use:\n" +
		"  jsi18n <languages> <roots>\n" +
		"For example:\n" +
		"  jsi18n \"en:es_ES\" \"src:/deme/lib/src\"\n" +
		"Note that different languages and roots are separated by ':'\n")
}

type pos struct {
	file string
	line int
}

type strs []string

func (ss strs) Len() int {
	return len(ss)
}

func (ss strs) Less(i, j int) bool {
	return strings.ToLower(ss[i]) < strings.ToLower(ss[j])
}

func (ss strs) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

type keys map[string][]*pos // key: position

type oldDic map[string]string // key: translation

func rdoc(text string) string {
	text = strings.Replace(text, "\\\"", "\"", -1)
	return strings.Replace(text, "\"", "\\\"", -1)
}

func extract(ks keys, p string) keys {
	fmt.Println(p)
	readFile := func() {
		const (
			CODE = iota
			PRECOMMENT
			LCOMMENT
			BCOMMENT
			BCOMMENT2
			QUOTE
			QUOTE2
			ENTRY1
			ENTRY2
			EQUOTE
			EQUOTE2
		)
		codeState := func(ch rune) int {
			switch ch {
			case '/':
				return PRECOMMENT
			case '"':
				return QUOTE
			case '_':
				return ENTRY1
			}
			return CODE
		}
		bf := new(bytes.Buffer)
		state := CODE
		nl := 1
		Lines(p, func(l string) {
			for _, ch := range l {
				switch state {
				case PRECOMMENT:
					switch ch {
					case '/':
						state = LCOMMENT
					case '*':
						state = BCOMMENT
					default:
						state = codeState(ch)
					}
				case LCOMMENT:
				case BCOMMENT:
					if ch == '*' {
						state = BCOMMENT2
					}
				case BCOMMENT2:
					if ch == '/' {
						state = CODE
					} else {
						state = BCOMMENT
					}
				case QUOTE:
					switch ch {
					case '"', '\n':
						state = CODE
					case '\\':
						state = QUOTE2
					}
				case QUOTE2:
					state = QUOTE
				case ENTRY1:
					if ch == '(' {
						state = ENTRY2
					} else {
						state = codeState(ch)
					}
				case ENTRY2:
					if ch == '"' {
						state = EQUOTE
					} else {
						state = codeState(ch)
					}
				case EQUOTE:
					switch ch {
					case '"', '\n':
						if ch == '"' {
							key := bf.String()
							if strings.Index(key, "=") == -1 {
								value := &pos{p, nl}
								data, ok := ks[key]
								if ok {
									ks[key] = append(data, value)
								} else {
									ks[key] = []*pos{value}
								}
							}
						}
						bf = new(bytes.Buffer)
						state = CODE
					case '\\':
						bf.WriteRune(ch)
						state = EQUOTE2
					default:
						bf.WriteRune(ch)
					}
				case EQUOTE2:
					bf.WriteRune(ch)
					state = EQUOTE
				default:
					state = codeState(ch)
				}
			}
			if state == LCOMMENT {
				state = CODE
			}
			nl++
		})
	}
	if IsDirectory(p) {
		files, _ := ioutil.ReadDir(p)
		for _, file := range files {
			ks = extract(ks, path.Join(p, file.Name()))
		}
	} else {
		// --------------------------------------------------------------------
		if path.Ext(p) == ".hx" {
			readFile()
		}
	}
	return ks
}

func makeDic(ks keys, lang string) string {
	dir := "i18n"
	if !Exists(dir) {
		os.Mkdir(dir, 0755)
	}

	filePath := path.Join(dir, lang+".txt")
	if !Exists(filePath) {
		Wopen(filePath).Close()
	}

	oldD := make(oldDic)
	Lines(filePath, func(l string) {
		l = strings.TrimSpace(l)
		if l == "" || l[0] == '#' {
			return
		}
		ix := strings.Index(l, "=")
		if ix == -1 {
			return
		}
		key := strings.TrimSpace(l[:ix])
		value := strings.TrimSpace(l[ix+1:])
		if key != "" && value != "" {
			oldD[rdoc(key)] = rdoc(value)
		}
	})

	orphan := ""
	todo := ""
	trans := ""
	dic := ""

	var kks strs
	for k := range ks {
		kks = append(kks, k)
	}
	sort.Sort(kks)
	for _, k := range kks {
		poss := ks[k]
		v, ok := oldD[k]
		if ok {
			for _, p := range poss {
				trans = trans + "# " + p.file + ": " + strconv.Itoa(p.line) + "\n"
			}
			trans = trans + k + " = " + v + "\n\n"
			dic = dic + "\t\"" + k + "\":  \"" + v + "\",\n"
			delete(oldD, k)
		} else {
			todo = todo + "# TO DO\n"
			for _, p := range poss {
				todo = todo + "# " + p.file + ": " + strconv.Itoa(p.line) + "\n"
			}
			todo = todo + k + " = \n\n"
		}
	}

	var oks strs
	for k := range oldD {
		oks = append(oks, k)
	}
	sort.Sort(oks)
	for _, k := range oks {
		orphan = orphan + "# ORPHAN\n" +
			k + " = " + oldD[k] + "\n\n"
	}

	fdic := Wopen(filePath)
	Write(fdic, "# File generated by jsi18n.\n\n")
	Write(fdic, orphan)
	Write(fdic, todo)
	Write(fdic, trans)
	fdic.Close()

	return dic[:len(dic)-2] + "\n"
}

func main() {
	if len(os.Args) != 3 {
		help("Missing package name and/or languages")
		return
	}
	langs := strings.Split(os.Args[1], ":")
	roots := strings.Split(os.Args[2], ":")
	for _, root := range roots {
		if !IsDirectory(root) {
			help("'" + root + "' is not a directory")
			return
		}
	}

	ks := make(keys)
	for _, root := range roots {
		ks = extract(ks, root)
	}

	jstarget := path.Join("src/i18n.js")

	fjstarget := Wopen(jstarget)
	defer fjstarget.Close()

	Write(fjstarget,
		`// Generate by jsi18n. Don't modify

goog.provide("I18n");

`)

	for _, lang := range langs {
		dic := makeDic(ks, lang)
		Write(fjstarget, "I18n."+lang+" = {\n"+
			dic+
			"};\n\n")
	}

	Write(fjstarget,
		`I18n.lang = {};

function _(key) {
  let v = I18n.lang[key];
  if (v !== undefined) {
    return v;
  }
  return key;
}

function _args(key, ...args) {
  let bf = "";
  v = _(key);
  isCode = false;
  for (let i = 0; i < v.length; ++i) {
    let ch = v.charAt(i);
    if (isCode) {
      bf += ch === "0" ? args[0]
        : ch === "1" ? args[1]
        : ch === "2" ? args[2]
        : ch === "3" ? args[3]
        : ch === "4" ? args[4]
        : ch === "5" ? args[5]
        : ch === "6" ? args[6]
        : ch === "7" ? args[7]
        : ch === "8" ? args[8]
        : ch === "9" ? args[9]
        : ch === "%" ? args[%]
        : "%" + ch;
      isCode = false;
    } else {
      if (ch === '%') {
        isCode = true;
      } else {
        bf += ch
      }
    }
  }
  return bf;
}
`)
}
