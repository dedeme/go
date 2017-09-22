// Copyright 10-Sep-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Server side of documentation generator for JavaScript.
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
	appName          = "JsDoc"
	cgiPath          = "/deme/wwwcgi/dmcgi"
	expiration int64 = 3600 // 1 hour
)

func jsdocInit() {
	dir := path.Join(libcgi.Home, "data")
	if !cgiio.Exists(dir) {
		cgiio.Mkdir(dir)

		cgiio.WriteAll(path.Join(dir, "paths.db"), "")
		cgiio.WriteAll(path.Join(dir, "conf.db"), "")
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
	jsdocInit()

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
		page := data["page"].(string)
		switch page {
		case "Main":
			sendConf()
		case "Paths":
      if data["pageId"] != pageId {
        libcgi.Expired()
      }
      rq := data["rq"].(string)
			switch rq {
			case "get":
				sendConfPaths()
			case "setConf":
				setConf(data["conf"].(string))
			case "exists":
				pathsExists(data["paths"].([]interface{}))
      case "setPaths":
				pathsSetPaths(data["paths"].(string))
			default:
        libcgi.Err(rq + ": Unknown option in page Paths")
			}
    case "Chpass":
      if data["pageId"] != pageId {
        libcgi.Expired()
      }
      rq := data["rq"].(string)
			switch rq {
			case "get":
				sendConf()
			case "setConf":
				setConf(data["conf"].(string))
      case "change":
				libcgi.ChangePass(
            data["user"].(string),
            data["pass"].(string),
            data["newPass"].(string))
			default:
        libcgi.Err(rq + ": Unknown option in page Paths")
			}
    case "Index":
      if data["pageId"] != pageId {
        libcgi.Expired()
      }
      rq := data["rq"].(string)
			switch rq {
			case "get":
				sendConfPaths()
			case "setConf":
				setConf(data["conf"].(string))
      case "path":
				indexTree(data["path"].(string))
			default:
        libcgi.Err(rq + ": Unknown option in page Paths")
			}
    case "Module":
      if data["pageId"] != pageId {
        libcgi.Expired()
      }
      rq := data["rq"].(string)
			switch rq {
			case "get":
				sendConfPaths()
			case "setConf":
				setConf(data["conf"].(string))
      case "code":
        sendFile(data["path"].(string))
			default:
        libcgi.Err(rq + ": Unknown option in page Paths")
			}
		case "Logout":
      libcgi.DelSession(sessionId)
			sendConf() // Returns conf file
		default:
			libcgi.Err(page + ": Unknown page")
		}
	}
}
