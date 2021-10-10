package httpSession

import (
	"crypto/rand"
	"time"
	"strings"
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
        "strconv"
	"net"
	"net/http"
	"io"
	)

type Session struct {
	sessionId string
	userAgent UserAgent
	key []byte
	ip string
	expire int64
}

type UserAgent struct {
	browser, bVersion string
	osName, platform string
	mobile bool
	tablet bool
}


var httpSess = make(map[string]*Session)

func CheckSessionId(sess string) (*Session, bool) {
	if v, ok := httpSess[sess]; ok {
		return v, ok
	}
	return nil,false
}

func CreateSessionId(ua string, ip string, createFlag bool) (bool, string) {
	if httpSess == nil && createFlag == false {
		return false,""
	}

	if createFlag == true {
	   sessionID:=GenerateSessionId(ua, ip)
	   if session, ok := AddSessionId(sessionID); ok {
		fmt.Printf("\n Added new session")
		return true, session.sessionId
	   }
	}
	return false,""
}

func ConvStr2Ascii(val string)(string) {
     var str string
     for _ , ch := range val{
	 lstr := strconv.Itoa(int(ch))
	 fmt.Print(string(ch))
         str = str + lstr
     }
     return str
}

//encrypt encrypts plain string with a secret key and returns encrypt string.
func encrypt(plainData string, secret []byte) (string, error) {
	cipherBlock, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(plainData), nil)), nil
}

func GenerateSessionId(ua string, ip string) *Session {
     netSession := new(Session)
     sessionId, _ := rand.Prime(rand.Reader, 10)
     netSession.ip = ip
     var e error

     value := sessionId.String()+netSession.userAgent.browser+netSession.userAgent.osName+netSession.userAgent.platform+ip /* sessionid+ua+ip */
     netSession.key =  make([]byte, 16)
     if _, err := rand.Read(netSession.key); err != nil {
       return nil
     }
     
     netSession.sessionId, e = encrypt(value, netSession.key)
     if e != nil {
	     return nil
     }
     
     return netSession
}

func AddSessionId(sess *Session) (*Session,bool) {
	sess.expire = time.Now().Add(time.Minute * 30).Unix()
	httpSess[sess.sessionId] = sess
	return sess,true
}

func DeleteSession(sess string) bool {
	if _,ok:=CheckSessionId(sess);ok {
		delete(httpSess, sess)
		return true
	}
	return false
}

func DeleteSessionExpiry(mins int) int {
	curTime := time.Now().Unix()
	cnt :=0
	for k, e := range httpSess {
		if (curTime - e.expire) >= int64(mins) {
			delete(httpSess, k)
			cnt = cnt + 1
		}
	}
	return cnt
}

func DeleteSessionAll() bool {
	if httpSess == nil {
		return false
	}
	for k := range httpSess {
		delete(httpSess, k)
	}
	httpSess = nil
	return true
}

func PurgeSessiion() {
	for true {
		time.Sleep(30 * 60 * time.Second)
		DeleteSessionExpiry(30)
	}
}

func GetIPAddress(req *http.Request) string {
     _, port, err := net.SplitHostPort(req.RemoteAddr)
     if err != nil {
        fmt.Printf("userip: %q is not IP:port", req.RemoteAddr)
        return ""
     }
     for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
        addresses := strings.Split(req.Header.Get(h), ",")
        // march from right to left until we get a public address
        // that will be the address right before our proxy.
        for i := len(addresses) -1 ; i >= 0; i-- {
            ip := strings.TrimSpace(addresses[i])
            // header can contain spaces too, strip those out.
            realIP := net.ParseIP(ip)
            if !realIP.IsGlobalUnicast() {
                // bad address, go to next
                continue
            }
            return ip+port
        }
     }
     return ""+port
}

func ParseUserAgent(ua string) UserAgent {
	var lUA UserAgent
	if strings.Contains(strings.ToLower(ua), "windows") {
	   lUA.osName = "windows"
        } else if strings.Contains(strings.ToLower(ua), "linux") {
	   lUA.osName = "linux"
        } else if strings.Contains(strings.ToLower(ua), "aix") {
	   lUA.osName = "aix"
        } else if strings.Contains(strings.ToLower(ua), "android") {
           lUA.osName = "android"
	   lUA.mobile = true
        } else if strings.Contains(strings.ToLower(ua), "iphone") {
           lUA.osName = "iphone"
	   lUA.mobile = true
        } else if strings.Contains(strings.ToLower(ua), "ipad") {
           lUA.osName = "ipad"
	   lUA.mobile = false
	   lUA.tablet = true
        } else if strings.Contains(strings.ToLower(ua), "mac") {
           lUA.osName = "mac"
	   lUA.mobile = false
	   lUA.tablet = false
        }

	if strings.Contains(strings.ToLower(ua),"msie") {
		substring:=string(ua[strings.Index(ua,"MSIE"):])
		substring=string(substring[:strings.Index(substring,";")])
                lUA.browser=substring
        } else if  strings.Contains(strings.ToLower(ua),"opr") || strings.Contains(strings.ToLower(ua),"opera") {
            if strings.Contains(ua,"opera") {
	       substring:=string(ua[strings.Index(ua,"Opr"):])
               lUA.browser=strings.Replace(substring, "/", "-", -1)
	    } else if strings.Contains(ua,"opr") {
	       substring:=string(ua[strings.Index(ua,"Opr"):])
               lUA.browser=strings.Replace(substring, "/", "-", -1)
	    }
        } else if strings.Contains(strings.ToLower(ua),"chrome") {
	   substring:=string(ua[strings.Index(ua,"Chrome"):])
           substring=string(substring[:strings.Index(substring, " ")])
	   lUA.browser=strings.Replace(substring, "/", "-", -1)
        } else if strings.Contains(strings.ToLower(ua), "firefox") {
	   substring:=string(ua[strings.Index(ua,"Firefox"):])
           lUA.browser=strings.Replace(substring, "/", "-", -1)
        } else if strings.Contains(strings.ToLower(ua),"edg") {
	   substring:=string(ua[strings.Index(ua,"Edg"):])
           lUA.browser=strings.Replace(substring, "/", "-", -1)
        } else if strings.Contains(strings.ToLower(ua),"safari") && strings.Contains(strings.ToLower(ua),"version") {
	   substring:=string(ua[strings.Index(ua,"Safari"):])
           lUA.browser=strings.Replace(substring, "/", "-", -1)
	} else {
            lUA.browser = "UnKnown"+ua
        }
	return lUA
}
