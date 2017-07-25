package opentracing

import (
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/cihub/seelog"

	"golang.org/x/net/context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"google.golang.org/grpc/metadata"
)

func Test_TraceGRPCRequestRoundtrip(t *testing.T) {
	assert := assert.New(t)
	tracer := mocktracer.New()

	// Initialize the ctx with a Span to inject.
	beforeSpan := tracer.StartSpan("to_inject").(*mocktracer.MockSpan)
	defer beforeSpan.Finish()
	beforeSpan.SetBaggageItem("key", "dongjiang")
	beforeCtx := opentracing.ContextWithSpan(context.Background(), beforeSpan)

	toGRPCFunc := ToGRPCRequest(tracer)
	md := metadata.Pairs()
	log.Info(md)
	// Call the RequestFunc.
	afterCtx := toGRPCFunc(beforeCtx, &md)

	// The Span should not have changed.
	afterSpan := opentracing.SpanFromContext(afterCtx)
	assert.Equal(beforeSpan, afterSpan)

	// No spans should have finished yet.
	finishedSpans := tracer.FinishedSpans()
	assert.Equal(0, len(finishedSpans))

	// Use FromGRPCRequest to verify that we can join with the trace given MD.
	fromGRPCFunc := FromGRPCRequest(tracer, "joined")
	joinCtx := fromGRPCFunc(afterCtx, md)
	joinedSpan := opentracing.SpanFromContext(joinCtx).(*mocktracer.MockSpan)

	joinedContext := joinedSpan.Context().(mocktracer.MockSpanContext)
	beforeContext := beforeSpan.Context().(mocktracer.MockSpanContext)

	assert.NotEqual(joinedContext.SpanID, beforeContext.SpanID)

	assert.Equal(beforeContext.SpanID, joinedSpan.ParentID)

	assert.Equal("joined", joinedSpan.OperationName)

	assert.Equal("dongjiang", joinedSpan.BaggageItem("key"))
	assert.Equal("", joinedSpan.BaggageItem("aaaa"))
}
