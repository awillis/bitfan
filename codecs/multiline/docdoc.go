// Code generated by "bitfanDoc -codec multiline"; DO NOT EDIT
package multilinecodec

import "github.com/awillis/bitfan/processors/doc"

func Doc() *doc.Codec {
	return &doc.Codec{
		Name:       "multiline",
		PkgName:    "multilinecodec",
		ImportPath: "github.com/awillis/bitfan/codecs/multiline",
		Doc:        "The multiline codec will collapse multiline messages and merge them into a single event.\n\nThe original goal of this codec was to allow joining of multiline messages from files into a single event. For example, joining Java exception and stacktrace messages into a single event.\n\nThe config looks like this:\n```\ninput {\n  stdin {\n    codec => multiline {\n      pattern => \"pattern, a regexp\"\n      negate => true or false\n      what => \"previous\" or \"next\"\n    }\n  }\n}\n```\nThe pattern should match what you believe to be an indicator that the field is part of a multi-line event.\n\nThe what must be previous or next and indicates the relation to the multi-line event.\n\nThe negate can be true or false (defaults to false). If true, a message not matching the pattern will constitute a match of the multiline filter and the what will be applied. (vice-versa is also true)\n\nFor example, Java stack traces are multiline and usually have the message starting at the far-left, with each subsequent line indented. Do this:\n\n```\ninput {\n  stdin {\n    codec => multiline {\n      pattern => \"^\\\\s\"\n      what => \"previous\"\n    }\n  }\n}\n```\nThis says that any line starting with whitespace belongs to the previous line.\n\nAnother example is to merge lines not starting with a date up to the previous line..\n\n```\ninput {\n  file {\n    path => \"/var/log/someapp.log\"\n    codec => multiline {\n      # Grok pattern names are valid! :)\n      pattern => \"^%{TIMESTAMP_ISO8601} \"\n      negate => true\n      what => \"previous\"\n    }\n  }\n}\n```\nThis says that any line not starting with a timestamp should be merged with the previous line.\n\nOne more common example is C line continuations (backslash). Here’s how to do that:\n\n```\nfilter {\n  multiline {\n    pattern => \"\\\\$\"\n    what => \"next\"\n  }\n}\n```\nThis says that any line ending with a backslash should be combined with the following line.",
		DocShort:   "",
		Decoder: &doc.Decoder{
			Doc: "Merges multiline messages into a single event",
			Options: &doc.CodecOptions{
				Doc: "",
				Options: []*doc.CodecOption{
					&doc.CodecOption{
						Name:           "Delimiter",
						Alias:          "delimiter",
						Doc:            "Change the delimiter that separates lines",
						Required:       false,
						Type:           "string",
						DefaultValue:   "\"\\n\"",
						PossibleValues: []string{},
						ExampleLS:      "",
					},
					&doc.CodecOption{
						Name:           "Negate",
						Alias:          "negate",
						Doc:            "Negate the regexp pattern (if not matched).",
						Required:       false,
						Type:           "bool",
						DefaultValue:   "false",
						PossibleValues: []string{},
						ExampleLS:      "",
					},
					&doc.CodecOption{
						Name:           "Pattern",
						Alias:          "pattern",
						Doc:            "The regular expression to match with the line",
						Required:       true,
						Type:           "string",
						DefaultValue:   nil,
						PossibleValues: []string{},
						ExampleLS:      "pattern => \"^\\\\s\"",
					},
					&doc.CodecOption{
						Name:         "What",
						Alias:        "what",
						Doc:          "If the pattern matched, does event belong to the next or previous event?",
						Required:     false,
						Type:         "string",
						DefaultValue: "\"previous\"",
						PossibleValues: []string{
							"previous",
							"next",
						},
						ExampleLS: "",
					},
				},
			},
		},
		Encoder: (*doc.Encoder)(nil),
	}
}
