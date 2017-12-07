// Copyright 23-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dedeme/go/libcgi"
	"github.com/dedeme/go/libcgi/cgiio"
	"github.com/dedeme/go/libcgi/cryp"
	"os"
	"path"
	"strings"
	"time"
  "net/http"
  "io/ioutil"
)

const (
	appName           = "BolsaData"
	dataVersion       = "201711"
	cgiPath           = "/deme/wwwcgi/dmcgi"
	expiration  int64 = 1800 // 1/2 hour
)

func bolsaDataInit() {
	dir := path.Join(libcgi.Home, "data")

	if !cgiio.Exists(dir) {
		cgiio.Mkdir(dir)
		cgiio.WriteAll(path.Join(dir, "version.txt"),
			appName+"\nData version: "+dataVersion+"\n")
		cgiio.Mkdir(path.Join(libcgi.Home, "tmp"))
		cgiio.Mkdir(path.Join(libcgi.Home, "trash"))
	}
}

func trash() []string {
	fs := cgiio.List(path.Join(libcgi.Home, "trash"))
	r := make([]string, len(fs))
	for i, f := range fs {
		r[i] = f.Name()
	}
	return r
}

func clearTmp() {
	dir := path.Join(libcgi.Home, "tmp")
	cgiio.Remove(dir)
	cgiio.Mkdir(dir)
}

func mkDate() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func mkDate2() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d-%d", t.Year(), t.Month(), t.Day(), t.Unix())
}

func toTrash() {
	cgiio.Zip(
		path.Join(libcgi.Home, "data"),
		path.Join(libcgi.Home, "trash", mkDate2()))
}

func unzip() {
	if err := cgiio.Unzip(
		path.Join(libcgi.Home, "tmp", "back.zip"),
		path.Join(libcgi.Home, "tmp")); err != nil {
		rp := make(map[string]interface{})
		rp["fail"] = "restore:unzip"
		libcgi.Ok(rp)
	}

	version := path.Join(libcgi.Home, "tmp", "data", "version.txt")
	if !cgiio.Exists(version) ||
		!strings.HasPrefix(cgiio.ReadAll(version), appName) {
		clearTmp()
		rp := make(map[string]interface{})
		rp["fail"] = "restore:version"
		libcgi.Ok(rp)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recover:\n%v", r)
			os.Exit(0)
		}
	}()

	libcgi.Init(path.Join(cgiPath, appName), expiration)
	bolsaDataInit()

	rq := os.Args[1]
	ix := strings.Index(rq, ":")

	if ix == -1 { //................................................ CONNECTION
		libcgi.SetKey(rq)
		libcgi.Connect(rq)
	} else if ix == 0 { //...................................... AUTHENTICATION
		key := cryp.Key(appName, libcgi.Klen)
		libcgi.SetKey(key)
		data := strings.Split(
			cryp.Decryp(key, rq[1:]),
			":")
		libcgi.Authentication(data[0], data[1], data[2] == "1")
	} else { //.................................................... NORMAL DATA
		sessionId := rq[:ix]
		pageId, key := libcgi.GetPageIdKey(sessionId)

		if key[0] == '!' {
			libcgi.SetKey(key[1:])
			libcgi.Expired()
		}

		libcgi.SetKey(key)
		jdata := cryp.Decryp(key, rq[ix+1:])
		var data map[string]interface{}
		err := json.Unmarshal([]byte(jdata), &data)
		if err != nil {
			libcgi.Err(err.Error())
		}
		// .................................................. PROCESS NORMAL DATA
		if data["pageId"] != pageId {
			libcgi.Expired()
		}
		r := data["rq"].(string)
		switch r {
		case "getDb":
			dbPath := path.Join(libcgi.Home, "data", "data.db")
			rp := make(map[string]interface{})
			if cgiio.Exists(dbPath) {
				rp["db"] = cgiio.ReadAll(dbPath)
			} else {
				rp["db"] = "{}"
			}
			rp["trash"] = trash()
			libcgi.Ok(rp)
		case "setDb":
			db := data["db"].(string)
			dbPath := path.Join(libcgi.Home, "data", "data.db")
			cgiio.WriteAll(dbPath, db)
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "getQuotes":
			nick := data["nick"].(string)
			dbPath := path.Join(libcgi.Home, "data", nick + ".db")
			rp := make(map[string]interface{})
			if cgiio.Exists(dbPath) {
				rp["quotes"] = cgiio.ReadAll(dbPath)
			} else {
				rp["quotes"] = ""
			}
			libcgi.Ok(rp)
		case "setQuotes":
			nick := data["nick"].(string)
      quotes := data["quotes"].(string)
			dbPath := path.Join(libcgi.Home, "data", nick + ".db")
			cgiio.WriteAll(dbPath, quotes)
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "getInvertia":
      rp := make(map[string]interface{})
      rp["fail"] = ""
      rp["page"] = ""
			url := data["url"].(string)
			resp, err := http.Get(url)
			if err != nil {
        rp["fail"] = err.Error();
				libcgi.Ok(rp)
			}
			body, err := ioutil.ReadAll(resp.Body)
      resp.Body.Close()
			if err != nil {
        rp["fail"] = err.Error();
				libcgi.Ok(rp)
			} else {
        rp["page"] = body
        libcgi.Ok(rp)
      }
		case "logout":
			libcgi.DelSession(sessionId)
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "chpass":
			libcgi.ChangePass(
				data["user"].(string),
				data["pass"].(string),
				data["newPass"].(string))
		case "backup":
			clearTmp()
			name := "BolsaDataBackup" + mkDate() + ".zip"
			cgiio.Zip(
				path.Join(libcgi.Home, "data"),
				path.Join(libcgi.Home, "tmp", name))
			rp := make(map[string]interface{})
			rp["name"] = name
			libcgi.Ok(rp)
		case "restoreStart":
			clearTmp()
			f := cgiio.OpenWrite(path.Join(libcgi.Home, "tmp", "back.zip"))
			f.Close()
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "restoreAppend":
			f := cgiio.OpenAppend(path.Join(libcgi.Home, "tmp", "back.zip"))
			d, _ := base64.StdEncoding.DecodeString(data["data"].(string))
			cgiio.WriteBin(f, d)
			f.Close()
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "restoreAbort":
			clearTmp()
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "restoreEnd":
			unzip()
			toTrash()
			cgiio.Remove(path.Join(libcgi.Home, "data"))
			cgiio.Rename(path.Join(libcgi.Home, "tmp", "data"),
				path.Join(libcgi.Home, "data"))
			clearTmp()
			rp := make(map[string]interface{})
			rp["fail"] = ""
			libcgi.Ok(rp)
		case "autorestore":
			toTrash()
			cgiio.Remove(path.Join(libcgi.Home, "data"))
			cgiio.Unzip(
				path.Join(libcgi.Home, "backups", data["file"].(string)),
				path.Join(libcgi.Home))
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "clearTrash":
			cgiio.Remove(path.Join(libcgi.Home, "trash"))
			cgiio.Mkdir(path.Join(libcgi.Home, "trash"))
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "restoreTrash":
			toTrash()
			cgiio.Remove(path.Join(libcgi.Home, "data"))
			cgiio.Unzip(
				path.Join(libcgi.Home, "trash", data["file"].(string)),
				path.Join(libcgi.Home))
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		default:
			libcgi.Err(r + ": Unknown request")
		}
	}
}
