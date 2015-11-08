package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getContent(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "Golang WHOIS Aaron 1.0")
	req.Header.Set("Accept", "application/json")

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func getWhois(url string) (*WhoisRecord, error) {
	content, err := getContent(fmt.Sprintf(url))

	var whoisRecord WhoisRecord

	err = json.Unmarshal(content, &whoisRecord)
	if err != nil {
		// An error occurred while converting our JSON to an object
		fmt.Println("func getWhois(url string) (*WhoisRecord, error)")
		fmt.Println(err)
		return nil, err
	}

	return &whoisRecord, err
}

func getCustRecord(url string) (*CustomerRecord, error) {
	content, err := getContent(fmt.Sprintf(url))

	var customerRecord CustomerRecord

	err = json.Unmarshal(content, &customerRecord)

	if err != nil {
		// An error occurred while converting our JSON to an object
		fmt.Println("func getCustRecord(url string) (*CustomerRecord, error)")
		fmt.Println(err)
		return nil, err
	}

	return &customerRecord, err
}

func getOrgRecord(url string) (*OrgRecord, error) {
	content, err := getContent(fmt.Sprintf(url))

	var orgRecord OrgRecord

	err = json.Unmarshal(content, &orgRecord)

	if err != nil {
		// An error occurred while converting our JSON to an object
		fmt.Println("func getOrgRecord(url string) (*OrgRecord, error)")
		fmt.Println(err)
		return nil, err
	}

	return &orgRecord, err
}

func main() {

	//TODO: CLI input validation
	var customerRecord *CustomerRecord
	var orgRecord *OrgRecord

	argsWithoutProg := os.Args[1:]
	url := "http://whois.arin.net/rest/ip/" + argsWithoutProg[0]
	whois, _ := getWhois(url)

	if string(whois.Net.OwnerInfo.Reference) != "" {
		customerRecord, _ = getCustRecord(string(whois.Net.OwnerInfo.Reference))
	}

	if string(whois.Net.OrgRef.Reference) != "" {
		orgRecord, _ = getOrgRecord(string(whois.Net.OrgRef.Reference))
	}

	var returnJson ReturnJSON
	returnJson.WhoisRecord = whois
	returnJson.CustomerRecord = customerRecord
	returnJson.OrgRecord = orgRecord

	jsonOutput, err := json.MarshalIndent(&returnJson, "", "\t")

	if err != nil {
		fmt.Println("err:", err.Error())
	}

	fmt.Printf("%v", string(jsonOutput))
}
