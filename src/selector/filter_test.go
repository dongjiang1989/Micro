package selector

import (
	"testing"

	"registry"

	"github.com/stretchr/testify/assert"

	log "github.com/cihub/seelog"
)

func Test_FilterEndpoint(t *testing.T) {
	assert := assert.New(t)

	testData := []struct {
		services []*registry.Service
		endpoint string
		count    int
	}{
		{
			services: []*registry.Service{
				&registry.Service{
					Name:    "dongjiang",
					Version: "1.0.0",
					Endpoints: []*registry.Endpoint{
						&registry.Endpoint{
							Name: "Foo.Bar",
						},
					},
				},
				&registry.Service{
					Name:    "dongjiang",
					Version: "1.1.0",
					Endpoints: []*registry.Endpoint{
						&registry.Endpoint{
							Name: "Baz.Bar",
						},
					},
				},
			},
			endpoint: "Foo.Bar",
			count:    1,
		},
		{
			services: []*registry.Service{
				&registry.Service{
					Name:    "dongjiang",
					Version: "1.0.0",
					Endpoints: []*registry.Endpoint{
						&registry.Endpoint{
							Name: "Foo.Bar",
						},
					},
				},
				&registry.Service{
					Name:    "dongjiang",
					Version: "1.1.0",
					Endpoints: []*registry.Endpoint{
						&registry.Endpoint{
							Name: "Foo.Bar",
						},
					},
				},
			},
			endpoint: "Bar.Baz",
			count:    0,
		},
	}

	for _, data := range testData {
		filter := FilterEndpoint(data.endpoint)
		services := filter(data.services)
		assert.Equal(len(services), data.count)

		for _, service := range services {
			var seen bool

			for _, ep := range service.Endpoints {
				if ep.Name == data.endpoint {
					seen = true
					break
				}
			}
			assert.False(seen == false && data.count > 0)
		}
	}
}

func Test_FilterLabel(t *testing.T) {
	assert := assert.New(t)
	testData := []struct {
		services []*registry.Service
		label    [2]string
		count    int
	}{
		{
			services: []*registry.Service{
				&registry.Service{
					Name:    "dongjiang",
					Version: "1.0.0",
					Nodes: []*registry.Node{
						&registry.Node{
							Id:      "dongjiang-1",
							Address: "localhost",
							Metadata: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
				&registry.Service{
					Name:    "dongjiang",
					Version: "1.1.0",
					Nodes: []*registry.Node{
						&registry.Node{
							Id:      "dongjiang-2",
							Address: "localhost",
							Metadata: map[string]string{
								"foo": "baz",
							},
						},
					},
				},
			},
			label: [2]string{"foo", "bar"},
			count: 1,
		},
		{
			services: []*registry.Service{
				&registry.Service{
					Name:    "dongjiang",
					Version: "1.0.0",
					Nodes: []*registry.Node{
						&registry.Node{
							Id:      "dongjiang-1",
							Address: "localhost",
						},
					},
				},
				&registry.Service{
					Name:    "dongjiang",
					Version: "1.1.0",
					Nodes: []*registry.Node{
						&registry.Node{
							Id:      "dongjiang-2",
							Address: "localhost",
						},
					},
				},
			},
			label: [2]string{"foo", "bar"},
			count: 0,
		},
	}

	for _, data := range testData {
		filter := FilterLabel(data.label[0], data.label[1])
		services := filter(data.services)

		assert.Equal(len(services), data.count)

		for _, service := range services {
			var seen bool

			for _, node := range service.Nodes {
				assert.Equal(node.Metadata[data.label[0]], data.label[1])
				seen = true
			}
			assert.True(seen)
		}
	}
}

func Test_FilterVersion(t *testing.T) {
	assert := assert.New(t)
	testData := []struct {
		services []*registry.Service
		version  string
		count    int
	}{
		{
			services: []*registry.Service{
				&registry.Service{
					Name:    "test",
					Version: "1.0.0",
				},
				&registry.Service{
					Name:    "test",
					Version: "1.1.0",
				},
			},
			version: "1.0.0",
			count:   1,
		},
		{
			services: []*registry.Service{
				&registry.Service{
					Name:    "test",
					Version: "1.0.0",
				},
				&registry.Service{
					Name:    "test",
					Version: "1.1.0",
				},
			},
			version: "2.0.0",
			count:   0,
		},
	}

	for _, data := range testData {
		filter := FilterVersion(data.version)
		services := filter(data.services)
		assert.Equal(len(services), data.count)

		var seen bool

		for _, service := range services {
			assert.Equal(service.Version, data.version)
			seen = true
		}

		log.Info("dongjiang:", data.count, seen)

		assert.False(seen == false && data.count > 0)
	}
}
