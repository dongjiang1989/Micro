package selector

import (
	"testing"

	"registry"

	"github.com/stretchr/testify/assert"

	log "github.com/cihub/seelog"
)

func Test_Strategies(t *testing.T) {
	assert := assert.New(t)

	testData := []*registry.Service{
		&registry.Service{
			Name:    "test1",
			Version: "latest",
			Nodes: []*registry.Node{
				&registry.Node{
					Id:      "test1-1",
					Address: "10.0.0.1",
					Port:    1001,
				},
				&registry.Node{
					Id:      "test1-2",
					Address: "10.0.0.2",
					Port:    1002,
				},
			},
		},
		&registry.Service{
			Name:    "test1",
			Version: "default",
			Nodes: []*registry.Node{
				&registry.Node{
					Id:      "test1-3",
					Address: "10.0.0.3",
					Port:    1003,
				},
				&registry.Node{
					Id:      "test1-4",
					Address: "10.0.0.4",
					Port:    1004,
				},
			},
		},
	}

	for name, strategy := range map[string]Strategy{"random": Random, "roundrobin": RoundRobin} {
		next := strategy(testData)
		counts := make(map[string]int)

		for i := 0; i < 100; i++ {
			node, err := next()
			assert.Nil(err)
			counts[node.Id]++
		}

		log.Info(name, ": ", counts)
	}
}
