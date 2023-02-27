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

	if *input == *output {
		log.Fatalf("Input file and output file are the same: %s == %s", *input, *output)
	}

	var rs *os.File
	if rs, err = os.Open(*input); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = rs.Close()
	}()

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
					color := sd.Dict.NameEntry("ColorSpace")
					if color == nil || *color != "DeviceGray" {
						continue
					}
					buf, err := FilterImage(sd.Raw)
					if err != nil {
						log.Fatal(err)
					}
					*ent.Offset = 0
					sd.Raw = buf
					*sd.StreamLength = int64(len(buf))
					sd.StreamLengthObjNr = nil
					sd.StreamOffset = 0
					ent.Object = sd
				}
			}
		}
	}
	err = api.ValidateContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var ws *os.File
	if ws, err = os.Create(*output); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = ws.Close()
	}()
	if err = api.WriteContext(ctx, ws); err != nil {
		log.Fatal(err)
	}
}
