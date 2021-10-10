package jsonParser

import (
	"encoding/json"
	"os"
	"fmt"
	"io/ioutil"
	"errors"
	)


type MailData struct {
	PostID int `json:"postId"`
	Id int     `json:"id"`
	Name string   `json:"name"`
	Email string  `json:"email"`
	Body string   `json:"body"`
}

type AlbumData struct {
	AlbumId int `json:"albumId"`
	Id int      `json:"id"`
	Title string   `json:"title"`
	Url string     `json:"url"`
	ThumbnailUrl string `json:"thumbnailurl"`
}

type JsonData struct {
	Posts []UserData    `json:"posts"`
	Comments []MailData `json:"comments"`
	Albums []AlbumData  `json:"albums"`
	Pcount int
	Ccount int
	Acount int
}

type UserData struct {
	UserId int `json:"userId"`
	Id int     `json:"id"`
	Title string   `json:"title"`
	Body string   `json:"body"`
	Height float64 `json:"height"`
}

type ParseFile interface {
	ParseJsonString()
	ParseJsonFile() JsonData
	MarshalJsonData()
	ParseUserData(id int) UserData
	ParseUserDataAll() []UserData
	ParseMailData(id int) MailData
	ParseMailDataAll() []MailData
	ParseAlbumData(id int) AlbumData
	ParseAlbumDataAll() []AlbumData
}

type UserDataSlice  []UserData

type UserMani interface {
	ParseUserString(string) UserData
	MarshalUserData() string
}

type MailMani interface {
	ParseMailString(string) MailData
	MarshalMailData() string
}

type AlbumMani interface {
	ParseAlbmString(string) AlbumData
	MarshalAlbmData() string
}

func (gcomData JsonData) ParseJsonString(jsonStr string) {
	json.Unmarshal([]byte(jsonStr), &gcomData)
}


func (gcomData JsonData) ParseJsonFile(fileName string) JsonData {
	fp, err := os.Open(fileName)
	if (err != nil) {
		fmt.Println("\nError in opening the file [%s] error is[%s]", fileName, err)
		return gcomData
	}
	defer fp.Close()

	byteVal,_ := ioutil.ReadAll(fp)
	jsonD := JsonData{}
	err = json.Unmarshal(byteVal, &jsonD)
	err = json.Unmarshal(byteVal, &gcomData)
	if err != nil {
		fmt.Printf("\n Error %s", err)
	}
	//fmt.Printf("\n%v", gcomData)
	jsonD.Pcount = len(jsonD.Posts)
	jsonD.Ccount = len(jsonD.Comments)
	jsonD.Acount = len(jsonD.Albums)
	return jsonD
}

func (gcomData JsonData) MarshalJsonData() string {
	b, err := json.Marshal(gcomData)
	if ( err != nil) {
		fmt.Println("\nError in forming json string [%s]", err)
		return ""
	} else
	{
	    return string(b)
        }
}

//User Data
//keyIdx == 0 primary key unique ex id
//keyIdx == 1 seconday key unique to few recs ex userId
func (gusers JsonData) ParseUserData(id int, keyIdx int) ([]UserData, int)  {
	idx := -1
	var usrD []UserData
	for _, elem := range gusers.Posts {
		if  (elem.Id == id && keyIdx == 0) {
			idx++
			usrD = append(usrD, elem)
			break
		} else if elem.UserId == id {
			idx++
			usrD = append(usrD, elem)
		}
	}
	return usrD, idx
}

func (gusers JsonData) ParseUserDataAll() []UserData {
	return gusers.Posts
}

func (gusers JsonData) MarshalUserDataAll() string {
	b, err := json.Marshal(gusers.Posts)
	if ( err != nil) {
		fmt.Println("\nError in forming json string [%s]", err)
		return ""
	} else
	{
	    return string(b)
        }
}

func (gusers UserData) CheckAndInsertUserData(jsData JsonData) error  {
	indx := -1
	fmt.Printf("\n user data:%d", gusers.Id)
	for idx, elem := range jsData.Posts {
		indx=idx
		if gusers.Id == elem.Id {
			return errors.New("Key exists")
		}
		
	}
	jsData.Posts[indx] = gusers
	return nil
}

func (gusers UserData) MarshalUserData() string {
	b, err := json.Marshal(gusers)
	if ( err != nil) {
		fmt.Println("\nError in forming json string [%s]", err)
		return ""
	} else
	{
	    return string(b)
        }
}

func (gusers UserDataSlice) MarshalUserData() string {
	b, err := json.Marshal(gusers)
	if ( err != nil) {
		fmt.Println("\nError in forming json string [%s]", err)
		return ""
	} else
	{
	    return string(b)
        }
}

func (guser UserData) ParseUserString(jsonStr string) UserData{
	var usrD UserData
	json.Unmarshal([]byte(jsonStr), &usrD)
	return usrD
}

//Mail Data
func (gcomData JsonData) ParseMailDataAll( ) []MailData{
	return gcomData.Comments
}

func (gcomData JsonData) ParseMailData(id int) MailData{
	var milD MailData
	for _, elem := range gcomData.Comments {
		if  (elem.Id == id) {
			milD = elem
			break
		}
	}
	return milD
}

func (gmail MailData) ParseMailString(jsonStr string) {
	json.Unmarshal([]byte(jsonStr), &gmail)
}

func (gmail MailData) MarshalMailData() string {
	b, err := json.Marshal(gmail)
	if ( err != nil) {
		fmt.Println("\nError in forming json string [%s]", err)
		return ""
	} else
	{
	    return string(b)
        }
}

func (gmail JsonData) MarshalMailDataAll() string {
	b, err := json.Marshal(gmail.Comments)
	if ( err != nil) {
		fmt.Println("\nError in forming json string [%s]", err)
		return ""
	} else
	{
	    return string(b)
        }
}

//Album Data
func (galbum AlbumData) ParseAlbumFile(fileName string, id string) {
	fp, err := os.Open(fileName)
	if (err != nil) {
		fmt.Println("\nError in opening the file [%s] error is[%s]", fileName, err)
		return
	}
	defer fp.Close()

	byteVal,_ := ioutil.ReadAll(fp)
	var comData JsonData

	json.Unmarshal(byteVal, &comData)
/*
	for _, elem := range comData.albums {
		if  (elem.id == id) {
			galbum = elem
			break
		}
	}
	*/
}

func (galbum AlbumData) ParseAlbumString(jsonStr string) {
	json.Unmarshal([]byte(jsonStr), &galbum)
}

func (galbum AlbumData) MarshalAlbumData() string {
	b, err := json.Marshal(galbum)
	if ( err != nil) {
		fmt.Println("\nError in forming json string [%s]", err)
		return ""
	}
	return string(b)
}

func (galbum JsonData) MarshalAlbumDataAll() string {
	b, err := json.Marshal(galbum.Albums)
	if ( err != nil) {
		fmt.Println("\nError in forming json string [%s]", err)
		return ""
	}
	return string(b)
}
