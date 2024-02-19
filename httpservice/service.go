package httpservice

import (
	"context"
	"net/http"
	"strings"

	"github.com/san-lab/ethwebtool/create"
	"github.com/san-lab/ethwebtool/merkledemo"
	"github.com/san-lab/ethwebtool/templates"
)

const passwdFile = "http.passwd.json"

//This is the glue between the http requests and the (hopefully) generic RPC client

type LilHttpHandler struct {
	//config     Config
	runContext context.Context
	renderer   *templates.Renderer
	refresh    int
}

// Creating a naw http handler with its embedded rpc client and html renderer
func NewHttpHandler() *LilHttpHandler {
	lhh := &LilHttpHandler{}

	lhh.renderer = templates.NewRenderer()

	lhh.refresh = 5
	return lhh
}

// Handles incoming requests. Some will be forwarded to the RPC client.
// Assumes the request path has either: 1 part - interpreted as a /command with logic implemented within the client
//
//	or: 2 parts - interpreted as /node/ethMethod
//
// The port No set at Client initialization is used for the RPC call
func (lhh *LilHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//isSlash := func(c rune) bool { return c == '/' }
	//f := strings.FieldsFunc(r.URL.Path, isSlash)
	//log.Println(f)
	pathTokens := strings.Split(r.URL.Path, "/")

	rdata := &templates.RenderData{HeaderData: nil, TemplateName: templates.Home}
	rdata.BodyData = ">>" + r.URL.Path + "<<"
	//	rdata.HeaderData.SetRefresh(lhh.refresh)
	r.ParseForm()

	if len(pathTokens) > 1 {
		switch pathTokens[1] {
		case "create":
			create.CallCreate(r, rdata)
		case "home":
			rdata.BodyData = "home"
		case "create2":
			create.CallCreate2(r, rdata)
		case "merkledemo":
			merkledemo.CallMerkleDemo(r, rdata)
		case "loadtemplates":
			lhh.renderer.LoadTemplates()
		}
		//TODO: dispatch
	}
	lhh.renderer.RenderResponse(w, rdata)

}

type handler func(w http.ResponseWriter, r *http.Request)

type Config struct {
	RPCFirstEntry string
	RPCTLS        bool
	MockMode      bool
	DumpRPC       bool
	StartWatchdog bool
	BasicAuth     bool
	DebugMode     bool
}
