package load_balancing

import (
	"context"
	"io"
	"time"

	"github.com/ermes-labs/api-go/api"
)

//lint:ignore U1000 Ignore unused function temporarily for debugging
func create_session_api_mock(ip string) api.SessionToken {
	// This function is a placeholder for the API mock.
	node := nodesByIp[ip]
	var sessionToken *api.SessionToken
	node.CreateSession(context.Background(), &sessionToken, api.DefaultCreateSessionOptions())
	return *sessionToken
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func simulate_function_invocation_api_mock(ip string, sessionToken *api.SessionToken) *api.SessionToken {
	// This function is a placeholder for the API mock.
	node := nodesByIp[ip]
	node.SetResourcesTotalUsage(context.Background(), *sessionToken, map[string]float64{
		"resource": 0.025,
	})
	node.SetSessionUsageTime(context.Background(), *sessionToken, time.Duration(500)*time.Millisecond)

	migration_trigger(node)

	return check_migrated(node, *sessionToken)
}

func check_migrated(node api.Node, sessionToken api.SessionToken) *api.SessionToken {
	// Check is the session has been migrated
	_, err := node.GetSessionMetadata(context.Background(), sessionToken.SessionId)
	if errMigrated, ok := err.(*api.SessionMigratedError); ok {
		// The session has been migrated
		var newToken = errMigrated.SessionToken()
		return &newToken
	}

	return nil
}

func migration_trigger(node api.Node) {
	for {
		lookupNodeInfo, sessions, err := node.SessionsToMigrate(context.Background(), api.DefaultBestOffloadTargetsOptions())
		if len(sessions) == 0 || err != nil {
			return
		}

		lookupNode := nodesByIp[lookupNodeInfo.AreaName]
		sessionsToNodesMap, err := lookupNode.BestOffloadTargetNodes(context.Background(), node.Host, sessions, api.DefaultBestOffloadTargetsOptions())
		if err != nil {
			return
		}

		// For each session-node couple, try to offload the session.
		for _, target := range sessionsToNodesMap {
			sessionId, nodeId := target[0], target[1]
			target := nodes[nodeId]

			node.MigrateSession(
				context.Background(),
				sessionId,
				api.DefaultOffloadSessionOptions(),
				func(ctx context.Context, metadata api.SessionMetadata, reader io.Reader) (api.SessionLocation, error) {
					return target.ReceiveSessionMigration(context.Background(), metadata, reader, api.DefaultOnloadSessionOptions())
				},
				func(ctx context.Context, lastVisitedLocation, newLocation api.SessionLocation) (bool, error) {
					lastVisitedNode := nodesByIp[lastVisitedLocation.Host]
					return lastVisitedNode.UpdateOffloadedSessionLocation(context.Background(), lastVisitedLocation.SessionId, newLocation)
				})
		}
	}
}
