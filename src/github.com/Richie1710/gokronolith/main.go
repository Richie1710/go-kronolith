package gokronolith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const timeFormat = "20060102T150405Z"

type listresponse struct {
	Version string   `json:"version"`
	Result  []string `json:"result"`
	ID      int      `json:"id"`
}

type listrespons2 struct {
	Version string `json:"version"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
}

//Vcard Struct with Data from Vcard
type Vcard struct {
	DTSTART      time.Time
	DTEND        time.Time
	DTSTAMP      time.Time
	UID          string
	CREATED      time.Time
	LASTMODIFIED time.Time
	SUMMARY      string
}

//Getcurrententry Gets current entrys from horde
func Getcurrententry(url string, calender string, userid string, password string) ([]string, error) {
	tn := time.Now().Unix()
	te := time.Now().AddDate(0, 0, 1).Unix()
	fmt.Println(tn)
	fmt.Println(te)
	fmt.Printf("Checking from %s to %s. \n", time.Unix(tn, 0).UTC(), time.Unix(te, 0).UTC())
	client := &http.Client{}
	stringwithtime := fmt.Sprintf("{ \"jsonrpc\": \"1.0\", \"method\": \"calendar.listUids\", \"params\": [\"%s\", %d, %d ], \"id\": 1}", calender, tn, te)
	var jsonStr = []byte(stringwithtime)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(userid, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var listi listresponse
	err = json.Unmarshal(b, &listi)
	if err != nil {
		return nil, err
	}
	return listi.Result, nil
}

//GetEntryByTime Gets current entrys from horde
func GetEntryByTime(url string, calender string, userid string, password string, startunix int64, stopunix int64) ([]string, error) {
	fmt.Printf("Checking from %s to %s. \n", time.Unix(startunix, 0).UTC(), time.Unix(stopunix, 0).UTC())
	fmt.Printf("Checking from %d to %d. \n", startunix, stopunix)
	client := &http.Client{}
	stringwithtime := fmt.Sprintf("{ \"jsonrpc\": \"1.0\", \"method\": \"calendar.listUids\", \"params\": [\"%s\", %d, %d ], \"id\": 1}", calender, startunix, stopunix)
	var jsonStr = []byte(stringwithtime)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(userid, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var listi listresponse
	err = json.Unmarshal(b, &listi)
	if err != nil {
		return nil, err
	}
	return listi.Result, nil
}

//FilterEntryObjectsByTime : Filters Objects from Type VCard by Timestmaps (DTStart, DTStop)
func FilterEntryObjectsByTime(entryobjects []Vcard, startunix int64, stopunix int64) []Vcard {
	var newvcard []Vcard
	for _, vcard := range entryobjects {
		if vcard.DTEND.Unix() >= startunix && vcard.DTEND.Unix() <= stopunix {
			newvcard = append(newvcard, vcard)
		}
	}
	return newvcard

}

//GetICSByEntry Gets ICS card by entry of GetEntryByTime
func GetICSByEntry(url string, entry string, userid string, password string) (string, error) {
	client := &http.Client{}
	stringwithentry := fmt.Sprintf("{ \"jsonrpc\": \"1.0\", \"method\": \"calendar.export\", \"params\": [\"%s\", \"text/calendar\", [], null], \"id\": 1}", entry)
	var jsonStr = []byte(stringwithentry)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(userid, password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	var listi listrespons2
	err = json.Unmarshal(b, &listi)
	if err != nil {
		return "", err
	}
	return listi.Result, nil
}

//GetICSObjectByEntry Gets ICS card by entry of GetEntryByTime
func GetICSObjectByEntry(ics string) (Vcard, error) {

	var vcard Vcard
	var err error

	newline := strings.Split(ics, "\n")
	for _, s := range newline {
		if strings.Contains(s, "CREATED:") {
			string := strings.Split(s, ":")
			vcard.CREATED, err = time.Parse(timeFormat, strings.TrimSpace(string[1]))
			if err != nil {
				return vcard, err
			}
		}
		if strings.Contains(s, "DTEND:") {
			string := strings.Split(s, ":")
			vcard.DTEND, err = time.Parse(timeFormat, strings.TrimSpace(string[1]))
			if err != nil {
				return vcard, err
			}
		}
		if strings.Contains(s, "DTSTART:") {
			string := strings.Split(s, ":")
			vcard.DTSTART, err = time.Parse(timeFormat, strings.TrimSpace(string[1]))
			if err != nil {
				return vcard, err
			}
		}
		if strings.Contains(s, "DTSTAMP:") {
			string := strings.Split(s, ":")
			vcard.DTSTAMP, err = time.Parse(timeFormat, strings.TrimSpace(string[1]))
			if err != nil {
				return vcard, err
			}
		}
		if strings.Contains(s, "LAST-MODIFIED:") {
			string := strings.Split(s, ":")
			vcard.LASTMODIFIED, err = time.Parse(timeFormat, strings.TrimSpace(string[1]))
			if err != nil {
				return vcard, err
			}
		}
		if strings.Contains(s, "SUMMARY:") {
			string := strings.Split(s, ":")
			vcard.SUMMARY = string[1]
		}
		if strings.Contains(s, "UID:") {
			string := strings.Split(s, ":")
			vcard.UID = string[1]
		}
	}
	return vcard, err
}

//BrowseCalender Get a list of all available Calenders
func BrowseCalender(url string, userid string, password string, path string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf("{ \"jsonrpc\": \"1.0\", \"method\": \"calendar.browse\", \"params\": [\"%s\"], \"id\": 1}", path))))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(userid, password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	fmt.Println(string(b))
	var listi listrespons2
	err = json.Unmarshal(b, &listi)
	if err != nil {
		return "", err
	}
	return listi.Result, nil
}
