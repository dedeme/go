// Copyright 12-Sep-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"github.com/dedeme/go/libcgi"
	"github.com/dedeme/go/libcgi/cgiio"
	"path"
)

// Set content of file data/paths.db
func pathsSetPaths(data string) {
	file := path.Join(libcgi.Home, "data", "paths.db")
	cgiio.WriteAll(file, data)
	rp := make(map[string]interface{})
	libcgi.Ok(rp)
}

// Send a map id->bool with values true o false depending on if the path
// exists or not. (Paths are absolute)
func pathsExists(idPaths []interface{}) {
	data := make(map[string]bool)
	for _, idPath := range idPaths {
    idP := idPath.([]interface{})
    id := idP[0].(string)
    path := idP[1].(string)
		data[id] = cgiio.IsDirectory(path)
	}
	rp := make(map[string]interface{})
	rp["paths"] = data
	libcgi.Ok(rp)
}
