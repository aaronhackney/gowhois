package main

type ReturnJSON struct {
	WhoisRecord   *Whois         `json:"whoIsRecord"`
	ContactRecord *ContactRecord `json:"ContactRecord,omitempty"`
}

type Whois struct {
	StartAddress     string              `json:"startAddress"`
	EndAddress       string              `json:"endAddress"`
	handle           string              `json:"handle"`
	name             string              `json:"name"`
	RegistrationDate string              `json:"registrationDate"`
	UpdateDate       string              `json:"updateDate"`
	version          string              `json:"version"`
	WhoisRefUrl      string              `json:"whoisRefUrl"`
	OriginASes       string              `json:"originASes"`
	ParentRefUrl     map[string]string   `json:"parentRefUrl"`
	OrgRef           map[string]string   `json:"OrgRef"`
	CustomerRef      map[string]string   `json:"CustomerRef"`
	ContactRef       map[string]string   `json:"ContactRef"`
	Comments         []string            `json:"comments"`
	netBlocks        []map[string]string `json:"netBlocks"`
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
	reference     string   `json:"reference"`
}
