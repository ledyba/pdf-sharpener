package main

import (
	"bytes"
	"flag"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"image"
	"image/jpeg"
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
	defer rs.Close()

	var ws *os.File
	if ws, err = os.Create(*output); err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	config := pdfcpu.NewDefaultConfiguration()
	ctx, err := api.ReadContext(rs, config)
	if err != nil {
		log.Fatal(err)
	}

	err = api.ValidateContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range ctx.Table {
		if d, ok := e.Object.(pdfcpu.Dict); ok {
			if typ := d.Type(); typ == nil || *typ != "Page" {
				continue
			}
			res := d.DictEntry("Resources")
			if res == nil {
				continue
			}
			xObj := res.DictEntry("XObject")
			if xObj == nil {
				continue
			}
			for _, o := range xObj {
				if ref, ok := o.(pdfcpu.IndirectRef); ok {
					ent, found := ctx.XRefTable.FindTableEntry(ref.ObjectNumber.Value(), ref.GenerationNumber.Value())
					if !found {
						continue
					}
					sd, ok := ent.Object.(pdfcpu.StreamDict)
					if !ok {
						continue
					}
					if !sd.Image() {
						continue
					}
					buf := bytes.Buffer{}
					img := image.NewRGBA(image.Rect(0, 0, 100, 100))
					err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 98})
					if err != nil {
						log.Fatal(err)
					}
					sd.Raw = buf.Bytes()
					*sd.StreamLength = int64(buf.Len())
					ent.Object = sd
				}
			}
		}
	}

	if err = api.WriteContext(ctx, ws); err != nil {
		log.Fatal(err)
	}
}
