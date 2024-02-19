package create

import (
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"github.com/proveniencenft/kmsclitool/common"
	"github.com/san-lab/ethwebtool/templates"
)

var addresses = []string{"0xBc0E100905580439968D19175c2F1b344173D8B2", "0xBC0E189507D624De860E1e5b3B06DEE3475A70e9"}
var privkeys = []string{"0xe45bf4d4c15999c8a8795dde55ea8c77b73726a97ce18e23d7ab5b126c3481e7",
	"0xc522c068090d4e888dadbab9967fd81a79a451aff84dce2040df59ad5a6ce1e8"}

func CallCreate(r *http.Request, rdata *templates.RenderData) {
	rdata.TemplateName = "create"
	crp := &CreatePayload{}
	crp.Nonces = map[string]string{}
	for i, adr := range addresses {
		nonces := GetNonce(adr)
		crp.Nonces[adr] = nonces + " private key: " + privkeys[i]
	}
	rdata.BodyData = crp
	adrs := r.FormValue("address")
	nons := r.FormValue("nonce")
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
		non, err := strconv.Atoi(nons)
		if err != nil {
			rdata.Error = err
			return
		}
		crp.Nonce = strconv.Itoa(non)
		cadr, err := common.CalcCREATEAddress(adr20, uint(non))
		if err != nil {
			crp.CAddress = err.Error()
		} else {
			crp.CAddress = common.CRCAddressString(cadr)
		}
	}
	rdata.BodyData = crp
}

type CreatePayload struct {
	Address  string `json:"address"`
	Nonce    string `json:"nonce"`
	CAddress string `json:"caddress"`
	Nonces   map[string]string
}
