package main

import (
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

func convertToSlice(object interface{}) ([]interface{}, error) {
	switch v := object.(type) {
	case []interface{}:
		return object.([]interface{}), nil
	case interface{}:
		var returnInterfaceArray []interface{} = make([]interface{}, 1)
		returnInterfaceArray[0] = object
		return returnInterfaceArray, nil
	default:
		fmt.Println(v)
		return nil, nil
	}
}

func getAddressLines(rawJson interface{}) ([]string, error) {
	var streetAddress []string

	address, err := convertToSlice(rawJson.(map[string]interface{})["streetAddress"].(map[string]interface{})["line"])
	if err != nil {
		return nil, err
	}

	for line := range address {
		//fmt.Printf("\nADDRESS LINE: %+v\n", address[line])
		streetAddress = append(streetAddress, address[line].(map[string]interface{})["$"].(string))
	}
	return streetAddress, nil
}

/////////////////////////////////////

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
	//	var customerRecord *CustomerRecord
	//	var orgRecord *OrgRecord
	var whois *Whois
	var validIP = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)

	// Flags
	isJson := flag.Bool("json", false, "json: change output from a screen print to JSON formatted output: gowhois -json 1.2.3.4")
	flag.Parse()

	//fmt.Println("IP Arguments:", flag.Args())

	if len(flag.Args()) < 1 {
		help()
		os.Exit(0)
	}

	ip := flag.Args()[0]

	if !validIP.MatchString(ip) {
		help()
		os.Exit(3)
	}

	url := "http://whois.arin.net/rest/ip/" + ip

	content, _ := getContent(fmt.Sprintf(url))

	whois, _ = whois.unmarshalResponse(content)
	//fmt.Printf("\nwhois record %+v\n", whois)

	//fmt.Printf("\n DEBUG: %+v\n", whois.ContactRef)
	contactRecord, _ := whois.getContactRecord(whois.ContactRef["url"])

	// Output generation
	if *isJson {
		jsonOutput, _ := whois.generateJson(whois, contactRecord)
		fmt.Println(string(jsonOutput))
	} /*else {
		printRecord(whois, contactRecord)
	}*/
}
