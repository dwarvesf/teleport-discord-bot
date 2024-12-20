package teleport

import (
	"context"
	"fmt"

	"github.com/gravitational/teleport/api/client"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/trace"
	"google.golang.org/grpc"

	"github.com/dwarvesf/teleport-discord-bot/internal/config"
)

// EventHandler defines the interface for handling Teleport events
type EventHandler interface {
	HandleNewAccessRequest(r types.AccessRequest) error
	HandleApproveAccessRequest(r types.AccessRequest) error
	HandleDenyAccessRequest(r types.AccessRequest) error
}

// Plugin manages Teleport access request monitoring
type Plugin struct {
	client       *client.Client
	eventHandler EventHandler
}

// NewPlugin creates a new Teleport access request plugin
func NewPlugin(cfg *config.Config, eventHandler EventHandler) (*Plugin, error) {
	ctx := context.Background()

	// Create a new Teleport client
	teleportClient, err := client.New(ctx, client.Config{
		Addrs: []string{cfg.ProxyAddr},
		Credentials: []client.Credentials{
			client.LoadIdentityFile(cfg.AuthPemPath),
		},
		DialOpts: []grpc.DialOption{
			grpc.WithReturnConnectionError(),
		},
	})
	if err != nil {
		return nil, trace.Wrap(err, "failed to create Teleport client")
	}

	return &Plugin{
		client:       teleportClient,
		eventHandler: eventHandler,
	}, nil
}

// Run starts watching for access request events
func (p *Plugin) Run(ctx context.Context) error {
	// Create a new watcher for access requests
	watch, err := p.client.NewWatcher(ctx, types.Watch{
		Kinds: []types.WatchKind{
			{Kind: types.KindAccessRequest},
		},
	})
	if err != nil {
		return trace.Wrap(err, "failed to create watcher")
	}
	defer watch.Close()

	fmt.Println("Starting the access request watcher")

	for {
		select {
		case e := <-watch.Events():
			if err := p.handleEvent(ctx, e); err != nil {
				return trace.Wrap(err, "error handling event")
			}
		case <-watch.Done():
			fmt.Println("The watcher job is finished")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// handleEvent processes individual Teleport events
func (p *Plugin) handleEvent(ctx context.Context, event types.Event) error {
	if event.Resource == nil {
		return nil
	}

	// Check if this is a watch status event (initial connection)
	if _, ok := event.Resource.(*types.WatchStatusV1); ok {
		fmt.Println("Successfully started listening for Access Requests...")
		return nil
	}

	// Try to cast the resource to an AccessRequest
	r, ok := event.Resource.(types.AccessRequest)
	if !ok {
		fmt.Printf("Unknown (%v) event received, skipping.\n", event.Resource)
		return nil
	}

	// Handle different access request states
	switch r.GetState() {
	case types.RequestState_PENDING:
		return p.eventHandler.HandleNewAccessRequest(r)
	case types.RequestState_APPROVED:
		return p.eventHandler.HandleApproveAccessRequest(r)
	case types.RequestState_DENIED:
		return p.eventHandler.HandleDenyAccessRequest(r)
	default:
		fmt.Printf("Unhandled access request state: %v\n", r.GetState())
		return nil
	}
}

// Close terminates the Teleport client connection
func (p *Plugin) Close() {
	if p.client != nil {
		p.client.Close()
	}
}
