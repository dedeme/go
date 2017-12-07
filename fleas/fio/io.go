// Copyright 27-Nov-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package fio

import (
	"encoding/json"
	"fmt"
	"github.com/dedeme/go/fleas/quote"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
	"time"
)

func mkDate() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func home() string {
	usr, err := user.Current()
	if err != nil {
		panic("fio.home: " + err.Error())
	}
	return path.Join(usr.HomeDir, ".dmGoApp", "fleas")
}

func dataDir() string {
	return path.Join(home(), "data")
}

// Exists returns true if path actually exists in file system
func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// Mkdirs makes a directory and its parents
func mkdirs(f string) {
	os.MkdirAll(f, os.FileMode(0755))
}

// Rename changes the name of a file or directory
func rename(oldname, newname string) {
	err := os.Rename(oldname, newname)
	if err != nil {
		panic(err)
	}
}

// write writes a text overwriting 'file'. (File is open and closed)
func write(path, text string) {
	err := ioutil.WriteFile(path, []byte(text), 0755)
	if err != nil {
		panic("fio.write: " + err.Error())
	}
}

// Read reads a data file completely. (File is open and closed)
func read(path string) string {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic("fio.read: " + err.Error())
	}
	return string(bs)
}

// Ini initializes data base
func Ini() {
	h := home()
	if !exists(h) {
		mkdirs(h)
	}
	d := dataDir()
	if !exists(d) {
		mkdirs(d)
	} else {
		tmp := path.Join(h, "tmp")
		if exists(tmp) {
			os.RemoveAll(tmp)
		}
	}
	backup := path.Join(h, "backup")
	if exists(backup) {
		os.RemoveAll(backup)
	}
	mkdirs(backup)

	version := path.Join(d, "version.txt")
	if !exists(version) {
		write(version, "fleas\nData version: 201711\n")
	}
	msg := path.Join(d, "msg.txt")
	if exists(msg) {
		os.Remove(msg)
	}
	db := path.Join(d, "db.json")
	if !exists(db) {
		write(db, "")
	}
	bolsaData := path.Join(h, "bolsaData")
	if !exists(bolsaData) {
		msg := "fio.Ini: Symbolic link to 'BolsaData' is missing"
		panic(msg)
	}
}

func ReadDb () string {
	return read(path.Join(dataDir(), "db.json"))
}

func WriteDb (db string) {
	write(path.Join(dataDir(), "db.json"), db)
}

func WriteMsg(msg string) {
	fmt.Println(msg)
	if exists(dataDir()) {
		f := path.Join(dataDir(), "msg.txt")
		if exists(f) {
			msg = read(f) + msg
		}
		write(f, msg+"\n")
	}
}

func Start() {
  write(path.Join(dataDir(), "start.lock"), "")
}

func IsStarted () bool {
  return exists(path.Join(dataDir(), "start.lock"))
}

func Stop () {
  if exists(dataDir()) {
    write(path.Join(dataDir(), "stop.lock"), "")
  }
}

func IsStoped () bool {
  return exists(path.Join(dataDir(), "stop.lock"))
}

func Force () {
  f := path.Join(dataDir(), "start.lock")
  if exists(f) {
    os.Remove(f)
  }
  f = path.Join(dataDir(), "stop.lock")
  if exists(f) {
    os.Remove(f)
  }
}

func Reset() {
	os.RemoveAll(dataDir())
}

func Backup() {
	backup := path.Join(home(), "backup")
	if exists(backup) {
		os.RemoveAll(backup)
	}
	mkdirs(backup)

	target := path.Join(home(), "backup", "Fleas"+mkDate()+".zip")
	if err := Zip(dataDir(), target); err != nil {
		panic(err)
	}
}

func Restore(fileName string) {
	source := path.Join(home(), "backup", fileName)
	if !exists(source) {
		panic("File '" + fileName + "' does not exist")
	}
	target := path.Join(home(), "tmp")
	if exists(target) {
		os.RemoveAll(target)
	}
	mkdirs(target)

	if err := Unzip(source, target); err != nil {
		panic(err)
	}
	vFile := path.Join(target, "data", "version.txt")
	if !exists(vFile) {
		os.RemoveAll(target)
		panic("'" + fileName + "' does not contain 'version.txt'")
	}
	vData := read(vFile)
	if !strings.HasPrefix(vData, "fleas\n") {
		os.RemoveAll(target)
		panic("'version.txt' in '" + fileName + "' is not valid")
	}
	os.RemoveAll(dataDir())
	rename(path.Join(target, "data"), dataDir())
}

// ListNicks returns a unsorted list of every cia nick (1) and Ibex cia nick (2)
func ListNicks() ([]string, map[string]struct{}) {
	dataS := read(path.Join(home(), "bolsaData", "data.db"))
	var data map[string]interface{}
	err := json.Unmarshal([]byte(dataS), &data)
	if err != nil {
		panic("fio.ListNicks(1): " + err.Error())
	}

	var r1 []string
	for nick, st := range data["status"].(map[string]interface{}) {
		if st.(string) == "" {
			r1 = append(r1, nick)
		}
	}

	r2 := make(map[string]struct{})
	for nick, st := range data["ibex"].(map[string]interface{}) {
		if st.(bool) {
			r2[nick] = struct{}{}
		}
	}

	return r1, r2
}

// Quotes returns market quotes sorted by date
func Quotes(nicks []string) []map[string]quote.Quote {
	qNumber := len(strings.Split(
		read(path.Join(home(), "bolsaData", nicks[0]+".db")), "\n"))
	r := make([]map[string]quote.Quote, qNumber)
	for i := range r {
		r[i] = make(map[string]quote.Quote)
	}
	for i := range nicks {
		nick := nicks[i]
		qs := strings.Split(read(path.Join(home(), "bolsaData", nick+".db")),
			"\n")
		for j, q := range qs {
			var qdata [6]interface{}
			qparts := strings.Split(q, ":")
			qdata[0] = qparts[0]
			for k := 1; k < 6; k++ {
				n, err := strconv.ParseFloat(qparts[k], 32)
				if err != nil {
					n = -1
				}
				qdata[k] = float32(n)
			}
			r[qNumber-j-1][nick] = qdata
		}
	}
	return r
}
