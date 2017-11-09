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
  "sort"
)

const (
	appName           = "Hconta"
	dataVersion       = "201709"
	cgiPath           = "/deme/wwwcgi/dmcgi"
	expiration  int64 = 1800 // 1/2 hour
)

func hcontaInit() {
	dir := path.Join(libcgi.Home, "data")

	if !cgiio.Exists(dir) {
		cgiio.Mkdir(dir)
		cgiio.WriteAll(path.Join(dir, "version.txt"),
			"Hconta\nData version: "+dataVersion+"\n")
		cgiio.Mkdir(path.Join(libcgi.Home, "tmp"))
		cgiio.Mkdir(path.Join(libcgi.Home, "backups"))
		cgiio.Mkdir(path.Join(libcgi.Home, "trash"))
	}
}

func backups() []string {
	fs := cgiio.List(path.Join(libcgi.Home, "backups"))
	r := make([]string, len(fs))
	for i, f := range fs {
		r[i] = f.Name()
	}
	return r
}

func trash() []string {
	fs := cgiio.List(path.Join(libcgi.Home, "trash"))
	r := make([]string, len(fs))
	for i, f := range fs {
		r[i] = f.Name()
	}
	return r
}

func filterBackups() {
  t0 := time.Now()
  d2 := fmt.Sprintf("%d%02d%02d", t0.Year() - 1, t0.Month(), t0.Day())
  t1 := t0.AddDate(0, 0, -7)
  d1:= fmt.Sprintf("%d%02d%02d", t1.Year(), t1.Month(), t1.Day())
  fs := backups();
  sort.Strings(fs);
  previous := "        "
  for _, f := range fs {
    if f < d2 {
      if previous[0:4] == f[0:4] {
        cgiio.Remove(path.Join(libcgi.Home, "backups", previous))
      }
    } else if (f < d1) {
      if previous[0:6] == f[0:6] {
        cgiio.Remove(path.Join(libcgi.Home, "backups", previous))
      }
    }
    previous = f
  }
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
		!strings.HasPrefix(cgiio.ReadAll(version), "Hconta") {
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
	hcontaInit()

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

		if key == "" {
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
		case "getConf":
			confPath := path.Join(libcgi.Home, "data", "conf.db")
			rp := make(map[string]interface{})
			if cgiio.Exists(confPath) {
				rp["conf"] = cgiio.ReadAll(confPath)
			} else {
				rp["conf"] = ""
			}
			libcgi.Ok(rp)
		case "getDb":
			year := data["year"].(string)
			dbPath := path.Join(libcgi.Home, "data", year+".db")
			rp := make(map[string]interface{})
			if cgiio.Exists(dbPath) {
				rp["db"] = cgiio.ReadAll(dbPath)
			} else {
				rp["db"] = ""
			}
			rp["backups"] = backups()
			rp["trash"] = trash()
			libcgi.Ok(rp)
		case "setConf":
			conf := data["conf"].(string)
			confPath := path.Join(libcgi.Home, "data", "conf.db")
			cgiio.WriteAll(confPath, conf)
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "setDb":
			year := data["year"].(string)
			actions := data["db"].(string)
			actionsPath := path.Join(libcgi.Home, "data", year+".db")
			cgiio.WriteAll(actionsPath, actions)
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "logout":
			cgiio.Zip(
				path.Join(libcgi.Home, "data"),
				path.Join(libcgi.Home, "backups", mkDate()))
      filterBackups()
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
			name := "HcontaBackup" + mkDate() + ".zip"
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
