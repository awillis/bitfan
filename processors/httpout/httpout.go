//go:generate bitfanDoc
// Display on http the last received event
//
// URL is available as http://webhookhost/pipelineName/pluginLabel/URI
//
// * webhookhost is defined by bitfan at startup
// * pluginLabel is defined in pipeline configuration, it's the named processor if you put one, or `httpout` by default
// * URI is defined in plugin configuration (see below)
package httpoutprocessor

import (
	"fmt"
	"net/http"
	"os"

	"github.com/awillis/bitfan/codecs"
	"github.com/awillis/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// @Default "json"
	// @Type codec
	Codec codecs.CodecCollection

	// URI path
	// @Default "out"
	Uri string

	// Add headers to output
	// @default {"Content-Type" => "application/json"}
	Headers map[string]string
}

// Reads events from standard input
type processor struct {
	processors.Base

	opt  *options
	q    chan bool
	host string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec: codecs.CodecCollection{
			Enc: codecs.New("json", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
		Uri: "out",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.host, err = os.Hostname(); err != nil {
		p.Logger.Warnf("can not get hostname : %v", err)
	}

	return err
}

func (p *processor) Start(e processors.IPacket) error {
	p.WebHook.Add(p.opt.Uri, p.HttpHandler)
	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	p.Memory.Set("last", e)
	p.Send(e)
	return nil
}

// Handle Request received by bitfan for this agent (url hook should be registered during p.Start)
func (p *processor) HttpHandler(w http.ResponseWriter, r *http.Request) {
	p.Logger.Debug("reading request")
	last, ok := p.Memory.Get("last")
	if !ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("Nothing received !")))
		return
	}

	// Encode content
	var err error
	var enc codecs.Encoder
	enc, err = p.opt.Codec.NewEncoder(w)
	if err != nil {
		p.Logger.Errorln("codec error : ", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	for ctn, ctv := range p.opt.Headers {
		w.Header().Set(ctn, ctv)
	}
	enc.Encode(last.(processors.IPacket).Fields().Old())
}
