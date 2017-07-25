package opentracing

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	"endpoint"

	"github.com/stretchr/testify/assert"

	log "github.com/cihub/seelog"
)

func Test_ServerTrace(t *testing.T) {
	assert := assert.New(t)

	tracer := mocktracer.New()
	log.Info(tracer)

	// Initialize the ctx with a nameless Span.
	contextSpan := tracer.StartSpan("123").(*mocktracer.MockSpan)
	ctx := opentracing.ContextWithSpan(context.Background(), contextSpan)

	tracedEndpoint := TraceServer(tracer, "dongjiang")(endpoint.Nop)
	_, err := tracedEndpoint(ctx, struct{}{})
	assert.Nil(err)

	finishedSpans := tracer.FinishedSpans()
	assert.Equal(1, len(finishedSpans))

	// Test that the op name is updated
	endpointSpan := finishedSpans[0]
	assert.Equal("dongjiang", endpointSpan.OperationName)

	contextContext := contextSpan.Context().(mocktracer.MockSpanContext)
	endpointContext := endpointSpan.Context().(mocktracer.MockSpanContext)
	assert.Equal(contextContext.SpanID, endpointContext.SpanID)
}

func Test_TraceServerNoContextSpan(t *testing.T) {
	assert := assert.New(t)

	tracer := mocktracer.New()

	// Empty/background context.
	tracedEndpoint := TraceServer(tracer, "dongjiang")(endpoint.Nop)
	_, err := tracedEndpoint(context.Background(), struct{ aa string }{aa: "aass"})
	assert.Nil(err)

	// tracedEndpoint created a new Span.
	finishedSpans := tracer.FinishedSpans()
	assert.Equal(1, len(finishedSpans))

	endpointSpan := finishedSpans[0]
	assert.Equal("dongjiang", endpointSpan.OperationName)
}

func Test_Client(t *testing.T) {
	assert := assert.New(t)

	tracer := mocktracer.New()
	// Initialize the ctx with a parent Span.
	parentSpan := tracer.StartSpan("parent").(*mocktracer.MockSpan)
	defer parentSpan.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), parentSpan)

	tracedEndpoint := TraceClient(tracer, "dongjiang")(endpoint.Nop)
	_, err := tracedEndpoint(ctx, struct{ test string }{test: "this is test!"})
	assert.Nil(err)

	// tracedEndpoint created a new Span.
	finishedSpans := tracer.FinishedSpans()
	assert.Equal(1, len(finishedSpans))

	endpointSpan := finishedSpans[0]
	assert.Equal("dongjiang", endpointSpan.OperationName)

	parentContext := parentSpan.Context().(mocktracer.MockSpanContext)
	endpointContext := parentSpan.Context().(mocktracer.MockSpanContext)

	assert.Equal(parentContext.SpanID, endpointContext.SpanID)
}

func Test_TraceClientNoContextSpan(t *testing.T) {
	assert := assert.New(t)
	tracer := mocktracer.New()

	// Empty/background context.
	tracedEndpoint := TraceClient(tracer, "dongjiang")(endpoint.Nop)
	_, err := tracedEndpoint(context.Background(), struct{ test string }{test: "this is testing!"})
	assert.Nil(err)

	// tracedEndpoint created a new Span.
	finishedSpans := tracer.FinishedSpans()
	assert.Equal(1, len(finishedSpans))

	endpointSpan := finishedSpans[0]
	assert.Equal("dongjiang", endpointSpan.OperationName)
}
