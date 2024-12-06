package mobile_client

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/ermes-labs/api-go/api"
	"github.com/ermes-labs/api-go/infrastructure"
)

//lint:ignore U1000 Ignore unused function temporarily for debugging

// Client struct to represent the client
type Client struct {
	Host         string
	SessionToken *api.SessionToken
	Position     infrastructure.GeoCoordinates
}

// LocationChange struct to store location change details
type LocationChange struct {
	Latitude  float64
	Longitude float64
	NewIP     string
}

var locationChanges []LocationChange

func (c *Client) moveAndTrack(endLat, stepKm float64) {
	for lat := c.Position.Latitude; lat <= endLat; lat += stepKmToLat(stepKm) {
		c.Position.Latitude = lat
		newToken := update_location_api_mock(c.Host, c.SessionToken, c.Position)
		if newToken != nil && (c.SessionToken == nil || newToken.Host != c.SessionToken.Host) {
			locationChanges = append(locationChanges, LocationChange{
				Latitude:  lat,
				Longitude: c.Position.Longitude,
				NewIP:     newToken.Host,
			})
			c.SessionToken = newToken
			c.Host = newToken.Host
			log.Printf("Position updated to lat: %f, long: %f, new IP: %s", lat, c.Position.Longitude, newToken.Host)
		}
	}
}

// Converts km step to equivalent latitude degrees (approximately)
func stepKmToLat(km float64) float64 {
	return km / 111.0 // 1 degree of latitude is approximately 111 km
}

func main() {
	startLat := 29.5
	endLat := 48.0
	long := -98.1
	stepKm := 1.0

	client := Client{
		Host: "<endpoint-1.4>",
		Position: infrastructure.GeoCoordinates{
			Latitude:  startLat,
			Longitude: long,
		},
	}

	// Create initial session and save the initial token and position
	token := create_session_api_mock(client.Host)
	client.SessionToken = &token
	locationChanges = append(locationChanges, LocationChange{
		Latitude:  client.Position.Latitude,
		Longitude: client.Position.Longitude,
		NewIP:     client.SessionToken.Host,
	})
	log.Printf("Initial session created at lat: %f, long: %f, IP: %s", client.Position.Latitude, client.Position.Longitude, client.SessionToken.Host)

	// Start moving and tracking
	client.moveAndTrack(endLat, stepKm)

	// Save location changes to a CSV file
	saveLocationChanges("location_changes.csv")
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func saveLocationChanges(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Latitude", "Longitude", "NewIP"})

	for _, lc := range locationChanges {
		writer.Write([]string{
			fmt.Sprintf("%f", lc.Latitude),
			fmt.Sprintf("%f", lc.Longitude),
			lc.NewIP,
		})
	}
}
