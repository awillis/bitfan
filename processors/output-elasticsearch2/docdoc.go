// Code generated by "bitfanDoc "; DO NOT EDIT
package elasticsearch2

import "github.com/vjeantet/bitfan/processors/doc"

func (p *processor) Doc() *doc.Processor {
	return &doc.Processor{
  Name:       "elasticsearch2",
  ImportPath: "/Users/sodadi/go/src/github.com/vjeantet/bitfan/processors/output-elasticsearch2",
  Doc:        "",
  DocShort:   "",
  Options:    &doc.ProcessorOptions{
    Doc:     "",
    Options: []*doc.ProcessorOption{
      &doc.ProcessorOption{
        Name:           "DocumentType",
        Alias:          "document_type",
        Doc:            "The document type to write events to. There is no default value for this setting.\n\nGenerally you should try to write only similar events to the same type.\nString expansion %{foo} works here. Unless you set document_type, the event type will\nbe used if it exists otherwise the document type will be assigned the value of logs",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\"%{type}\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "FlushCount",
        Alias:          "flush_count",
        Doc:            "The number of requests that can be enqueued before flushing them. Default value is 1000",
        Required:       false,
        Type:           "int",
        DefaultValue:   "1000",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "FlushSize",
        Alias:          "flush_size",
        Doc:            "The number of bytes that the bulk requests can take up before the bulk processor decides to flush. Default value is 5242880 (5MB).",
        Required:       false,
        Type:           "int",
        DefaultValue:   "5242880",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Host",
        Alias:          "host",
        Doc:            "Host of the remote instance. Default value is \"localhost\"",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\"localhost\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "IdleFlushTime",
        Alias:          "idle_flush_time",
        Doc:            "The amount of seconds since last flush before a flush is forced. Default value is 1\n\nThis setting helps ensure slow event rates don’t get stuck.\nFor example, if your flush_size is 100, and you have received 10 events,\nand it has been more than idle_flush_time seconds since the last flush,\nthose 10 events will be flushed automatically.\nThis helps keep both fast and slow log streams moving along in near-real-time.",
        Required:       false,
        Type:           "int",
        DefaultValue:   "1",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Index",
        Alias:          "index",
        Doc:            "The index to write events to. Default value is \"logstash-%Y.%m.%d\"\n\nThis can be dynamic using the %{foo} syntax and strftime syntax (see http://strftime.org/).\nThe default value will partition your indices by day.",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\"logstash-%Y.%m.%d\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Password",
        Alias:          "password",
        Doc:            "Password to authenticate to a secure Elasticsearch cluster. There is no default value for this setting.",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Path",
        Alias:          "path",
        Doc:            "HTTP Path at which the Elasticsearch server lives. Default value is \"/\"\n\nUse this if you must run Elasticsearch behind a proxy that remaps the root path for the Elasticsearch HTTP API lives.",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\"/\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Port",
        Alias:          "port",
        Doc:            "ElasticSearch port to connect on. Default value is 9200",
        Required:       false,
        Type:           "int",
        DefaultValue:   "9200",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "User",
        Alias:          "user",
        Doc:            "Username to authenticate to a secure Elasticsearch cluster. There is no default value for this setting.",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "SSL",
        Alias:          "ssl",
        Doc:            "Enable SSL/TLS secured communication to Elasticsearch cluster. Default value is false",
        Required:       false,
        Type:           "bool",
        DefaultValue:   "false",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
    },
  },
  Ports: []*doc.ProcessorPort{},
}
}