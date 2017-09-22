// Copyright 15-Sep-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"github.com/dedeme/go/libcgi"
	"github.com/dedeme/go/libcgi/cgiio"
	"path"
)

// Send to client content of file data/conf.db in field "conf"
func sendConf() {
  file := path.Join(libcgi.Home, "data", "conf.db")
	rp := make(map[string]interface{})
	rp["conf"] = cgiio.ReadAll(file)
	libcgi.Ok(rp)
}

// Send to client content of file data/paths.db and data/conf.db in fields
// "paths" and "conf"
func sendConfPaths() {
  pfile := path.Join(libcgi.Home, "data", "paths.db")
  cfile := path.Join(libcgi.Home, "data", "conf.db")
	rp := make(map[string]interface{})
	rp["paths"] = cgiio.ReadAll(pfile)
	rp["conf"] = cgiio.ReadAll(cfile)
	libcgi.Ok(rp)
}

func setConf(data string) {
	file := path.Join(libcgi.Home, "data", "conf.db")
	cgiio.WriteAll(file, data)
	rp := make(map[string]interface{})
	libcgi.Ok(rp)
}

// Send to client content of file path
func sendFile(path string){
  rp := make(map[string]interface{})
  rp["text"] = cgiio.ReadAll(path)
  libcgi.Ok(rp)
}
