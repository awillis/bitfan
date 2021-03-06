// Code generated by "bitfanDoc "; DO NOT EDIT
package websocket

import "github.com/awillis/bitfan/processors/doc"

func (p *processor) Doc() *doc.Processor {
	return &doc.Processor{
		Name:       "websocket",
		ImportPath: "github.com/awillis/bitfan/processors/websocket",
		Doc:        "Send event received over a ws connection to connected clients",
		DocShort:   "Reads events from standard input",
		Options: &doc.ProcessorOptions{
			Doc: "",
			Options: []*doc.ProcessorOption{
				&doc.ProcessorOption{
					Name:           "Codec",
					Alias:          "",
					Doc:            "The codec used for outputed data.",
					Required:       false,
					Type:           "codec",
					DefaultValue:   "\"json\"",
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Uri",
					Alias:          "",
					Doc:            "URI path",
					Required:       false,
					Type:           "string",
					DefaultValue:   "\"wsout\"",
					PossibleValues: []string{},
					ExampleLS:      "",
				},
			},
		},
		Ports: []*doc.ProcessorPort{},
	}
}
