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

	if whoisRecord.Net.NetBlocks.NetBlock.StartAddress.StartAddress != "" {
		fmt.Println("NetBlock:\t" + whoisRecord.Net.NetBlocks.NetBlock.StartAddress.StartAddress + "/" + whoisRecord.Net.NetBlocks.NetBlock.CidrLength.CidrLength)
	} else if len(whoisRecord.Net.NetBlocks.NetBlockArray) > 0 {
		for i := 0; i < len(whoisRecord.Net.NetBlocks.NetBlockArray); i++ {
			fmt.Println("CIDR:\t\t" + whoisRecord.Net.NetBlocks.NetBlockArray[i].StartAddress.StartAddress + "/" + whoisRecord.Net.NetBlocks.NetBlockArray[i].CidrLength.CidrLength)
		}
	}
	fmt.Println("-----------------------------------------------------------------")

	/////////////////////////
	// ORG Data
	if string(whoisRecord.Net.OrgRef.Reference) != "" {
		fmt.Println("\nOrg Handle: " + orgRecord.Org.Handle.Handle + " (" + whoisRecord.Net.OrgRef.Reference + ")")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("\t " + orgRecord.Org.Name.Name)

		// print the org address but we need to check for a string vs an array
		if orgRecord.Org.StreetAddress.Line.Line != "" {
			fmt.Println("\t " + orgRecord.Org.StreetAddress.Line.Line + "\t ")
		} else if len(orgRecord.Org.StreetAddress.LineArray) > 0 {
			for i := 0; i < len(orgRecord.Org.StreetAddress.LineArray); i++ {
				fmt.Println("\t " + orgRecord.Org.StreetAddress.LineArray[i].Line)
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

		// print the org address but we need to check for a string vs an array
		if customerRecord.Customer.StreetAddress.Line.Line != "" {
			fmt.Println("\t " + customerRecord.Customer.StreetAddress.Line.Line + "\t ")
		} else if len(customerRecord.Customer.StreetAddress.LineArray) > 0 {
			for i := 0; i < len(customerRecord.Customer.StreetAddress.LineArray); i++ {
				fmt.Println("\t " + customerRecord.Customer.StreetAddress.LineArray[i].Line)
			}
		}

		fmt.Println("\t " + customerRecord.Customer.City.City + " " + " " + customerRecord.Customer.State.State + " " + customerRecord.Customer.PostalCode.PostalCode + " " + customerRecord.Customer.Country.Code2.Code2)
		fmt.Println("\n")

		/*fmt.Println("\nOwner Handle: " + customerRecord.Customer.Handle.Handle + " (" + whoisRecord.Net.OwnerInfo.Reference + ")")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("\t " + customerRecord.Customer.Name.Name)
		//fmt.Println("\t " + customerRecord.Customer.StreetAddress.Line.StreetAddress + "\t ")
		fmt.Println("\t " + customerRecord.Customer.City.City + " " + " " + customerRecord.Customer.State.State + " " + customerRecord.Customer.PostalCode.PostalCode + " " + customerRecord.Customer.Country.Code2.Code2)
		fmt.Println("\n")*/
	}

	return nil
}
