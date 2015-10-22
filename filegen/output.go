package filegen

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/maciekmm/sitemap-generator/config"
)

var header = xml.ProcInst{
	Target: "xml", Inst: []byte(`version="1.0" encoding="UTF-8"`),
}

var startElement = xml.StartElement{
	Name: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "urlset"},
	Attr: []xml.Attr{
		xml.Attr{Name: xml.Name{Space: "", Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
		xml.Attr{Name: xml.Name{Space: "", Local: "xsi:schemaLocation"}, Value: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd/sitemap.xsd"},
	},
}

type output struct {
	config.Filter
	encoder *xml.Encoder
	urls    int
	fileNo  int
}

func (out *output) put(entry string) error {
	var modT *time.Time
	if out.Filter.IncludeModificationDate {
		mod := time.Now()
		modT = &mod
	}
	err := out.encoder.Encode(&url{entry, modT, out.Filter.Modifiers})
	if err != nil {
		return fmt.Errorf("Error occured while encoding: %s", err)
	}
	if out.Filter.PerFile != 0 {
		out.urls++
		if out.urls > out.Filter.PerFile {
			err := out.nextFile()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (out *output) nextFile() error {
	if out.encoder != nil {
		out.clean()
	}
	out.fileNo++
	out.urls = 0
	file, err := os.OpenFile("./output/"+out.Filter.FilePrefix+"-"+strconv.Itoa(out.fileNo)+".xml", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	enc := xml.NewEncoder(file)
	enc.EncodeToken(header)
	enc.Indent("", "\n")
	enc.EncodeToken(startElement)
	enc.Indent("", "  ")
	out.encoder = enc
	return nil
}

func (out *output) clean() {
	out.encoder.EncodeToken(startElement.End())
	out.encoder.Flush()
}
