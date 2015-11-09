package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"io/ioutil"
	"net/http"
	"regexp"
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

func printRecord(whoisRecord *WhoisRecord,customerRecord *CustomerRecord, orgRecord *OrgRecord)(error) {
	
	fmt.Println("\nNetblock: " + whoisRecord.Net.NetBlocks.Netblock.StartAddress.StartAddress + "/" + whoisRecord.Net.NetBlocks.Netblock.CidrLength.CidrLength + "\t (" + whoisRecord.Net.NetworkRef.NetworkRef + ")")
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("\t " + whoisRecord.Net.NetBlocks.Netblock.StartAddress.StartAddress + " - " + whoisRecord.Net.NetBlocks.Netblock.EndAddress.EndAddress)
	
	if string(whoisRecord.Net.OrgRef.Reference) != "" {
		fmt.Println("\nOrg Handle: " + orgRecord.Org.Handle.Handle + "\t("  + whoisRecord.Net.OrgRef.Reference + ")")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("\t " + orgRecord.Org.Name.Name)		
		fmt.Println("\t " + orgRecord.Org.StreetAddress.Line.StreetAddress + "\t " )
		fmt.Println("\t " + orgRecord.Org.City.City + " " + " " + orgRecord.Org.State.State + " " + orgRecord.Org.PostalCode.PostalCode + " " + orgRecord.Org.Country.Code2.Code2)
		fmt.Println("\n")
	}

	if string(whoisRecord.Net.OwnerInfo.Reference) != "" {
		fmt.Println("\nOwner Handle: " + customerRecord.Customer.Handle.Handle + "\t("  + whoisRecord.Net.OwnerInfo.Reference + ")")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("\t " + customerRecord.Customer.Name.Name)		
		fmt.Println("\t " + customerRecord.Customer.StreetAddress.Line.StreetAddress + "\t " )
		fmt.Println("\t " + customerRecord.Customer.City.City + " " + " " + customerRecord.Customer.State.State + " " + customerRecord.Customer.PostalCode.PostalCode + " " + customerRecord.Customer.Country.Code2.Code2)
		fmt.Println("\n")
	}

	return nil
}

func generateJson(whoisRecord *WhoisRecord,customerRecord *CustomerRecord, orgRecord *OrgRecord)([]byte, error) {
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
func main() {
	var customerRecord *CustomerRecord
	var orgRecord *OrgRecord
	var validIP = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	
	// Flags
    isJson := flag.Bool("json", false, "json: change output from a screen print to JSON formatted output")	
	flag.Parse()

	//fmt.Println("IP Arguments:", flag.Args())
	
	ip := flag.Args()[0]
	
	if  !validIP.MatchString(ip) {
		fmt.Println("You must input a valid IP address: gowhois 1.2.3.4")
		os.Exit(3)
	}
	
	url := "http://whois.arin.net/rest/ip/" + ip
	
	whois, _ := getWhois(url)

	if string(whois.Net.OwnerInfo.Reference) != "" {
		customerRecord, _ = getCustRecord(string(whois.Net.OwnerInfo.Reference))
	}

	if string(whois.Net.OrgRef.Reference) != "" {
		orgRecord, _ = getOrgRecord(string(whois.Net.OrgRef.Reference))
	}

	if *isJson {
		jsonOutput, _ := generateJson(whois, customerRecord, orgRecord)
		fmt.Println(string(jsonOutput))
	} else {
		printRecord(whois, customerRecord, orgRecord)
	}
}
