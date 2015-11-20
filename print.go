// print
package main

import (
	"fmt"
)

func printRecord(whoisRecord *WhoisRecord, customerRecord *CustomerRecord, orgRecord *OrgRecord) error {
	fmt.Println("Gowhois 0.2 https://github.com/aaronhackney/gowhois")
	/////////////////////////
	// NETBLOCKS
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("Net Range:\t" + whoisRecord.Net.StartAddress.StartAddress + " - " + whoisRecord.Net.EndAddress.EndAddress)

	if len(whoisRecord.Net.NetBlocks.NetBlock) > 0 {
		for i := 0; i < len(whoisRecord.Net.NetBlocks.NetBlock); i++ {
			fmt.Println("CIDR:\t\t" + whoisRecord.Net.NetBlocks.NetBlock[i].StartAddress.StartAddress + "/" + whoisRecord.Net.NetBlocks.NetBlock[i].CidrLength.CidrLength)
		}
	}

	fmt.Println("-----------------------------------------------------------------")

	/////////////////////////
	// ORG Data
	if string(whoisRecord.Net.OrgRef.Reference) != "" {
		fmt.Println("\nOrg Handle: " + orgRecord.Org.Handle.Handle + " (" + whoisRecord.Net.OrgRef.Reference + ")")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("\t " + orgRecord.Org.Name.Name)

		if len(orgRecord.Org.StreetAddress.Line) > 0 {
			for i := 0; i < len(orgRecord.Org.StreetAddress.Line); i++ {
				fmt.Println("\t " + orgRecord.Org.StreetAddress.Line[i].Line)
			}
		}
		fmt.Println("\t " + orgRecord.Org.City.City + " " + " " + orgRecord.Org.State.State + " " + orgRecord.Org.PostalCode.PostalCode + " " + orgRecord.Org.Country.Code2.Code2)
		fmt.Println("\n")
	}

	/////////////////////////
	// Customer Data
	if string(whoisRecord.Net.OwnerInfo.Reference) != "" {
		fmt.Println("\nOwner Handle: " + customerRecord.Customer.Handle.Handle + " (" + whoisRecord.Net.OwnerInfo.Reference + ")")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("\t " + customerRecord.Customer.Name.Name)

		if len(customerRecord.Customer.StreetAddress.Line) > 0 {
			for i := 0; i < len(customerRecord.Customer.StreetAddress.Line); i++ {
				fmt.Println("\t " + customerRecord.Customer.StreetAddress.Line[i].Line)
			}
		}

		fmt.Println("\t " + customerRecord.Customer.City.City + " " + " " + customerRecord.Customer.State.State + " " + customerRecord.Customer.PostalCode.PostalCode + " " + customerRecord.Customer.Country.Code2.Code2)
		fmt.Println("\n")

	}

	return nil
}