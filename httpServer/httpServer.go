package main

import (
	//"os"
	"strings"
	"fmt"
	"httpJson/jsonParser"
	"encoding/json"
	"net/http"
	//"encoding/base64"
	"httpServer/httpSession"
	"strconv"
//	"golang.org/x/net/http2"
)

var jsonD jsonParser.JsonData
func parseCompleteJsonFile() {
	jsonD = jsonD.ParseJsonFile("../httpJson/data.json")
	//fmt.Printf("\n [%s] json", json.MarshalJsonData())
}

func fullDataFunc(res http.ResponseWriter, req *http.Request) {
	parseCompleteJsonFile()
	res.Write([]byte(jsonD.MarshalJsonData()))
}

func SetCooke(res http.ResponseWriter, ua string, ip string, name string) {
	ok,value:=httpSession.CreateSessionId(ua,ip, true)
	if ok {
	   cook := http.Cookie{Name:name, Value:value, Path: "/posts", Domain: "www.Saravanan.com", Secure:true, HttpOnly:true}
	   http.SetCookie(res, &cook)
        }
}

func GetUserAgent(req *http.Request) string {
	userDetail:=req.Header.Get("User-Agent")
	return userDetail
}

func userDataFunc(res http.ResponseWriter, req *http.Request) {
	usrAgt := GetUserAgent(req)
	ip:=httpSession.GetIPAddress(req)
	cook, err := req.Cookie("user-data-cook")
	if err == nil {
		sess:=cook.Value
                httpSession.CheckSessionId(sess)
	} else {
	    SetCooke(res, usrAgt, ip, "user-data-cook")
        }

	if jsonD.Pcount == 0 {
	   parseCompleteJsonFile()
   	}
	if req.Method == http.MethodGet  {
		fmt.Println("Query:", req.URL.Query())
		param := req.URL.Query()
		fmt.Printf("Request Type Mail[%s] \n", req.URL.String())
		lURL := strings.Split(req.URL.Path, "/")
		fmt.Printf("%d-%s-%s", len(lURL), lURL[0], lURL[1])
		if len(lURL) > 3 {
			res.Write([]byte("Invalid URL"))
			return
		}
		idx, _ :=strconv.Atoi(lURL[2])
		fmt.Printf("\n %d ", idx)
		if idx == 0 {
			if len(param) == 0 {
				res.Write([]byte(jsonD.MarshalUserDataAll()))
			} else {
				keys, ok := req.URL.Query()["id"]
				var secKeys []string
				var usrD jsonParser.UserDataSlice
				//var cnt int
				if !ok || len(keys) == 0 {
					secKeys, ok = param["userId"]
					if !ok || len(secKeys) == 0 {
						res.Write([]byte("Invalid URL query parameter"))
						return
					}
				}
				if len(keys) == 0 {
					var err error
					if idx, err = strconv.Atoi(secKeys[0]);  err == nil {
						usrD, _ = jsonD.ParseUserData(idx,1)
					}
				} else {
					var err error
					if idx, err  = strconv.Atoi(keys[0]); err == nil {
						usrD, _ = jsonD.ParseUserData(idx,0)
					}
				}
				fmt.Printf("\nQuery processing %d %v", idx , usrD)
				res.Write([]byte(usrD.MarshalUserData()))
			}
			return
		}
		if idx > jsonD.Ccount {
			res.Write([]byte("Invalid URL idx"))
			return
		}
		fmt.Printf("Request Type User \n")
		usrD, _ := jsonD.ParseUserData(idx, 0)
		res.Write([]byte(usrD[0].MarshalUserData()))
		return
	}
	if req.Method == http.MethodPut {
	}
	if req.Method == http.MethodPost {
		var usrD jsonParser.UserData
		if err:=req.ParseForm(); err == nil {

			for key, val := range req.Form {
				switch key {
				case "id","Id":
					usrD.Id,_ = strconv.Atoi(val[0])
				case "userId", "UserId":
					usrD.UserId,_ = strconv.Atoi(val[0])
				case "title", "Title":
					usrD.Title = val[0]
				case "body", "Body":
					usrD.Body=val[0]
				case "height", "Height":
					usrD.Height,_ = strconv.ParseFloat(val[0], 64)
				}
			}
		}
		fmt.Printf("Body:%v", req.Body)
		if err := json.NewDecoder(req.Body).Decode(&usrD); err != nil {
			fmt.Printf("\n Error on decode json ", err)
			res.Write([]byte("Invalid Json"))
			return
		}
		if err := usrD.CheckAndInsertUserData(jsonD); err != nil {
			fmt.Printf("\n Error on adding user data ", err)
			res.Write([]byte("Invaid key"))
			return
		}
		res.Write([]byte("Post successfull !!"))
		return
	}
	if req.Method == http.MethodPatch {
	}
}

func mailDataFunc(res http.ResponseWriter, req *http.Request) {
	//http.SetCookie(res, "Mail-Cookie")
	if jsonD.Ccount == 0 {
	   parseCompleteJsonFile()
   	}
	fmt.Printf("Request Type Mail[%s] \n", req.URL.String())
	lURL := strings.Split(req.URL.Path, "/")
	fmt.Printf("%d-%s-%s", len(lURL), lURL[0], lURL[1])
	if len(lURL) > 3 {
		res.Write([]byte("Invalid URL"))
		return
	}
	idx, _ :=strconv.Atoi(lURL[2]) 
	fmt.Printf("\n %d ", idx)
	if idx == 0 {
		if jsonD.Ccount == 0 {
			parseCompleteJsonFile()
		}
		res.Write([]byte(jsonD.MarshalMailDataAll()))
		return
	}
	if idx > jsonD.Ccount {
		res.Write([]byte("Invalid URL"))
		return
	}
	
	milD := jsonD.ParseMailData(idx)
	res.Write([]byte(milD.MarshalMailData()))
}

func albumDataFunc(res http.ResponseWriter, req *http.Request) {
	if jsonD.Acount == 0 {
	   parseCompleteJsonFile()
   	}
	fmt.Printf("Request Type Album [%s]\n", req.URL.Path)
	res.Write([]byte(jsonD.MarshalAlbumDataAll()))
}

func mailMiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	lURL := strings.Split(req.URL.Path, "/")
	if len(lURL) > 3 {
		res.Write([]byte("Invalid URL"))
		return
	}
	if len(lURL) == 2 {
		if jsonD.Ccount == 0 {
			parseCompleteJsonFile()
		}
		res.Write([]byte(jsonD.MarshalMailDataAll()))
		return
	}
	next.ServeHTTP(res, req)
	})
}

func main() {

	//usrHandler := http.HandlerFunc(userDataFunc)
	mux := http.NewServeMux()
	mailMux := http.NewServeMux()

	mailMux.HandleFunc("/", mailDataFunc)
	mux.HandleFunc("/posts/", userDataFunc)
	mux.HandleFunc("/mails/",mailDataFunc)
	mux.HandleFunc("/albums", albumDataFunc)

	s := &http.Server{
	  Addr:    ":8081",
	  Handler: mux,
	}
	if s.ListenAndServe() != nil {
		return
	}
}
