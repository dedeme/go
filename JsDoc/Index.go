// Copyright 13-Sep-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"bufio"
	"bytes"
	"github.com/dedeme/go/libcgi"
	"github.com/dedeme/go/libcgi/cgiio"
	"log"
	"path"
	"strings"
)

func readLine(buffer *bytes.Buffer, l string, status int) int {
	switch status {
	case 1:
		if strings.HasPrefix(l, "*") {
			l = strings.TrimSpace(l[1:])
		}
		if l == "" {
			return 1
		}
		ix := strings.Index(l, "*/")
		if ix == -1 {
			buffer.WriteString(l)
			return -1
		}
		buffer.WriteString(l[:ix])
		return -1
	case 2:
		ix := strings.Index(l, "*/")
		if ix != -1 {
			return readLine(buffer, strings.TrimSpace(l[ix+2:]), 0)
		}
		return 2
	default:
		if strings.HasPrefix(l, "///") {
			buffer.WriteString(l[3:])
			return -1
		}
		if strings.HasPrefix(l, "/**") {
			return readLine(buffer, strings.TrimSpace(l[2:]), 1)
		}
		if strings.HasPrefix(l, "//") || l == "" {
			return 0
		}
		if strings.HasPrefix(l, "/*") {
			status = 3
			return readLine(buffer, strings.TrimSpace(l[2:]), 2)
		}
		return -1
	}
}

func readFile(path string) string {
	file := cgiio.OpenRead(path)
	defer file.Close()

	var buffer bytes.Buffer
	status := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		status = readLine(&buffer, strings.TrimSpace(scanner.Text()), status)
		if status == -1 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	r := strings.TrimSpace(buffer.String())
	ix := strings.Index(r, ".")
	if ix == -1 {
		return r
	}
	return r[:ix]
}

func tree(id, pt string) [3]interface{} {
	elements := make([]interface{}, 0)
	for _, f := range cgiio.List(pt) {
		fName := f.Name()
		p := path.Join(pt, fName)
		if cgiio.IsDirectory(p) {
			elements = append(elements, tree(fName, p))
		} else if strings.HasSuffix(fName, ".js") {
			elements = append(elements, [3]interface{}{fName, readFile(p), nil})
		}
	}
	return [3]interface{}{id, nil, elements}
}

func indexTree(path string) {
	rp := make(map[string]interface{})
	if cgiio.IsDirectory(path) {
		rp["tree"] = tree("root", path)
	} else {
		rp["tree"] = [3]interface{}{"root", nil, nil}
	}
	libcgi.Ok(rp)
}
