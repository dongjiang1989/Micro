package registry

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	//	"golang.org/x/net/context"
)

func encode(buf []byte) string {
	var b bytes.Buffer
	defer b.Reset()

	w := zlib.NewWriter(&b)
	_, err := w.Write(buf)
	if err != nil {
		return ""
	}
	w.Close()

	return hex.EncodeToString(b.Bytes())
}

func decode(d string) []byte {
	hr, err := hex.DecodeString(d)
	if err != nil {
		return nil
	}

	b := bytes.NewReader(hr)
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil
	}

	buf, err := ioutil.ReadAll(z)
	if err != nil {
		return nil
	}
	return buf
}

func encodeEndpoints(ep []*Endpoint) []string {
	var tags []string
	for _, e := range ep {
		if b, err := json.Marshal(e); err == nil {
			tags = append(tags, "e-"+encode(b)) // add perfix : "e-"
		}
	}
	return tags
}

func decodeEndpoints(tags []string) []*Endpoint {
	var ep []*Endpoint
	var version byte

	for _, tag := range tags {
		// tag
		if len(tag) == 0 || tag[0] != 'e' { //perfix : "e-"
			continue
		}

		// version
		if version > 0 && tag[1] != version {
			continue
		}

		var e *Endpoint
		var buf []byte

		// Old encoding was plain
		if tag[1] == '=' {
			buf = []byte(tag[2:])
		}

		// New encoding is hex
		if tag[1] == '-' { //perfix : "e-"
			buf = decode(tag[2:])
		}

		if err := json.Unmarshal(buf, &e); err == nil {
			ep = append(ep, e)
		}

		// set version
		version = tag[1]
	}

	return ep
}

func encodeMetadata(md map[string]string) []string {
	var tags []string
	for k, v := range md {
		b, err := json.Marshal(map[string]string{
			k: v,
		})
		if err == nil {
			//TODO
			tags = append(tags, "t-"+encode(b)) // Metadata prefix: "t-"
		}
	}
	return tags
}

func decodeMetadata(tags []string) map[string]string {
	md := make(map[string]string)
	var version byte

	for _, tag := range tags {
		if len(tag) == 0 || tag[0] != 't' { // Metadata prefix: "t-"
			continue
		}

		// check version
		if version > 0 && tag[1] != version {
			continue
		}

		var kv map[string]string
		var buf []byte

		// Old encoding was plain
		if tag[1] == '=' {
			buf = []byte(tag[2:])
		}

		// New encoding is hex
		if tag[1] == '-' {
			buf = decode(tag[2:])
		}

		// Now unmarshal
		if err := json.Unmarshal(buf, &kv); err == nil {
			for k, v := range kv {
				md[k] = v
			}
		}

		// set version
		version = tag[1]
	}
	return md
}

func encodeVersion(v string) []string {
	return []string{"v-" + encode([]byte(v))} // version prefix: "t-"
}

func decodeVersion(tags []string) (string, bool) {
	for _, tag := range tags {
		if len(tag) < 2 || tag[0] != 'v' { // version prefix: "t-"
			continue
		}

		// Old encoding was plain
		if tag[1] == '=' {
			return tag[2:], true
		}

		// New encoding is hex
		if tag[1] == '-' { // version prefix: "t-"
			return string(decode(tag[2:])), true
		}
	}
	return "", false
}
