package main

import (
	"flag"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"log"
	"os"
)

func main() {
	var err error
	input := flag.String("i", "in.pdf", "Input")
	output := flag.String("o", "out.pdf", "Output")

	flag.Parse()

	var rs *os.File
	if rs, err = os.Open(*input); err != nil {
		log.Fatal(err)
	}

	var ws *os.File
	if ws, err = os.Create(*output); err != nil {
		log.Fatal(err)
	}

	config := pdfcpu.NewDefaultConfiguration()
	ctx, err := api.ReadContext(rs, config)

	if err = api.WriteContext(ctx, ws); err != nil {
		log.Fatal(err)
	}
}
