// Copyright 23-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package libcgi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dedeme/go/libcgi/cgiio"
	"github.com/dedeme/go/libcgi/cryp"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const Klen = 300
const tNoExpiration int64 = 2592000 // seconds == 30 days
const demeKey = "nkXliX8lg2kTuQSS/OoLXCk8eS4Fwmc+N7l6TTNgzM1vdKewO0cjok51vcdl" +
	"OKVXyPu83xYhX6mDeDyzapxL3dIZuzwyemVw+uCNCZ01WDw82oninzp88Hef" +
	"bn3pPnSMqEaP2bOdX+8yEe6sGkc3IO3e38+CqSOyDBxHCqfrZT2Sqn6SHWhR" +
	"KqpJp4K96QqtVjmXwhVcST9l+u1XUPL6K9HQfEEGMGcToMGUrzNQxCzlg2g+" +
	"Hg55i7iiKbA0ogENhEIFjMG+wmFDNzgjvDnNYOaPTQ7l4C8aaPsEfl3sugiw"

var b64 = base64.StdEncoding

var fkey string
var Home string
var key string
var tExpiration int64

// Init initializes libcgi. 'texpiration' is in seconds.
//
// When application is initialized for the first time, 'home' must no exist.
func Init(home string, texpiration int64) {
	fkey = cryp.Key(demeKey, len(demeKey))
	Home = home
	key = fkey // key must be initialized with 'setKey' before call Ok or Err
	tExpiration = texpiration

	if !cgiio.Exists(home) {
		cgiio.Mkdirs(home)
		cgiio.WriteAll(path.Join(home, "users.db"), "")
		cgiio.WriteAll(path.Join(home, "sessions.db"), "")
		putUser("admin", demeKey, "0")
	}
}

// SetKey sets the key which 'Ok' and 'Err' will use
func SetKey(k string) {
	key = k
}

// If expiration is false tNoExpiration is used
func addSession(sessionId, pageId, key string, expiration bool) {
	lapse := tNoExpiration
	if expiration {
		lapse = tExpiration
	}
	time := time.Now().Unix() + lapse

	f := cgiio.OpenAppend(path.Join(Home, "sessions.db"))
	cgiio.Write(f,
		cryp.Cryp(fkey,
			sessionId+":"+
				pageId+":"+
				key+":"+
				strconv.FormatInt(time, 10)+":"+
				strconv.FormatInt(lapse, 10))+
			"\n")
	f.Close()
}

// In session.db:
//    Each line is a record
//    Each field is a B64 string, separated by ":"
//    Fields are: sessionId:PageId:key:time:lapse
// Set 'pageId' to "" for avoiding to change it.
// If sessionId is expired key returns "!" + key
func readSession(sessionId, pageId string) (rpageId, key string) {
	path := path.Join(Home, "sessions.db")
	now := time.Now().Unix()
	newSs := ""
	rpageId = ""
	key = ""

	cgiio.Lines(path, func(l string) {
		if l == "" {
			return
		}

		es := strings.Split(cryp.Decryp(fkey, l), ":")
		t, err := strconv.ParseInt(es[3], 10, 64)
		if err != nil {
			Err(err.Error())
		}
		if now < t {
			if es[0] == sessionId {
				if pageId != "" {
					es[1] = pageId
				}
				t, err := strconv.ParseInt(es[4], 10, 64)
				if err != nil {
					Err(err.Error())
				}
				es[3] = strconv.FormatInt(now+t, 10)
				newSs += cryp.Cryp(
					fkey,
					es[0]+":"+es[1]+":"+es[2]+":"+es[3]+":"+es[4]) + "\n"
				rpageId = es[1]
				key = es[2]
			} else {
				newSs += l + "\n"
			}
		} else {
			if es[0] == sessionId {
				key = "!" + es[2]
			}
		}
	})

	cgiio.WriteAll(path, newSs)

	return
}

func DelSession(sessionId string) {
	path := path.Join(Home, "sessions.db")
	newSs := ""
	cgiio.Lines(path, func(l string) {
		if l == "" {
			return
		}
		es := strings.Split(cryp.Decryp(fkey, l), ":")
		if es[0] != sessionId {
			newSs += l + "\n"
		}
	})

	cgiio.WriteAll(path, newSs)

	return
}

func readUsers() []string {
	r := make([]string, 0)
	cgiio.Lines(path.Join(Home, "users.db"), func(l string) {
		if l != "" {
			r = append(r, cryp.Decryp(fkey, l))
		}
	})
	return r
}

func writeUsers(users []string) {
	f := cgiio.OpenWrite(path.Join(Home, "users.db"))
	for _, l := range users {
		cgiio.Write(f, cryp.Cryp(fkey, l)+"\n")
	}
	f.Close()
}

// Returns user level or "" if checking failed
func checkUser(user, key string) string {
	rkey := cryp.Key(key, Klen)
	for _, l := range readUsers() {
		es := strings.Split(l, ":")
		if es[0] == user && es[1] == rkey {
			return es[2]
		}
	}
	return ""
}

func changeUserPass(user, key string) bool {
	rkey := cryp.Key(key, Klen)
	r := false
	newUsers := make([]string, 0)
	for _, l := range readUsers() {
		es := strings.Split(l, ":")
		if es[0] == user {
			r = true
			newUsers = append(newUsers, es[0]+":"+rkey+":"+es[2])
		} else {
			newUsers = append(newUsers, l)
		}
	}
	writeUsers(newUsers)
	return r
}

func changeUserLevel(user, level string) bool {
	r := false
	newUsers := make([]string, 0)
	for _, l := range readUsers() {
		es := strings.Split(l, ":")
		if es[0] == user {
			r = true
			newUsers = append(newUsers, es[0]+":"+es[1]+":"+level)
		} else {
			newUsers = append(newUsers, l)
		}
	}
	writeUsers(newUsers)
	return r
}

func delUser(user string) bool {
	r := false
	newUsers := make([]string, 0)
	for _, l := range readUsers() {
		es := strings.Split(l, ":")
		if es[0] == user {
			r = true
		} else {
			newUsers = append(newUsers, l)
		}
	}
	writeUsers(newUsers)
	return r
}

// Removes an adds an user.
// "0" is the admin level.
func putUser(user, key, level string) {
	delUser(user)
	users := readUsers()
	users = append(users, user+":"+cryp.Key(key, Klen)+":"+level)
	writeUsers(users)
}

// Returns pageId and key. If conection failed both are ""
func GetPageIdKey(sessionId string) (pageId, key string) {
	return readSession(sessionId, "")
}

// Send to client pageId and key. If conection failed both are ""
func Connect(sessionId string) {
	rp := make(map[string]interface{})
	rp["pageId"], rp["key"] = readSession(sessionId, cryp.GenK(Klen))
	Ok(rp)
}

// Send to client level, key, pageId and sessionI, If authentication failed
// every value is "".
func Authentication(user, key string, expiration bool) {
	rp := make(map[string]interface{})

	if level := checkUser(user, key); level != "" {
		sessionId := cryp.GenK(Klen)
		pageId := cryp.GenK(Klen)
		key := cryp.GenK(Klen)
		addSession(sessionId, pageId, key, expiration)
		rp["level"] = level
		rp["key"] = key
		rp["pageId"] = pageId
		rp["sessionId"] = sessionId
	} else {
		rp["level"] = ""
		rp["key"] = ""
		rp["pageId"] = ""
		rp["sessionId"] = ""
	}
	Ok(rp)
}

// Admin level is "0". Send to client ok
func AddUser(admin, akey, user, ukey, level string) {
	rp := make(map[string]interface{})
	if alevel := checkUser(admin, akey); alevel == "0" {
		putUser(user, ukey, level)
		rp["ok"] = true
	} else {
		rp["ok"] = false
	}
	Ok(rp)
}

// Admin level is "0". Send to client ok
func DelUser(admin, akey, user string) {
	rp := make(map[string]interface{})
	if alevel := checkUser(admin, akey); alevel == "0" {
		if delUser(user) {
			rp["ok"] = true
		} else {
			rp["ok"] = false
		}
	} else {
		rp["ok"] = false
	}
	Ok(rp)
}

// Admin level is "0". Send to client ok
func ChangeLevel(admin, akey, user, level string) {
	rp := make(map[string]interface{})
	if alevel := checkUser(admin, akey); alevel == "0" {
		if changeUserLevel(user, level) {
			rp["ok"] = true
		} else {
			rp["ok"] = false
		}
	} else {
		rp["ok"] = false
	}
	Ok(rp)
}

// Send to clien ok
func ChangePass(user, key, newKey string) {
	rp := make(map[string]interface{})
	rp["ok"] = false
	if checkUser(user, key) != "" {
		if changeUserPass(user, newKey) {
			rp["ok"] = true
		} else {
			rp["ok"] = false
		}
	}
	Ok(rp)
}

// Err sends an error message to client. The message is a JSON object of type:
//    {error:msg}
func Err(msg string) {
	r := make(map[string]interface{})
	r["error"] = msg
	rp, err := json.Marshal(r)
	if err == nil {
		fmt.Print(cryp.Cryp(key, string(rp)))
	} else {
		fmt.Print("{\"${Error}\" : \"Error in Err\"}")
	}
	os.Exit(0)
}

// Expired sends a expiration message to client. The message is a JSON object:
//   {expired:true}
func Expired() {
	r := make(map[string]interface{})
	r["expired"] = true
	Ok(r)
}

// Ok sends a response ('rp') to client. It adds a field rp["error"] = "".
func Ok(rp map[string]interface{}) {
	rp["error"] = ""
	js, err := json.Marshal(rp)
	if err == nil {
		fmt.Print(cryp.Cryp(key, string(js)))
	} else {
		fmt.Print("{\"${error}\" : \"Error in Ok\"}")
	}
	os.Exit(0)
}
