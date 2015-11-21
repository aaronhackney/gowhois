package main

import (
	"fmt"
)

func printRecord(whoisRecord *Whois, contactRecord *ContactRecord) error {
	fmt.Println("\nGowhois", version, "https://github.com/aaronhackney/gowhois")
	/////////////////////////
	// NETBLOCKS
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("Net Range:\t" + whoisRecord.StartAddress + " - " + whoisRecord.EndAddress)

	if len(whoisRecord.NetBlocks) > 0 {
		for i := range whoisRecord.NetBlocks {
			fmt.Println("CIDR:\t\t" + whoisRecord.NetBlocks[i]["startAddress"] + "/" + whoisRecord.NetBlocks[i]["cidrLength"] + "\t(" + whoisRecord.NetBlocks[i]["description"] + ")")
		}
	}

	fmt.Println()
	for i := range whoisRecord.Comments {
		fmt.Println(whoisRecord.Comments[i])
	}

	fmt.Println("-----------------------------------------------------------------")

	/////////////////////////
	// Contact Data
	// change to exists
	fmt.Println("Contact Handle: " + contactRecord.Handle + " (" + contactRecord.Reference + ")")
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("\t " + contactRecord.Name)

	if len(contactRecord.StreetAddress) > 0 {
		for i := range contactRecord.StreetAddress {
			fmt.Println("\t " + contactRecord.StreetAddress[i])
		}
	}

	fmt.Println("\t "+contactRecord.City, contactRecord.State, contactRecord.PostalCode, contactRecord.Country)
	fmt.Println("\n")

	return nil
}
