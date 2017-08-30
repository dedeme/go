package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var enc = base64.StdEncoding

type rpT struct {
	Error string
	Page  string
}

func Err(msg string) {
	st, err := json.Marshal(msg)
	if err == nil {
		fmt.Print(enc.EncodeToString(st))
	} else {
		fmt.Print("{\"Error\" : \"Error in Err\"}")
	}
}

func Ok(rp map[string]interface{}) {
	st, err := json.Marshal(rp)
	if err == nil {
		fmt.Print(enc.EncodeToString(st))
	} else {
		fmt.Print("{\"Error\" : \"Error in Ok\"}")
	}
}

func readPage(page string) string {
	resp, err := http.Get(page)
	if err != nil {
		Err(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Err(err.Error())
	}
	return string(body)
}

func main() {
	var rq string
	if rq0, err := enc.DecodeString(os.Args[1]); err == nil {
		rq = string(rq0)
	} else {
		Err(err.Error())
	}

	page := rq

	switch rq {
	case "publico":
		page = readPage("http://www.publico.com")
	case "diario":
		page = readPage("http://www.eldiario.es")
	case "20minutos":
		page = readPage("http://www.20minutos.es")
	case "meneame":
		page1 := readPage("https://www.meneame.net")
		page2 := readPage("https://www.meneame.net/?page2")
		page = page1 + "\n" + page2
	default:
		page = "'" + rq + "' Unknown vendor"
	}

	rp := make(map[string]interface{})
	rp["Error"] = ""
	rp["Page"] = page
	Ok(rp)
}
