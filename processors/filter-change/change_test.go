package change

import (
	"testing"
	"time"

	"bitfan/processors/doc"
	"bitfan/processors/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a processor")
}
func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}
func TestMaxConcurent(t *testing.T) {
	max := New().(*processor).MaxConcurent()
	assert.Equal(t, 1, max, "this processor does not support concurency")
}

func TestReceiveMatch(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
		},
	)

	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")
	p.Receive(testutils.NewPacketOld("test", nil))
	p.Receive(testutils.NewPacketOld("test", nil))
	p.Receive(testutils.NewPacketOld("test", nil))
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 1")
	p.Receive(testutils.NewPacketOld("toto", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 2")
	p.Receive(testutils.NewPacketOld("toto", nil))
	p.Receive(testutils.NewPacketOld("toto", nil))
	p.Receive(testutils.NewPacketOld("toto", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 3")
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 2, ctx.SentPacketsCount(0), "changed ! 4")
}

func TestReceiveIgnoreMissingTrue(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field":  "toto",
			"Ignore_missing": true,
		},
	)

	p.Receive(testutils.NewPacketOld("test", map[string]interface{}{"toto": "A"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 1")
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 2")
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 3")
	p.Receive(testutils.NewPacketOld("test", map[string]interface{}{"toto": "A"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 4")
	p.Receive(testutils.NewPacketOld("test", map[string]interface{}{"toto": "A"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 5")
	p.Receive(testutils.NewPacketOld("test", map[string]interface{}{"toto": "B"}))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 6")
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 7")
}

func TestReceiveIgnoreMissingFalse(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field":  "toto",
			"Ignore_missing": false,
		},
	)

	p.Receive(testutils.NewPacketOld("test", map[string]interface{}{"toto": "A"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 1")
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 2")
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 3")
	p.Receive(testutils.NewPacketOld("test", map[string]interface{}{"toto": "A"}))
	assert.Equal(t, 2, ctx.SentPacketsCount(0), "changed ! 4")
	p.Receive(testutils.NewPacketOld("test", map[string]interface{}{"toto": "A"}))
	assert.Equal(t, 2, ctx.SentPacketsCount(0), "changed ! 5")
	p.Receive(testutils.NewPacketOld("test", map[string]interface{}{"toto": "B"}))
	assert.Equal(t, 3, ctx.SentPacketsCount(0), "changed ! 6")
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 4, ctx.SentPacketsCount(0), "changed ! 7")
}

func TestStopNoTimeFrame(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
			"Timeframe":     0,
		},
	)

	p.Receive(testutils.NewPacketOld("test", nil))
	p.Receive(testutils.NewPacketOld("test2", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 1")
	assert.NoError(t, p.Stop(nil), "no error")
}

func TestStopWithTimeFrame(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
			"Timeframe":     3,
		},
	)

	p.Receive(testutils.NewPacketOld("test", nil))
	p.Receive(testutils.NewPacketOld("test2", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 1")
	time.Sleep(time.Second * 1)
	assert.NoError(t, p.Stop(nil), "no error")
}

func TestReceiveMatchWithTimeframe(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
			"timeframe":     1,
		},
	)

	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")

	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")

	p.Receive(testutils.NewPacketOld("test1", nil))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 0") {
		assert.Equal(t, "test1", ctx.SentPackets(0)[0].Message())
	}

	p.Receive(testutils.NewPacketOld("test1", nil))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 0") {
		assert.Equal(t, "test1", ctx.SentPackets(0)[0].Message())
	}

	p.Receive(testutils.NewPacketOld("test1", nil))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 0") {
		assert.Equal(t, "test1", ctx.SentPackets(0)[0].Message())
	}

	time.Sleep(time.Second * 2)
	p.Receive(testutils.NewPacketOld("test1", nil))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 0") {
		assert.Equal(t, "test1", ctx.SentPackets(0)[0].Message())
	}

	time.Sleep(time.Second * 2)
	p.Receive(testutils.NewPacketOld("test1", nil))
	p.Receive(testutils.NewPacketOld("test1", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 0")

	p.Receive(testutils.NewPacketOld("A", nil))
	assert.Equal(t, 2, ctx.SentPacketsCount(0), "changed ! 0")

	p.Receive(testutils.NewPacketOld("B", nil))
	assert.Equal(t, 3, ctx.SentPacketsCount(0), "changed ! 0")

	time.Sleep(time.Second * 2)
	p.Receive(testutils.NewPacketOld("B", nil))
	assert.Equal(t, 3, ctx.SentPacketsCount(0), "changed ! 0")

	p.Receive(testutils.NewPacketOld("B", nil))
	p.Receive(testutils.NewPacketOld("test", nil))
	time.Sleep(time.Second * 2)
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 4, ctx.SentPacketsCount(0), "changed ! 0")

}

func TestReceiveRepetitionsWithinTimeFrameDoNotMatch(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
			"timeframe":     3,
		},
	)

	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")

	// 3 repetitions in a row
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")

	// 1 repetition per second
	time.Sleep(time.Second * 1)
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")

	time.Sleep(time.Second * 1)
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")

	time.Sleep(time.Second * 1)
	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")
}

func TestReceiveChangesOutOfTimeFrameDoNotMatch(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
			"timeframe":     1,
		},
	)

	p.Receive(testutils.NewPacketOld("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")

	// Changes happen out of time frame
	time.Sleep(time.Second * 2)
	p.Receive(testutils.NewPacketOld("test2", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")

	time.Sleep(time.Second * 2)
	p.Receive(testutils.NewPacketOld("test3", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")

	time.Sleep(time.Second * 2)
	p.Receive(testutils.NewPacketOld("test4", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")
}
