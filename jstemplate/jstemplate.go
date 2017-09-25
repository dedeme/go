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

Modify a js file changing templates to javascript code.

If a fail happens the js file will not be modified.

Templates

Valid templates are:

class: Generates a seriazable class.
  <<<class:className
  field1:type
  field2:type
  @field3:type
  ...
  >>>
  field is a get field
  @field is a get/set field

func: Generates a function skeleton
  <<<func:functionName
  field1:type
  field2:type
  ...
  returnType
  >>>
  fields and return are optional.
  If returnType is missing is value is set to 'void'

meth: Generates a method skeleton
  <<<meth:methodName
  field1:type
  field2:type
  ...
  returnType
  >>>
  fields and return are optional.
  If returnType is missing is value is set to 'void'

static: Generates a static method skeleton
  <<<static:staticMethodName
  field1:type
  field2:type
  ...
  returnType
  >>>
  fields and return are optional.
  returnType must be the last line.
  If returnType is missing, is value is set to 'void'.

vars: Declares 'let' variables out of class. They will be linked with <<<link:
  <<<vars:
  field1:type
  field2:type
  @field2:type
  ...
  >>>
  field is a get field
  @field is a get/set field

consts: Declares constants out of class. They will be linked with <<<link:
  <<<consts:
  field1:type
  field2:type
  ...
  >>>

pars: Declares 'this_' variables. They will be linked with <<<link:
  <<<pars:
  field1:type
  field2:type
  @field2:type
  ...
  >>>
  field is a get field
  @field is a get/set field

<<<links: Generates links to all the vars, consts and pars previously declared
  <<<links:
  >>>

*/
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
		"  jstemplate file\n" +
		"For example:\n" +
		"  jstemplate \"src/main.js\"\n")
}

type par struct {
	sta     bool
	writing bool
	id      string
	tp      string
}

var pars = make([]par, 0)

func tClass(name string, lines []string) string {
	type field struct {
		id      string
		writing bool
		tp      string
	}
	fields := make([]field, 0)
	for _, l := range lines {
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
			"  " + f.id + " () {\n" +
			"    return this._" + f.id + ";\n" +
			"  }\n"

		if f.writing {
			id := f.id
			r += "\n  /** @param {" + f.tp + "} value */\n" +
				"  set" + strings.ToUpper(id[0:1]) + id[1:] + " (value) {\n" +
				"    this._" + id + " = value;\n" +
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

func tFunc(kind, name string, lines []string) string {
	cm := "  /**\n"
	js := "  static " + name + " ("
	if kind == "func" {
		js = "  function " + name + " ("
	} else if kind == "meth" {
		js = "  " + name + " ("
	}
	withReturn := false
	i := 0
	for ; i < len(lines); i++ {
		l := lines[i]
		if l == "" {
			continue
		}
		ix := strings.Index(l, ":")
		if ix == -1 {
			cm += "   * @return {" + l + "}\n"
      i++;
			withReturn = true
      break;
		}
		id := strings.TrimSpace(l[:ix])
		tp := strings.TrimSpace(l[ix+1:])
		if id == "" {
			panic("Field name is missing in '" + l + "'")
		}
		if tp == "" {
			panic("Field type is missing in '" + l + "'")
		}
		cm += "   * @param {" + tp + "} " + id + "\n"
		js = js + id + ", "
	}

	for ; i < len(lines); i++ {
		l := lines[i]
		if l != "" {
			panic("Field after return in '" + l + "'")
		}
	}

	if strings.HasSuffix(js, ", ") {
		js = js[0 : len(js)-2]
	}
	if !withReturn {
		cm += "   * @return {void}\n"
	}

	cm += "   */\n"
	js += ") {\n  }\n"

	return cm + js
}

func tVars(kind string, lines []string) string {
	r := ""
	sta := true
	if kind == "pars" {
		sta = false
	}
	for _, l := range lines {
		if l != "" {
			writing := false
			if l[0] == '@' {
				if kind == "consts" {
					panic("Value constant with '@' in '" + l + "'")
				}
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
			pars = append(pars, par{sta, writing, id, tp})
			val := "null"
			if tp == "boolean" {
				val = "false"
			} else if tp == "number" {
				val = "0"
			} else if tp == "string" {
				val = "\"\""
			}
			if kind == "vars" {
				r += "  /** @type {" + tp + "} */\n" +
					"  let " + id + " = " + val + ";\n"
			} else if kind == "consts" {
				r += "  /** @const {" + tp + "} */\n" +
					"  const " + id + " = " + val + ";\n"
			} else {
				r += "    /** @type {" + tp + "} */\n" +
					"    this._" + id + " = " + val + ";\n"
			}
		}
	}

	return r + "\n"
}

func tLinks(lines []string) string {
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if l != "" {
			panic("Value in links '" + l + "'")
		}
	}

	r := ""
	for _, p := range pars {
		r += "  /** @return {" + p.tp + "} */\n" +
			"  "
		if p.sta {
			r += "static "
		}
		r += p.id + " () {\n" +
			"    return "
		if !p.sta {
			r += "this._"
		}
		r += p.id + ";\n" +
			"  }\n\n"
		if p.writing {
			id2 := strings.ToUpper(p.id)
			if len(p.id) > 1 {
				id2 = id2[0:1] + p.id[1:]
			}
			r += "  /**\n" +
				"   * @param {" + p.tp + "} value\n" +
				"   * @return {void}\n" +
				"   */\n" +
				"  "
			if p.sta {
				r += "static "
			}
			r += "set" + id2 + " (value) {\n"
			r += "    "
			if !p.sta {
				r += "this._"
			}
			r += p.id + " = value;\n" +
				"  }\n\n"
		}
	}

	pars = make([]par, 0)
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
		lTrim := strings.TrimSpace(l)
		if templateType == "" {
			if strings.HasPrefix(lTrim, "<<<class:") {
				templateType = "class"
				templateName = strings.TrimSpace(lTrim[9:])
				changes++
			} else if strings.HasPrefix(lTrim, "<<<func:") {
				templateType = "func"
				templateName = strings.TrimSpace(lTrim[8:])
				changes++
			} else if strings.HasPrefix(lTrim, "<<<meth:") {
				templateType = "meth"
				templateName = strings.TrimSpace(lTrim[8:])
				changes++
			} else if strings.HasPrefix(lTrim, "<<<static:") {
				templateType = "static"
				templateName = strings.TrimSpace(lTrim[10:])
				changes++
			} else if strings.HasPrefix(lTrim, "<<<vars:") {
				templateType = "vars"
				changes++
			} else if strings.HasPrefix(lTrim, "<<<consts:") {
				templateType = "consts"
				changes++
			} else if strings.HasPrefix(lTrim, "<<<pars:") {
				templateType = "pars"
				changes++
			} else if strings.HasPrefix(lTrim, "<<<links:") {
				templateType = "links"
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
				case "func":
					code = tFunc("func", templateName, lines)
				case "meth":
					code = tFunc("meth", templateName, lines)
				case "static":
					code = tFunc("static", templateName, lines)
				case "vars":
					code = tVars("vars", lines)
				case "consts":
					code = tVars("consts", lines)
				case "pars":
					code = tVars("pars", lines)
				case "links":
					code = tLinks(lines)
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
