package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

type mockRegistry struct {
	body   []byte
	status int
	err    error
	url    string
}

func encodeData(obj interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func newMockServer(rg *mockRegistry, l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc(rg.url, func(w http.ResponseWriter, r *http.Request) {
		if rg.err != nil {
			http.Error(w, rg.err.Error(), 500)
			return
		}
		w.WriteHeader(rg.status)
		w.Write(rg.body)
	})
	return http.Serve(l, mux)
}

// return *object, callbackfunc()
func newConsulTestRegistry(r *mockRegistry) (*consulRegistry, func()) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		// blurgh?!!
		panic(err.Error())
	}
	cfg := consul.DefaultConfig()
	cfg.Address = l.Addr().String()
	cl, _ := consul.NewClient(cfg)

	go newMockServer(r, l)

	return &consulRegistry{
			Address:  cfg.Address,
			Client:   cl,
			register: make(map[string]uint64),
		}, func() {
			l.Close()
		}
}

func newServiceList(svc []*consul.ServiceEntry) []byte {
	bts, _ := encodeData(svc)
	return bts
}

func Test_Consul_GetService_WithError(t *testing.T) {
	assert := assert.New(t)
	cr, cl := newConsulTestRegistry(&mockRegistry{
		err: errors.New("client-error"),
		url: "/v1/health/service/service-name",
	})
	defer cl()

	_, err := cr.GetService("test-service")
	assert.NotNil(err)
}

func Test_ConsulGetService_WithHealthyServiceNodes(t *testing.T) {
	assert := assert.New(t)

	// warning is still seen as healthy, critical is not
	svcs := []*consul.ServiceEntry{
		newServiceEntry(
			"node-name-1", "node-address-1", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-1", "service-name", "passing"),
				newHealthCheck("node-name-1", "service-name", "warning"),
			},
		),
		newServiceEntry(
			"node-name-2", "node-address-2", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-2", "service-name", "passing"),
				newHealthCheck("node-name-2", "service-name", "warning"),
			},
		),
	}

	cr, cl := newConsulTestRegistry(&mockRegistry{
		status: 200,
		body:   newServiceList(svcs),
		url:    "/v1/health/service/service-name",
	})
	defer cl()

	svc, err := cr.GetService("service-name")
	assert.Nil(err)

	assert.Equal(1, len(svc))
	assert.Equal(2, len((svc[0]).Nodes))
}

func Test_Consul_GetService_WithUnhealthyServiceNode(t *testing.T) {
	assert := assert.New(t)
	// warning is still seen as healthy, critical is not
	svcs := []*consul.ServiceEntry{
		newServiceEntry(
			"node-name-1", "node-address-1", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-1", "service-name", "passing"),
				newHealthCheck("node-name-1", "service-name", "warning"),
			},
		),
		newServiceEntry(
			"node-name-2", "node-address-2", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-2", "service-name", "passing"),
				newHealthCheck("node-name-2", "service-name", "critical"),
			},
		),
	}

	cr, cl := newConsulTestRegistry(&mockRegistry{
		status: 200,
		body:   newServiceList(svcs),
		url:    "/v1/health/service/service-name",
	})
	defer cl()

	svc, err := cr.GetService("service-name")
	assert.Nil(err)
	assert.Equal(1, len(svc))
	assert.Equal(1, len(svc[0].Nodes))
}

func Test_ConsulGetService_WithUnhealthyServiceNodes(t *testing.T) {
	assert := assert.New(t)
	// warning is still seen as healthy, critical is not
	svcs := []*consul.ServiceEntry{
		newServiceEntry(
			"node-name-1", "node-address-1", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-1", "service-name", "passing"),
				newHealthCheck("node-name-1", "service-name", "critical"),
			},
		),
		newServiceEntry(
			"node-name-2", "node-address-2", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-2", "service-name", "passing"),
				newHealthCheck("node-name-2", "service-name", "critical"),
			},
		),
	}

	cr, cl := newConsulTestRegistry(&mockRegistry{
		status: 200,
		body:   newServiceList(svcs),
		url:    "/v1/health/service/service-name",
	})
	defer cl()

	svc, err := cr.GetService("service-name")
	assert.Nil(err)
	assert.Equal(1, len(svc))
	assert.Equal(0, len(svc[0].Nodes))
}
