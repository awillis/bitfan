package mongodb

import (
	"github.com/awillis/bitfan/processors/doc"
	"github.com/awillis/bitfan/processors/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a processor")
}
func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}
func TestMaxConcurrent(t *testing.T) {
	max := New().(*processor).MaxConcurrent()
	assert.Equal(t, 0, max, "this processor does support concurency")
}
func TestConfigure(t *testing.T) {
	conf := map[string]interface{}{}
	ctx := testutils.NewProcessorContext()
	p := New()
	err := p.Configure(ctx, conf)
	assert.Nil(t, err, "Configure() processor without error")
}
