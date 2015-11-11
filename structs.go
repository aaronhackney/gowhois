package main

import (
	"encoding/json"
)

type Handle struct {
	Handle string `json:"$"`
}

type Name struct {
	Name string `json:"$"`
}

type StreetAddress struct {
	Line struct {
		StreetAddress string `json:"$"`
	}
}

type City struct {
	City string `json:"$"`
}

type State struct {
	State string `json:"$"`
} // `json:"iso3166-2"`

type PostalCode struct {
	PostalCode string `json:"$"`
}

type Country struct {
	Name struct {
		Name string `json:"$"`
	}
	Code2 struct {
		Code2 string `json:"$"`
	}
} // `json:"iso3166-1"`

type CustomerRecord struct {
	Customer struct {
		Handle        Handle
		Name          Name
		StreetAddress StreetAddress
		City          City
		State         State `json:"iso3166-2"`
		PostalCode    PostalCode
		Country       Country `json:"iso3166-1"`
	}
}

type OrgRecord struct {
	Org struct {
		Handle        Handle
		Name          Name
		StreetAddress StreetAddress
		City          City
		State         State `json:"iso3166-2"`
		PostalCode    PostalCode
		Country       Country `json:"iso3166-1"`
	}
}

type WhoisRecord struct {
	Net struct {
		NetworkRef struct {
			NetworkRef string `json:"$"`
		} `json:"ref"`

		ParentNetworkRef struct {
			Reference string `json:"$"`
			Handle    string `json:"@handle"`
			Name      string `json:"@name"`
		} `json:"parentNetRef"`

		EndAddress struct {
			EndAddress string `json:"$"`
		}

		StartAddress struct {
			StartAddress string `json:"$"`
		}

		OwnerInfo struct {
			Name      string `json:"@name"`
			Handle    string `json:"@handle"`
			Reference string `json:"$"`
		} `json:"customerRef"`

		OrgRef struct {
			Name      string `json:"@name"`
			Handle    string `json:"@handle"`
			Reference string `json:"$"`
		} `json:"orgRef"`

		Version struct {
			Version string `json:"$"`
		}
		UpdateDate struct {
			UpdateDate string `json:"$"`
		}

		Name struct {
			Name string `json:"$"`
		}

		Handle struct {
			Handle string `json:"$"`
		}

		NetBlocks struct {
			NetblockRaw   json.RawMessage `json:"netblock"`
			NetBlock      Netblock        `json:"-"`
			NetblockArray []Netblock      `json:"-"`
		}

		Comment struct {
			LineRaw   json.RawMessage `json:"line"`
			Line      CommentLine     `json:"-"`
			LineArray []CommentLine   `json:"-"`
		} `json:"comment"`
	}
}

type CommentLine struct {
	Line string `json:"$"`
}

type Netblock struct {
	CidrLength struct {
		CidrLength string `json:"$"`
	}

	Description struct {
		Description string `json:"$"`
	}

	EndAddress struct {
		EndAddress string `json:"$"`
	}

	StartAddress struct {
		StartAddress string `json:"$"`
	}

	BlockType struct {
		Type string `json:"$"`
	} `json:"type"`
}

type ReturnJSON struct {
	WhoisRecord    *WhoisRecord    `json:"whoIsRecord"`
	CustomerRecord *CustomerRecord `json:"customerRecord,omitempty"`
	OrgRecord      *OrgRecord      `json:"orgRecord,omitempty"`
}
