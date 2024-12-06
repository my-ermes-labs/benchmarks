package load_balancing

import (
	"github.com/ermes-labs/api-go/api"
)

//lint:ignore U1000 Ignore unused function temporarily for debugging
var parentNode api.Node = api.Node{
	Cmd: NewCommands(),
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func test() {
	// Create the nodes
	cloud_node := nodes["Central-USA"]
	left_edge_node := nodes["Minnesota-Edge"]
	right_edge_node := nodes["Idaho-Edge"]
	// set of the three nodes
	initialGroup := []api.Node{cloud_node, left_edge_node, right_edge_node}

	// Generate the clients
	group1 := generateNClientsBetween(100, cloud_node, left_edge_node, "Central-USA")
	group2 := generateNClientsBetween(100, cloud_node, right_edge_node, "Central-USA")
	group3 := generateNClientsBetween(100, left_edge_node, right_edge_node, "Central-USA")

	randomlyInitiateSessions(group1, group2, group3)
	// Print the nodes
	printSessions(initialGroup)

	simulateWeightedTraffic([]int{1, 2, 1}, [][]Client{group1, group2, group3})

	// Print the nodes
	printSessions(initialGroup)
}
