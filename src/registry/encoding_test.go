package registry

import (
	"encoding/json"

	"testing"

	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"
)

func Test_EncodingEndpoints(t *testing.T) {
	assert := assert.New(t)
	eps := []*Endpoint{
		&Endpoint{
			Name: "endpoint-dongjiang1",
			Request: &Value{
				Name: "request",
				Type: "request",
			},
			Response: &Value{
				Name: "response",
				Type: "response",
			},
			Metadata: map[string]string{
				"key2": "value1",
			},
		},
		&Endpoint{
			Name: "endpoint-dongjiang2",
			Request: &Value{
				Name: "request",
				Type: "request",
			},
			Response: &Value{
				Name: "response",
				Type: "response",
			},
			Metadata: map[string]string{
				"key2": "value2",
			},
		},
		&Endpoint{
			Name: "endpoint-dongjiang3",
			Request: &Value{
				Name: "request",
				Type: "request",
			},
			Response: &Value{
				Name: "response",
				Type: "response",
			},
			Metadata: map[string]string{
				"key3": "value3",
			},
		},
	}

	testEp := func(ep *Endpoint, enc string) {
		log.Critical(enc)
		e := encodeEndpoints([]*Endpoint{ep})

		assert.Equal(1, len(e))

		var seen bool

		for _, en := range e {
			if en == enc {
				seen = true
				break
			}
		}
		assert.True(seen)

		// decode
		d := decodeEndpoints([]string{enc})
		log.Critical(d, len(d))
		assert.NotEqual(0, len(d))
		assert.Equal(ep.Name, d[0].Name)

		// check all the metadata exists
		for k, v := range ep.Metadata {
			assert.Equal(v, d[0].Metadata[k])
		}
	}

	for _, ep := range eps {
		// JSON encoded
		jencoded, err := json.Marshal(ep)
		log.Info(ep, jencoded)
		assert.Nil(err)

		// HEX encoded
		hencoded := encode(jencoded)
		// endpoint tag
		hepTag := "e-" + hencoded
		log.Info(hepTag)
		testEp(ep, hepTag)
	}
}

func Test_EncodingVersion(t *testing.T) {
	assert := assert.New(t)
	testData := []struct {
		decoded string
		encoded string
	}{
		{"1.0.0", "v-789c32d433d03300040000ffff02ce00ee"},
		{"latest", "v-789cca492c492d2e01040000ffff08cc028e"},
	}

	for _, data := range testData {
		e := encodeVersion(data.decoded)

		assert.Equal(e[0], data.encoded)

		d, ok := decodeVersion(e)
		assert.True(ok)
		assert.Equal(data.decoded, d)

		d, ok = decodeVersion([]string{data.encoded})
		assert.True(ok)
		assert.Equal(data.decoded, d)
	}

	log.Info(testData)
}
