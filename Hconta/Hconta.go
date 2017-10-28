// Copyright 23-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"encoding/json"
	"fmt"
	"github.com/dedeme/go/libcgi"
	"github.com/dedeme/go/libcgi/cgiio"
	"github.com/dedeme/go/libcgi/cryp"
	"os"
	"path"
	"strings"
)

const (
	appName          = "Hconta"
	cgiPath          = "/deme/wwwcgi/dmcgi"
	expiration int64 = 1800 // 1/2 hour
)

func hcontaInit() {
	dir := path.Join(libcgi.Home, "data")

	if !cgiio.Exists(dir) {
		cgiio.Mkdir(dir)
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
			actionsPath := path.Join(libcgi.Home, "data", year+".db")
			rp := make(map[string]interface{})
			if cgiio.Exists(actionsPath) {
				rp["actions"] = cgiio.ReadAll(actionsPath)
			} else {
				rp["actions"] = ""
			}
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
			libcgi.DelSession(sessionId)
			rp := make(map[string]interface{})
			libcgi.Ok(rp)
		case "chpass":
			libcgi.ChangePass(
				data["user"].(string),
				data["pass"].(string),
				data["newPass"].(string))
		default:
			libcgi.Err(r + ": Unknown request")
		}
	}
}
