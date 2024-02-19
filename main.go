package main

import (
	"github.com/san-lab/ethwebtool/httpservice"

	"github.com/san-lab/commongo/gohttpservice"
)

func main() {
	h := httpservice.NewHttpHandler()
	gohttpservice.DefPort = "8100"
	gohttpservice.Startserver(h)
}
