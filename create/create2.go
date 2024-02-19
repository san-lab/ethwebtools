package create

import (
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/proveniencenft/kmsclitool/common"
	"github.com/san-lab/ethwebtool/templates"
)

func CallCreate2(r *http.Request, rdata *templates.RenderData) {
	rdata.TemplateName = "create2"
	crp := &Create2Payload{}
	rdata.BodyData = crp
	adrs := r.FormValue("address")
	chashs := r.FormValue("codehash")
	salts := r.FormValue("salt")
	action := r.FormValue("action")
	if action == "Calculate" {
		adrs = strings.TrimPrefix(adrs, "0x")
		adr, err := hex.DecodeString(adrs)
		if err != nil {
			rdata.Error = err
			return
		}
		adr20 := make([]byte, 20)
		copy(adr20[20-len(adr):], adr)
		crp.Address = common.CRCAddressString(adr20)

		chashs = strings.TrimPrefix(chashs, "0x")
		chash, err := hex.DecodeString(chashs)
		if err != nil {
			rdata.Error = err
			return
		}
		if len(chash) < 32 {
			x := make([]byte, 32)
			copy(x[32-len(chash):], chash)
			chash = x
		}
		crp.Codehash = "0x" + hex.EncodeToString(chash)
		salts = strings.TrimPrefix(salts, "0x")
		salt, err := hex.DecodeString(salts)
		if err != nil {
			rdata.Error = err
			return
		}
		if len(salt) < 32 {
			x := make([]byte, 32)
			copy(x[32-len(salt):], salt)
			salt = x
		}
		crp.Salt = "0x" + hex.EncodeToString(salt)

		cadr, err := common.CalcCREATE2Address(adr20, chash, salt, 0)
		if err != nil {
			crp.CAddress = err.Error()
		} else {
			crp.CAddress = common.CRCAddressString(cadr)
		}
	}
}

type Create2Payload struct {
	Address  string `json:"address"`
	Codehash string `json:"codehash"`
	Salt     string `json:"salt"`
	CAddress string `json:"caddress"`
}
