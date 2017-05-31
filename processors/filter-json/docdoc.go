// Code generated by "bitfanDoc "; DO NOT EDIT
package json

import "github.com/vjeantet/bitfan/processors/doc"

func (p *processor) Doc() *doc.Processor {
	return &doc.Processor{
  Name:       "json",
  ImportPath: "/Users/sodadi/go/src/github.com/vjeantet/bitfan/processors/filter-json",
  Doc:        "",
  DocShort:   "Parses JSON events",
  Options:    &doc.ProcessorOptions{
    Doc:     "",
    Options: []*doc.ProcessorOption{
      &doc.ProcessorOption{
        Name:           "Add_field",
        Alias:          "",
        Doc:            "If this filter is successful, add any arbitrary fields to this event.\nField names can be dynamic and include parts of the event using the %{field}.",
        Required:       false,
        Type:           "hash",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Add_tag",
        Alias:          "",
        Doc:            "If this filter is successful, add arbitrary tags to the event.\nTags can be dynamic and include parts of the event using the %{field} syntax.",
        Required:       false,
        Type:           "array",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Remove_field",
        Alias:          "",
        Doc:            "If this filter is successful, remove arbitrary fields from this event.",
        Required:       false,
        Type:           "array",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Remove_tag",
        Alias:          "",
        Doc:            "If this filter is successful, remove arbitrary tags from the event.\nTags can be dynamic and include parts of the event using the %{field} syntax",
        Required:       false,
        Type:           "array",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Source",
        Alias:          "",
        Doc:            "The configuration for the JSON filter",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Target",
        Alias:          "",
        Doc:            "Define the target field for placing the parsed data. If this setting is omitted,\nthe JSON data will be stored at the root (top level) of the event",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
    },
  },
  Ports: []*doc.ProcessorPort{},
}
}