// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package cryp

import (
	"fmt"
	"testing"
)

func eq(actual, expected string) string {
	if actual != expected {
		return fmt.Sprintf("\nActual  : %s\nExpected: %s\n", actual, expected)
	}
	return ""
}

func TestGenK(t *testing.T) {
	ac := fmt.Sprintf("%d", len(GenK(12)))
	if r := eq(ac, "12"); r != "" {
		t.Fatal(r)
	}
	ac = fmt.Sprintf("%d", len(GenK(6)))
	if r := eq(ac, "6"); r != "" {
		t.Fatal(r)
	}
}

func TestKey(t *testing.T) {
	if r := eq(Key("deme", 6), "wiWTB9"); r != "" {
		t.Fatal(r)
	}
	if r := eq(Key("Generaro", 5), "Ixy8I"); r != "" {
		t.Fatal(r)
	}
	if r := eq(Key("Generara", 5), "0DIih"); r != "" {
		t.Fatal(r)
	}
}

func TestCryp(t *testing.T) {
	r := eq(Cryp("deme", "Cañón€%ç"), "v12ftuzYeq2Xz7q7tLe8tNnHtqY=")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("deme", Cryp("deme", "Cañón€%ç")), "Cañón€%ç")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("deme", Cryp("deme", "1")), "1")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("deme", Cryp("deme", "")), "")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("", Cryp("", "Cañón€%ç")), "Cañón€%ç")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("", Cryp("", "1")), "1")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("", Cryp("", "")), "")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("abc", Cryp("abc", "01")), "01")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("abcd", Cryp("abcd", "11")), "11")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("abc", Cryp("abc", "")), "")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("c", Cryp("c", "a")), "a")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("xxx", Cryp("xxx", "ab c")), "ab c")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(Decryp("abc", Cryp("abc", "\n\ta€b c")), "\n\ta€b c")
	if r != "" {
		t.Fatal(r)
	}
}
