// Code generated by "bitfanDoc -codec encoder"; DO NOT EDIT
package linecodec

import "github.com/awillis/bitfan/processors/doc"

func Doc() *doc.Codec {
	return &doc.Codec{
		Name:       "encoder",
		PkgName:    "linecodec",
		ImportPath: "github.com/awillis/bitfan/codecs/line",
		Doc:        "doc codec",
		DocShort:   "",
		Decoder: &doc.Decoder{
			Doc: "doc decoder",
			Options: &doc.CodecOptions{
				Doc: "doc decoderOptions",
				Options: []*doc.CodecOption{
					&doc.CodecOption{
						Name:           "Delimiter",
						Alias:          "",
						Doc:            "Change the delimiter that separates lines",
						Required:       false,
						Type:           "string",
						DefaultValue:   "\"\\\\n\"",
						PossibleValues: []string{},
						ExampleLS:      "",
					},
				},
			},
		},
		Encoder: &doc.Encoder{
			Doc: "doc encoder",
			Options: &doc.CodecOptions{
				Doc: "doc encoderOptions",
				Options: []*doc.CodecOption{
					&doc.CodecOption{
						Name:           "Delimiter",
						Alias:          "",
						Doc:            "Change the delimiter that separates lines",
						Required:       false,
						Type:           "string",
						DefaultValue:   "\"\\\\n\"",
						PossibleValues: []string{},
						ExampleLS:      "",
					},
					&doc.CodecOption{
						Name:           "Format",
						Alias:          "format",
						Doc:            "Format as a golang text/template",
						Required:       false,
						Type:           "string",
						DefaultValue:   "\"{{TS \"dd/MM/yyyy:HH:mm:ss\" .}} {{.host}} {{.message}}\"",
						PossibleValues: []string{},
						ExampleLS:      "",
					},
					&doc.CodecOption{
						Name:           "Var",
						Alias:          "var",
						Doc:            "You can set variable to be used in Statements by using ${var}.\neach reference will be replaced by the value of the variable found in Statement's content\nThe replacement is case-sensitive.",
						Required:       false,
						Type:           "hash",
						DefaultValue:   nil,
						PossibleValues: []string{},
						ExampleLS:      "var => {\"hostname\"=>\"myhost\",\"varname\"=>\"varvalue\"}",
					},
				},
			},
		},
	}
}
