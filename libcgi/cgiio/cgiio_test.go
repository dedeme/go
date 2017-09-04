// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package cgiio

import (
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

func TestCgiio(t *testing.T) {
	u, _ := user.Current()
	dir := filepath.Join(u.HomeDir, ".dmGoApli", "dmGoLib")
	tmp := filepath.Join(dir, "tmp.txt")
	tmp1 := filepath.Join(dir, "tmp1.txt")
	tmp2 := filepath.Join(dir, "tmp2.txt")

	Mkdirs(dir)
	ftmp := OpenWrite(tmp)
	Write(ftmp, "Una\n")
	Write(ftmp, "\n")
	Write(ftmp, "Dos...\ny Tres")
	ftmp.Close()

	ftmp = OpenAppend(tmp)
	Write(ftmp, "\nY un añadido")
	ftmp.Close()

	tx := ReadAll(tmp)
	t.Log(tx)
	if tx != "Una\n\nDos...\ny Tres\nY un añadido" {
		t.Fatal("LineReader gives:\n" + tx)
	}

	ftmp1 := OpenWrite(tmp1)
	Write(ftmp1, tx)
	ftmp1.Close()

	WriteAll(tmp2, tx)

	tx1 := ReadAll(tmp1)
	if tx != tx1 {
		t.Fatal("LineReader gives:\n" + tx1)
	}
	ftmp1.Close()

	tx2 := ""
	Lines(tmp2, func(l string) {
		tx2 += l + "\n"
	})

	if tx != strings.TrimSpace(tx2) {
		t.Fatal("LineReader gives:\n" + tx2)
	}

	Remove(dir)
}
