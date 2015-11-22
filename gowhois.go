package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	//"regexp"
)

var version = "0"

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

func help() {
	fmt.Println("------------------------------------------------------------\n")
	fmt.Println("You must input a valid IP address.\n")
	fmt.Println("Examples:\n\t\t gowhois 208.0.0.0")
	fmt.Println("\t\t gowhois -json 1.2.3.4\n")
	fmt.Println("gowhois --help for info on switches\n")
	fmt.Println("------------------------------------------------------------\n\n")
	return
}

func main() {
	var whois *Whois
	//var validIP = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)

	// Flags
	isJson := flag.Bool("json", false, "json: change output from a screen print to JSON formatted output: gowhois -json 1.2.3.4")
	isVersion := flag.Bool("v", false, "v: Prints the gowhois version: gowhois -v")
	flag.Parse()

	if *isVersion {
		fmt.Println("\nGowhois version", version)
		fmt.Println("https://github.com/aaronhackney/gowhois\n")
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		help()
		os.Exit(0)
	}

	ip := flag.Args()[0]

	if net.ParseIP(ip) == nil {
		help()
		os.Exit(3)
	}

	/*if !validIP.MatchString(ip) {
		help()
		os.Exit(3)
	}*/

	url := "http://whois.arin.net/rest/ip/" + ip

	content, _ := getContent(fmt.Sprintf(url))

	whois, _ = whois.unmarshalResponse(content)
	contactRecord, _ := whois.getContactRecord(whois.ContactRef["url"])

	if *isJson {
		jsonOutput, _ := whois.generateJson(whois, contactRecord)
		fmt.Println(string(jsonOutput))
	} else {
		printRecord(whois, contactRecord)
	}
}
