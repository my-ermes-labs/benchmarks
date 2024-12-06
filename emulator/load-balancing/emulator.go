package load_balancing

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ermes-labs/api-go/api"
	"github.com/ermes-labs/api-go/infrastructure"
)

//lint:ignore U1000 Ignore unused function temporarily for debugging
var sessionMap map[string]int = map[string]int{
	"<endpoint-1>":   0,
	"<endpoint-1.1>": 0,
	"<endpoint-1.2>": 0,
}

// Client struct to represent the client
type Client struct {
	Host         string
	SessionToken *api.SessionToken
	Position     infrastructure.GeoCoordinates
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func generateNClientsBetween(numPoints int, node1, node2 api.Node, host string) []Client {
	// Generate n points between node1 and node2
	points := generateNPointsBetween(numPoints, node1.GeoCoordinates, distance(node1.GeoCoordinates, node2.GeoCoordinates))

	clients := []Client{}
	// Generate nodes from the points
	for _, point := range points {
		client := Client{
			Host:         host,
			Position:     point,
			SessionToken: nil,
		}

		clients = append(clients, client)
	}

	return clients
}

func generateNPointsBetween(numPoints int, center infrastructure.GeoCoordinates, radius float64) []infrastructure.GeoCoordinates {
	var points []infrastructure.GeoCoordinates

	for i := 0; i < numPoints; i++ {
		// Generate random angle
		randAngle := 2 * math.Pi * rand.Float64()

		// Generate random distance from center within the radius
		randDist := radius * math.Sqrt(rand.Float64())

		// Calculate x and y coordinates of the point
		Longitude := randDist*math.Cos(randAngle) + center.Longitude
		Latitude := randDist*math.Sin(randAngle) + center.Latitude

		points = append(points, infrastructure.GeoCoordinates{
			Latitude:  Latitude,
			Longitude: Longitude,
		})
	}

	return points
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func randomlyInitiateSessions(groups ...[]Client) {
	for _, group := range groups {
		for _, client := range group {
			sessionToken := create_session_api_mock(client.Host)
			client.SessionToken = &sessionToken

			// Update the session count for the node
			sessionMap[client.Host]++
		}
	}
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func printSessions(nodes []api.Node) {
	// Print the areaName of the nodes as the column names,
	// Then print the number of sessions per node
	fmt.Println("AreaName, Sessions")
	for _, node := range nodes {
		fmt.Printf("%s, %d\n", node.AreaName, sessionMap[node.Host])
	}
}

func simulateWeightedTraffic(weighs []int, clients [][]Client) {
	// Simulate traffic between the clients
	for i, group := range clients {
		for _, client := range group {
			weight := weighs[i]

			for j := 0; j < weight; j++ {
				token := simulate_function_invocation_api_mock(client.Host, client.SessionToken)

				if token != nil {
					sessionMap[client.Host]--
					sessionMap[token.Host]++

					client.SessionToken = token
					client.Host = token.Host
				}
			}
		}
	}
}

func distance(pint1, point2 infrastructure.GeoCoordinates) float64 {
	return math.Sqrt(math.Pow(pint1.Latitude-point2.Latitude, 2) + math.Pow(pint1.Longitude-point2.Longitude, 2))
}
