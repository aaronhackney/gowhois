// whois.go
package main

import (
	"encoding/json"
	"fmt"
)

type ReturnJSON struct {
	WhoisRecord   *Whois         `json:"whoIsRecord"`
	ContactRecord *ContactRecord `json:"ContactRecord,omitempty"`
}

type Whois struct {
	StartAddress     string              `json:"startAddress"`
	EndAddress       string              `json:"endAddress"`
	Handle           string              `json:"handle"`
	Name             string              `json:"name"`
	RegistrationDate string              `json:"registrationDate"`
	UpdateDate       string              `json:"updateDate"`
	Version          string              `json:"version"`
	OriginASes       string              `json:"originASes,omitempty"`
	ParentRefUrl     map[string]string   `json:"parentRefUrl"`
	ContactRef       map[string]string   `json:"ContactRef"`
	Comments         []string            `json:"comments,omitempty"`
	NetBlocks        []map[string]string `json:"netBlocks"`
}

type ContactRecord struct {
	Handle        string   `json:"handle"`
	Name          string   `json:"name"`
	StreetAddress []string `json:"address"`
	City          string   `json:"city"`
	State         string   `json:"state"`
	PostalCode    string   `json:"postalCode"`
	Country       string   `json:"country"`
	ContactType   string   `json:"type"`
	Reference     string   `json:"reference"`
}

func (*Whois) getContactRecord(url string) (*ContactRecord, error) {
	content, err := getContent(fmt.Sprintf(url))

	if err != nil {
		return nil, err
	}

	var contactRecord ContactRecord
	var tempJsonMap map[string]interface{}

	// unmarshall into a map of interfaces
	if err := json.Unmarshal(content, &tempJsonMap); err != nil {
		return nil, err
	}

	var prefix interface{}
	if org, exists := tempJsonMap["org"]; exists {
		contactRecord.ContactType = "org"
		prefix = org
	} else if cust, exists := tempJsonMap["customer"]; exists {
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
	contactRecord.Reference = prefix.(map[string]interface{})["ref"].(map[string]interface{})["$"].(string)

	return &contactRecord, nil

}

func (*Whois) unmarshalResponse(b []byte) (*Whois, error) {
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
		comments, _ := convertToSlice(rawComment.(map[string]interface{})["line"])
		var returnComments []string
		for line := range comments {
			returnComments = append(returnComments, comments[line].(map[string]interface{})["$"].(string))
		}
		whois.Comments = returnComments
	}

	// NetBlocks
	netBlockList, err := convertToSlice(jsonUnwound["netBlocks"].(map[string]interface{})["netBlock"])
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	for i := range netBlockList {
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
	whois.Handle = jsonUnwound["handle"].(map[string]interface{})["$"].(string)
	whois.Name = jsonUnwound["name"].(map[string]interface{})["$"].(string)
	whois.RegistrationDate = jsonUnwound["registrationDate"].(map[string]interface{})["$"].(string)
	whois.UpdateDate = jsonUnwound["updateDate"].(map[string]interface{})["$"].(string)
	whois.Version = jsonUnwound["version"].(map[string]interface{})["$"].(string)
	whois.ParentRefUrl = parentRefUrl
	whois.NetBlocks = returnNetBlocks

	return &whois, nil
}

func (*Whois) generateJson(whoisRecord *Whois, contactRecord *ContactRecord) ([]byte, error) {
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
