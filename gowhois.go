package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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

	// Peekahead to see if a [] or a {}
	if string(customerRecord.Customer.StreetAddress.LineRaw[0]) == "[" {
		err = json.Unmarshal(customerRecord.Customer.StreetAddress.LineRaw, &customerRecord.Customer.StreetAddress.LineArray)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			return nil, err
		}
	} else {
		err = json.Unmarshal(customerRecord.Customer.StreetAddress.LineRaw, &customerRecord.Customer.StreetAddress.Line)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			return nil, err
		}
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

	// Peekahead to see if a [] or a {}
	if string(orgRecord.Org.StreetAddress.LineRaw[0]) == "[" {
		err = json.Unmarshal(orgRecord.Org.StreetAddress.LineRaw, &orgRecord.Org.StreetAddress.LineArray)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			return nil, err
		}
	} else {
		err = json.Unmarshal(orgRecord.Org.StreetAddress.LineRaw, &orgRecord.Org.StreetAddress.Line)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			return nil, err
		}
	}

	return &orgRecord, err
}

func generateJson(whoisRecord *WhoisRecord, customerRecord *CustomerRecord, orgRecord *OrgRecord) ([]byte, error) {
	var returnJson ReturnJSON
	returnJson.WhoisRecord = whoisRecord
	returnJson.CustomerRecord = customerRecord
	returnJson.OrgRecord = orgRecord

	jsonOutput, err := json.MarshalIndent(&returnJson, "", "\t")

	if err != nil {
		fmt.Println("err:", err.Error())
		return nil, err
	}

	return jsonOutput, nil

}

func unmarshalWhoisJson(content []byte) (*WhoisRecord, error) {
	var whois WhoisRecord

	err := json.Unmarshal(content, &whois)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}

	// Peekahead to see if a [] or a {}
	if len(whois.Net.NetBlocks.NetBlockRaw) > 0 {
		if string(whois.Net.NetBlocks.NetBlockRaw[0]) == "[" {
			err = json.Unmarshal(whois.Net.NetBlocks.NetBlockRaw, &whois.Net.NetBlocks.NetBlockArray)
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				return nil, err
			}
		} else {
			makeArray(whois.Net.NetBlocks.NetBlockRaw)
			err = json.Unmarshal(whois.Net.NetBlocks.NetBlockRaw, &whois.Net.NetBlocks.NetBlock)
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				return nil, err
			}
		}
	}

	// Peekahead to see if a [] or a {}

	if len(whois.Net.Comment.LineRaw) > 0 {
		if string(whois.Net.Comment.LineRaw[0]) == "[" {
			err = json.Unmarshal(whois.Net.Comment.LineRaw, &whois.Net.Comment.LineArray)
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				return nil, err
			}
		} else {
			err = json.Unmarshal(whois.Net.Comment.LineRaw, &whois.Net.Comment.Line)
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				return nil, err
			}
		}
	}

	return &whois, nil
}

func makeArray(rawJson json.RawMessage) json.RawMessage {
	// take raw message and make it an array
	var myObject interface{}
	var mySlice []interface{}
	err := json.Unmarshal(rawJson, &myObject)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return nil
	}

	fmt.Printf("Type: %T", myObject)

	mySlice = append(mySlice, myObject)

	// now reencode to rawjson

	return nil

}

func help() {
	fmt.Println("------------------------------------------------------------\n")
	fmt.Println("You must input a valid IP address.\n")
	fmt.Println("Examples:\n\t\t gowhois 1.2.3.4")
	fmt.Println("\t\t gowhois -json 1.2.3.4\n")
	fmt.Println("gowhois --help for info on switches\n")
	fmt.Println("------------------------------------------------------------\n\n")
	return
}

func main() {
	var customerRecord *CustomerRecord
	var orgRecord *OrgRecord
	var validIP = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)

	// Flags
	isJson := flag.Bool("json", false, "json: change output from a screen print to JSON formatted output: gowhois -json 1.2.3.4")
	flag.Parse()

	//fmt.Println("IP Arguments:", flag.Args())

	if len(flag.Args()) < 1 {
		help()
		os.Exit(3)
	}

	ip := flag.Args()[0]

	if !validIP.MatchString(ip) {
		help()
		os.Exit(3)
	}

	url := "http://whois.arin.net/rest/ip/" + ip

	//whois, _ := getWhois(url)
	content, _ := getContent(fmt.Sprintf(url))

	// Unmarshall the raw server response
	whois, _ := unmarshalWhoisJson(content)

	fmt.Println()

	// Move these into the unmarshal function
	if string(whois.Net.OwnerInfo.Reference) != "" {
		customerRecord, _ = getCustRecord(string(whois.Net.OwnerInfo.Reference))
	}

	if string(whois.Net.OrgRef.Reference) != "" {
		orgRecord, _ = getOrgRecord(string(whois.Net.OrgRef.Reference))
	}

	// Output generation
	if *isJson {
		jsonOutput, _ := generateJson(whois, customerRecord, orgRecord)
		fmt.Println(string(jsonOutput))
	} else {
		printRecord(whois, customerRecord, orgRecord)
	}
}
