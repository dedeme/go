// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Utilities to easy input-output
package cgiio

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

// UserDir returns the name of user dir
func HomeDir() string {
	u, _ := user.Current()
	return u.HomeDir
}

// Exists returns true if path actually exists in file system
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// Is directory return true if path exists and is a directory
func IsDirectory(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}

// Mkdir makes a directory
func Mkdir(f string) {
	os.Mkdir(f, os.FileMode(0755))
}

// Mkdirs makes a directory and its parents
func Mkdirs(f string) {
	os.MkdirAll(f, os.FileMode(0755))
}

// TempDir makes a directorio in the temporal directory. If fails throws a
// panic(error),
func TempDir(prefix string) string {
	name, err := ioutil.TempDir("", prefix)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

// List return the list of files of a directory
func List(dir string) []os.FileInfo {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	return fis
}

// TempFile creates a file in 'dir'. If 'dir' is "" file is created in the
// temporal directory.
func TempFile(dir string, prefix string) *os.File {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// Rename changes the name of a file or directory
func Rename(oldname, newname string) {
	err := os.Rename(oldname, newname)
	if err != nil {
		log.Fatal(err)
	}
}

// Remove removes path and all its subdirectories.
func Remove(path string) error {
	return os.RemoveAll(path)
}

// OpenRead opens path for read. If fails throws a panic(error).
func OpenRead(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// Read reads a file completely. (File is open and closed)
func ReadAllBin(path string) []byte {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return bs
}

// Read reads a data file completely. (File is open and closed)
func ReadAll(path string) string {
	return string(ReadAllBin(path))
}

// Lines are read without end of line.
func Lines(path string, f func(s string)) {
	file := OpenRead(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// OpenRead opens path for write. If fails throws a panic(error).
func OpenWrite(path string) *os.File {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// OpenRead opens path for append. If fails throws a panic(error).
func OpenAppend(path string) *os.File {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// WriteAll writes data overwriting 'file'. (File is open and closed)
func WriteAllBin(path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

// WriteAll writes a text overwriting 'file'. (File is open and closed)
func WriteAll(path, text string) {
	WriteAllBin(path, []byte(text))
}

// Write  writes a text in 'file'
func Write(file *os.File, text string) {
	_, err := file.WriteString(text)
	if err != nil {
		log.Fatal(err)
	}
}
