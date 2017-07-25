package opentracing

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"golang.org/x/net/context"

	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func Test_TraceHTTPRequestRoundtrip(t *testing.T) {
	assert := assert.New(t)

	tracer := mocktracer.New()

	// Initialize the ctx with a Span to inject.
	beforeSpan := tracer.StartSpan("to_inject").(*mocktracer.MockSpan)
	defer beforeSpan.Finish()
	beforeSpan.SetBaggageItem("key", "dongjiang")
	beforeCtx := opentracing.ContextWithSpan(context.Background(), beforeSpan)

	toHTTPFunc := ToHTTPRequest(tracer)
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("aa", "bb")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("i am callback!"))
			}))

	req, _ := http.NewRequest("GET", server.URL, nil)
	// Call the RequestFunc.
	afterCtx := toHTTPFunc(beforeCtx, req)

	// The Span should not have changed.
	afterSpan := opentracing.SpanFromContext(afterCtx)
	assert.Equal(beforeSpan, afterSpan)

	// No spans should have finished yet.
	finishedSpans := tracer.FinishedSpans()
	assert.Equal(0, len(finishedSpans))

	// Use FromHTTPRequest to verify that we can join with the trace given a req.
	fromHTTPFunc := FromHTTPRequest(tracer, "joined")
	joinCtx := fromHTTPFunc(afterCtx, req)
	joinedSpan := opentracing.SpanFromContext(joinCtx).(*mocktracer.MockSpan)

	joinedContext := joinedSpan.Context().(mocktracer.MockSpanContext)
	beforeContext := beforeSpan.Context().(mocktracer.MockSpanContext)

	log.Info(joinedContext, beforeContext)

	assert.NotEqual(joinedContext.SpanID, beforeContext.SpanID)

	assert.Equal(beforeContext.SpanID, joinedSpan.ParentID)

	assert.Equal("joined", joinedSpan.OperationName)

	assert.Equal("dongjiang", joinedSpan.BaggageItem("key"))
}

func Test_ToHTTPRequestTags(t *testing.T) {
	assert := assert.New(t)
	tracer := mocktracer.New()

	span := tracer.StartSpan("to_inject").(*mocktracer.MockSpan)
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("aa", "bb")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("i am callback!"))
			}))

	req, _ := http.NewRequest("GET", server.URL, nil)

	ToHTTPRequest(tracer)(ctx, req)

	expectedTags := map[string]interface{}{
		string(ext.HTTPMethod):   "GET",
		string(ext.HTTPUrl):      server.URL,
		string(ext.PeerHostname): "127.0.0.1",
	}
	assert.Equal(expectedTags, span.Tags())
	assert.True(reflect.DeepEqual(expectedTags, span.Tags()))
}

func Test_FromHTTPRequestTags(t *testing.T) {
	assert := assert.New(t)
	tracer := mocktracer.New()
	parentSpan := tracer.StartSpan("to_extract").(*mocktracer.MockSpan)
	defer parentSpan.Finish()

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("aa", "bb")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("i am callback!"))
			}))

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.Nil(err)

	tracer.Inject(parentSpan.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(req.Header))

	ctx := FromHTTPRequest(tracer, "dongjiang")(context.Background(), req)
	opentracing.SpanFromContext(ctx).Finish()

	childSpan := tracer.FinishedSpans()[0]
	expectedTags := map[string]interface{}{
		string(ext.HTTPMethod): "GET",
		string(ext.HTTPUrl):    server.URL,
		string(ext.SpanKind):   ext.SpanKindRPCServerEnum,
	}

	assert.Equal(expectedTags, childSpan.Tags())
	assert.True(reflect.DeepEqual(expectedTags, childSpan.Tags()))
	assert.Equal("dongjiang", childSpan.OperationName)
}
