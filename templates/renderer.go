package templates

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
)

type Renderer struct {
	Templates *template.Template
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.LoadTemplates()
	return r
}

const Raw = "raw"
const Home = "home"
const Network = "network"
const Peers = "peers"
const ListMap = "listMap"
const TxpoolStatus = "txpoolStatus"
const BlockNumber = "blockNumber"

// Taken out of the constructor with the idae of forced template reloading
func (r *Renderer) LoadTemplates() {
	var allFiles []string
	files, err := os.ReadDir("./templates")
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".htemplate") {
			allFiles = append(allFiles, "./templates/"+filename)
		}
	}
	r.Templates, err = template.ParseFiles(allFiles...) //parses all .tmpl files in the 'Templates' folder
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Templates loaded")
}

func (r *Renderer) RenderResponse(w io.Writer, data *RenderData) error {
	err := r.Templates.ExecuteTemplate(w, data.TemplateName, data)
	if err != nil {
		log.Println(err)
	}
	return err

}

// This is a try to bring some uniformity to passing data to the Templates
// The "RenderData" container is a wrapper for the header/body/footer containers
type RenderData struct {
	Error        error
	TemplateName string
	HeaderData   HeaderData
	BodyData     interface{}
	FooterData   interface{}
}

type HeaderData interface {
	GetRefresh() (interval int)
	SetRefresh(int)
}
