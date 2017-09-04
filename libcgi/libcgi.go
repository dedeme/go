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
  "time"
)

const (
	kGlobal = "var expiration = persistent ? 2592000 : 900;"

	uDb    = "users.db"
	uPass  = 0
	uLevel = 1

	ssDb      = "sessions.db"
	ssUser    = 0
	ssLevel   = 1
	ssTime    = 2
	ssIncTime = 3
	ssPageId  = 4

	// Request field to identify page source
	Page = "page"
	// (Reserved) Request field to send the page identifier
	PageId = "pageId"
	// (Reserved) Request field to send the session identifier and to store it
	SessionId = "sessionId"
	// (Reserved) Request field to indicate expiration time
	Expiration = "expirationTime"
	// (Reserved) Request field to indicate user name
	User = "user"
	// (Reserved) Request field to indicate user password
	Pass = "pass"
	// (Reserved) Request field to indicate old user password
	OldPass = "oldPass"
	// (Reserved) Response field for errors
	Error = "error"
	// (Reserved) Response field for session-page control
	SessionOk = "sessionOk"
	// (Reserved) Response field for change-password control
	ChpassOk = "chpassOk"
	// (Reserved) Page value to set connection.
	PageConnection = "_ClientConnection"
	// (Reserved) Page value for authentication.
	PageAuthentication = "_ClientAuthentication"
	// (Reserved) Page value for change-password.
	PageChpass = "_ClientChpass"
)

var b64 = base64.StdEncoding

var Home string

func b64write(file string, tx string) {
	cgiio.WriteAll(file, cryp.Cryp(kGlobal, tx))
}

func b64read(file string) string {
	return cryp.Decryp(kGlobal, cgiio.ReadAll(file))
}

func readUsers() map[string][]string {
  fusers := path.Join(Home, ssDb)
  var r map[string][]string
	if err := json.Unmarshal([]byte(b64read(fusers)), &r); err != nil {
		Err(err.Error())
	}
  return r
}

func readSessions() map[string][]interface{} {
  fsessions := path.Join(Home, ssDb)
  var r map[string][]interface{}
	if err := json.Unmarshal([]byte(b64read(fsessions)), &r); err != nil {
		Err(err.Error())
	}
  return r
}

// Init initializes libcgi
func Init(home string) {
	Home = home

	fusers := path.Join(home, uDb)
	fsessions := path.Join(home, ssDb)

	if !cgiio.Exists(fusers) {
		userdb := make(map[string]interface{})
		userdb["admin"] = []string{
			"CaZkw7OHNkp618+7zWhasQcHYc/BBWNV+zeyqVPQjwU982S4/d1PwvSWtVPFE4upqI" +
				"kuvFYlQ9IZnCdI80vN8iid54Xn/3Cwki/SDVYNWFvKXsTIs8Z0Z/v+",
			"0",
		}
		bs, err := json.Marshal(userdb)
		if err != nil {
			Err(err.Error())
		}
		b64write(fusers, string(bs))
	}
	if !cgiio.Exists(fsessions) {
		sessiondb := make(map[string]interface{})
		bs, err := json.Marshal(sessiondb)
		if err != nil {
			Err(err.Error())
		}
		b64write(fsessions, string(bs))
	}
}

// Connect sets pageId and returns true if sessionId is correct.
func Connect(sessionId, pageId string)bool {
  sessionsOld := readSessions()
  now := time.Now().Unix()
  sessions := make(map[string][]interface{})
  for k, v := range sessionsOld {
    if v[ssTime].(int64) > now {
      sessions[k] = v
    }
  }

  data, ok := sessions[sessionId]
  if !ok {
    return false
  }


/*
    var now = Date.now().getTime();
    var sessions = new Map<String, Dynamic>();
    It.fromIterator(sessionsOld.keys()).each(function (k) {
      var v = sessionsOld.get(k);
      if (v[SS_TIME] > now) {
        sessions.set(k, v);
      }
    });

    var ssData = sessions.get(sessionId);
    if (ssData == null) {
      return false;
    }
    ssData[SS_PAGE_ID] = pageId;
    ssData[SS_TIME] = ssData[SS_INC_TIME] + now;

    b64write(fsessions, Json.from(sessions));
    */
    return true;

}

// ReadRequest reads a request sent in os.Args[1]. This request is a map which
// have been first 'jsonized' and then codified in B64.
func ReadRequest() map[string]interface{} {
	rq, err := b64.DecodeString(os.Args[1])
	if err != nil {
		Err(err.Error())
	}

	var r map[string]interface{}
	if err := json.Unmarshal(rq, &r); err != nil {
		Err(err.Error())
	}

	return r
}

// Err sends an error message to client. The message is an object JSON type:
//    {error:msg}
func Err(msg string) {
	r := make(map[string]interface{})
	r[Error] = msg
	rp, err := json.Marshal(r)
	if err == nil {
		fmt.Print(b64.EncodeToString(rp))
	} else {
		fmt.Print("{\"${Error}\" : \"Error in Err\"}")
	}
  panic("");
}

// Ok sends a response ('rp') to client. It adds a field rp["error"] = "".
func Ok(rp map[string]interface{}) {
	rp[Error] = ""
	st, err := json.Marshal(rp)
	if err == nil {
		fmt.Print(b64.EncodeToString(st))
	} else {
		fmt.Print("{\"${error}\" : \"Error in Ok\"}")
	}
}
