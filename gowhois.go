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

func generateJson(whoisRecord *Whois, contactRecord *ContactRecord) ([]byte, error) {
	var returnJson ReturnJSON
	returnJson.WhoisRecord = whoisRecord
	returnJson.ContactRecord = contactRecord

	jsonOutput, err := json.MarshalIndent(&returnJson, "", "\t")

	if err != nil {
		fmt.Println("err:", err.Error())
		return nil, err
	}

	return jsonOutput, nil
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

func unmarshalJSON(b []byte) (*Whois, error) {
	var whois Whois
	var tempJsonMap map[string]interface{}
	var contactPrefix interface{}
	var returnNetBlocks []map[string]string

	jsonUnwound := make(map[string]interface{})

	// unmarshall into a map of interfaces
	if err := json.Unmarshal(b, &tempJsonMap); err != nil {
		return nil, err
	}

	// Extract the top level json nest []net
	for key, value := range tempJsonMap["net"].(map[string]interface{}) {
		jsonUnwound[key] = value
	}

	parentRefUrl := map[string]string{
		"url":    jsonUnwound["parentNetRef"].(map[string]interface{})["$"].(string),
		"handle": jsonUnwound["parentNetRef"].(map[string]interface{})["@handle"].(string),
		"name":   jsonUnwound["parentNetRef"].(map[string]interface{})["@name"].(string),
	}

	if prefix, exists := jsonUnwound["orgRef"]; exists {
		contactPrefix = prefix
	} else if prefix, exists := jsonUnwound["customerRef"]; exists {
		contactPrefix = prefix
	}

	whois.ContactRef = map[string]string{
		"url":    contactPrefix.(map[string]interface{})["$"].(string),
		"handle": contactPrefix.(map[string]interface{})["@handle"].(string),
		"name":   contactPrefix.(map[string]interface{})["@name"].(string),
	}

	// Comments
	if rawComment, exists := jsonUnwound["comment"]; exists {
		//comments, _ := convertToSlice(jsonUnwound["comment"].(map[string]interface{})["line"])
		comments, _ := convertToSlice(rawComment.(map[string]interface{})["line"])
		var returnComments []string
		for i := 0; i < len(comments); i++ {
			returnComments = append(returnComments, comments[i].(map[string]interface{})["$"].(string))
		}
		whois.Comments = returnComments
	}

	// NetBlocks
	netBlockList, err := convertToSlice(jsonUnwound["netBlocks"].(map[string]interface{})["netBlock"])
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	for i := 0; i < len(netBlockList); i++ {
		description := netBlockList[i].(map[string]interface{})["description"].(map[string]interface{})["$"].(string)
		endAddress := netBlockList[i].(map[string]interface{})["endAddress"].(map[string]interface{})["$"].(string)
		startAddress := netBlockList[i].(map[string]interface{})["startAddress"].(map[string]interface{})["$"].(string)
		blockType := netBlockList[i].(map[string]interface{})["type"].(map[string]interface{})["$"].(string)
		cidrLength := netBlockList[i].(map[string]interface{})["cidrLength"].(map[string]interface{})["$"].(string)
		netBlockObject := map[string]string{
			"description":  description,
			"startAddress": startAddress,
			"endAddress":   endAddress,
			"cidrLength":   cidrLength,
			"type":         blockType,
		}
		returnNetBlocks = append(returnNetBlocks, netBlockObject)
	}

	if originAS, exists := jsonUnwound["originASes"]; exists {
		whois.OriginASes = originAS.(map[string]interface{})["originAS"].(map[string]interface{})["$"].(string)
	}

	whois.StartAddress = jsonUnwound["startAddress"].(map[string]interface{})["$"].(string)
	whois.EndAddress = jsonUnwound["endAddress"].(map[string]interface{})["$"].(string)
	whois.handle = jsonUnwound["handle"].(map[string]interface{})["$"].(string)
	whois.name = jsonUnwound["name"].(map[string]interface{})["$"].(string)
	whois.RegistrationDate = jsonUnwound["registrationDate"].(map[string]interface{})["$"].(string)
	whois.UpdateDate = jsonUnwound["updateDate"].(map[string]interface{})["$"].(string)
	whois.version = jsonUnwound["version"].(map[string]interface{})["$"].(string)
	whois.ParentRefUrl = parentRefUrl
	whois.netBlocks = returnNetBlocks

	return &whois, nil
}

func getContactRecord(url string) (*ContactRecord, error) {
	content, err := getContent(fmt.Sprintf(url))

	if err != nil {
		return nil, err
	}
	//b := []byte(entity)
	//b := []byte(customer)

	var contactRecord ContactRecord
	var tempJsonMap map[string]interface{}

	// unmarshall into a map of interfaces
	if err := json.Unmarshal(content, &tempJsonMap); err != nil {
		return nil, err
	}

	var prefix interface{}
	if org, exists := tempJsonMap["org"]; exists {
		//fmt.Printf("We have a org record type: %+v:\n", org)
		contactRecord.ContactType = "org"
		prefix = org
	} else if cust, exists := tempJsonMap["customer"]; exists {
		//fmt.Printf("We have a customer record type: %+v:\n", cust)
		contactRecord.ContactType = "customer"
		prefix = cust
	}

	contactRecord.Handle = prefix.(map[string]interface{})["handle"].(map[string]interface{})["$"].(string)
	contactRecord.Name = prefix.(map[string]interface{})["name"].(map[string]interface{})["$"].(string)
	contactRecord.City = prefix.(map[string]interface{})["city"].(map[string]interface{})["$"].(string)
	contactRecord.State = prefix.(map[string]interface{})["iso3166-2"].(map[string]interface{})["$"].(string)
	contactRecord.PostalCode = prefix.(map[string]interface{})["postalCode"].(map[string]interface{})["$"].(string)
	contactRecord.Country = prefix.(map[string]interface{})["iso3166-1"].(map[string]interface{})["code2"].(map[string]interface{})["$"].(string)
	contactRecord.StreetAddress, _ = getAddressLines(prefix)
	contactRecord.reference = prefix.(map[string]interface{})["ref"].(map[string]interface{})["$"].(string)

	//fmt.Printf("Contact Record: %+v\n", contactRecord)

	return &contactRecord, nil

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

	whois, _ := unmarshalJSON(content)
	//fmt.Printf("\nwhois record %+v\n", whois)

	//fmt.Printf("\n DEBUG: %+v\n", whois.ContactRef)
	contactRecord, _ := getContactRecord(whois.ContactRef["url"])

	// Output generation
	if *isJson {
		jsonOutput, _ := generateJson(whois, contactRecord)
		fmt.Println(string(jsonOutput))
	} /*else {
		printRecord(whois, contactRecord)
	}*/
}
