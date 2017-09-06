// Copyright 05-Sep-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

/*
Templates to easy making js code.

jstemplate Use

From command line use:
  jstemplate file
For example:
  jstemplate "src/main.js"

Working with jstemplate

Modify a js file changing template to javascript code.

If a fail happens the js file will not be modified.

Templates

Valid templates is:

  <<<class
  field1:type
  field2:type
  @field3:type
  ...
  >>>

It will be replaced by js class code (note that the field which starts with '@'
is codified as a read/write one.)
*/
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
  "strconv"
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
		"  jstemplate file\n" +
		"For example:\n" +
		"  jstemplate \"src/main.js\"\n")
}

func tClass(name string, lines []string) string {
	type field struct {
		id      string
		writing bool
		tp      string
	}
	fields := make([]field, 0)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			writing := false
			if l[0] == '@' {
				writing = true
				l = l[1:]
			}
			ix := strings.Index(l, ":")
			if ix == -1 {
				panic("':' is missing in '" + l + "'")
			}
			id := strings.TrimSpace(l[:ix])
			tp := strings.TrimSpace(l[ix+1:])
			if id == "" {
				panic("Field name is missing in '" + l + "'")
			}
			if tp == "" {
				panic("Field type is missing in '" + l + "'")
			}
			fields = append(fields, field{id, writing, tp})
		}
	}

	if len(fields) == 0 {
		panic("There are no fields in 'class'")
	}

	r := name + " = class {\n  /**\n"
	for _, f := range fields {
		r += "   * @param {" + f.tp + "} " + f.id + "\n"
	}
	r += "   */\n"

	r += "  constructor (\n" +
    "    " + fields[0].id
	for i := 1; i < len(fields); i++ {
		r += ", \n    " + fields[i].id
	}
	r += "\n  ) {\n"

	for _, f := range fields {
		r += "    /** @private */\n    this._" + f.id + " = " + f.id + ";\n"
	}
	r += "  }\n"

	for _, f := range fields {
		r += "\n  /** @return {" + f.tp + "} */\n" +
			"  get " + f.id + " () {\n" +
			"    return this._" + f.id + ";\n" +
			"  }\n"

		if f.writing {
			r += "\n  /** @param {" + f.tp + "} value */\n" +
				"  set " + f.id + " (value) {\n" +
				"    this._" + f.id + " = value;\n" +
				"  }\n"
		}
	}

  r += "\n  /** @return {!Array<?>} */\n" +
    "  serialize () {\n" +
    "    return [\n" +
    "      this._" + fields[0].id

  for i := 1; i < len(fields); i++ {
    r += ",\n      this._" + fields[i].id
  }

  r += "\n    ];\n" +
    "  }\n"

  r += "\n  /**\n" +
    "   * @param {!Array<?>} serial\n" +
    "   * @return {!" + name + "}\n" +
    "   */\n" +
    "  static restore (serial) {\n" +
    "    return new " + name + " (\n" +
    "      serial[0]"

  for i := 1; i < len(fields); i++ {
    r += ",\n      serial[" + strconv.Itoa(i) + "]"
  }

  r += "\n    );\n"
  r += "  }\n"

	r += "}\n"
	return r
}

func main() {
	if len(os.Args) != 2 {
		help("Missing file name")
		return
	}
	jsfile := os.Args[1]

	tmpfile, err := ioutil.TempFile("", "jstemplate")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	changes := 0
	lines := make([]string, 0)
	templateType := ""
  templateName := ""
	Lines(jsfile, func(l string) {
    lTrim := strings.TrimSpace(l);
		if templateType == "" {
      if strings.HasPrefix(lTrim, "<<<class:") {
				templateType = "class"
        templateName = strings.TrimSpace(lTrim[9:])
				changes++
      } else {
				if _, err := tmpfile.Write([]byte(l + "\n")); err != nil {
					log.Fatal(err)
				}
      }
		} else {
			if lTrim == ">>>" {
				code := ""
				switch templateType {
				case "class":
					code = tClass(templateName, lines)
				default:
					panic("case " + templateType + " is unknown")
				}

				if _, err := tmpfile.Write([]byte(code)); err != nil {
					log.Fatal(err)
				}
				lines = make([]string, 0)
				templateType = ""
        templateName = ""
			} else {
				lines = append(lines, lTrim)
			}
		}
	})
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	if changes > 0 {
		jswfile := Wopen(jsfile)
		defer jswfile.Close()

		Lines(tmpfile.Name(), func(l string) {
			Write(jswfile, l+"\n")
		})
	}

}
