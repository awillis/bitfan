//go:generate bitfanDoc
// When used in input (input->filter->o) the processor will receive events from the last filter from the pipeline used in configuration file.
//
// When used in filter (i->filter->o) the processor will
//
// * pass the event to the first filter plugin found in the used configuration file
// * receive events from the last filter plugin found in the used configuration file
//
// When used in output (i->filter->output->o) the processor will
//
// * pass the event to the first filter plugin found in the used configuration file
package use

import "bitfan/processors"

const (
	PORT_SUCCESS = 0
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// Path to configuration to import in this pipeline, it could be a local file or an url
	// can be relative path to the current configuration.
	//
	// SPLIT and JOIN : in filter Section, set multiples path to make a split and join into your pipeline
	// @ExampleLS path=> ["meteo-input.conf"]
	Path []string `mapstructure:"path" validate:"required"`

	// You can set variable references in the used configuration by using ${var}.
	// each reference will be replaced by the value of the variable found in this option
	//
	// The replacement is case-sensitive.
	// @ExampleLS var => {"hostname"=>"myhost","varname"=>"varvalue"}
	Var map[string]string `mapstructure:"var"`
}

// Include a config file
type processor struct {
	processors.Base

	opt *options
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt = &options{}
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	p.opt.ProcessCommonOptions(e.Fields())

	p.Send(e, PORT_SUCCESS)
	return nil
}
