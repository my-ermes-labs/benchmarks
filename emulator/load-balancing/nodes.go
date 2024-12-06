package load_balancing

import (
	"io"
	"log"
	"os"

	"github.com/ermes-labs/api-go/api"
	"github.com/ermes-labs/api-go/infrastructure"
)

var nodes map[string]api.Node = map[string]api.Node{}
var nodesByIp map[string]api.Node = map[string]api.Node{}

func init() {
	// Read json file infrastructure ./infrastructure.json
	file, err := os.Open("./infrastructure.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read file bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON file
	_, areaMap, err := infrastructure.UnmarshalInfrastructure(fileBytes)
	if err != nil {
		log.Fatal(err)
	}

	// For each node in the map
	for id, area := range areaMap {
		// Add the node to the nodes map
		nodes[id] = *api.NewNode(
			area.Node,
			NewCommands(),
		)
	}

	// For each node in the map
	for _, area := range areaMap {
		// Add the node to the nodesByIp map
		nodesByIp[area.Node.Host] = *api.NewNode(
			area.Node,
			NewCommands(),
		)
	}
}
