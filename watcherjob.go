// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gravitational/trace"

	"github.com/gravitational/teleport/api/types"
)

func SendDiscordNotification(token string, channelID string, message string) error {
	// Create a new Discord session using the provided bot token
	// Note: token should include "Bot " prefix if it's a bot token
	dg, err := discordgo.New(token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %v", err)
	}
	defer dg.Close()

	// Send the message
	_, err = dg.ChannelMessageSend(channelID, message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}

func (d *discordClient) HandleEvent(ctx context.Context, event types.Event) error {
	if event.Resource == nil {
		return nil
	}

	if _, ok := event.Resource.(*types.WatchStatusV1); ok {
		fmt.Println("Successfully started listening for Access Requests...")
		return nil
	}

	r, ok := event.Resource.(types.AccessRequest)
	if !ok {
		fmt.Printf("Unknown (%v) event received, skipping.\n", event.Resource)
		return nil
	}

	switch r.GetState() {
	case types.RequestState_PENDING:
		if err := d.handleNewAccessRequest(r); err != nil {
			return err
		}
	case types.RequestState_APPROVED:
		if err := d.handleApproveAccessRequest(r); err != nil {
			return err
		}
	case types.RequestState_DENIED:
		if err := d.handleDenyAccessRequest(r); err != nil {
			return err
		}
	}

	return nil
}

func (p *AccessRequestPlugin) Run() error {
	ctx := context.Background()

	watch, err := p.TeleportClient.NewWatcher(ctx, types.Watch{
		Kinds: []types.WatchKind{
			types.WatchKind{Kind: types.KindAccessRequest},
		},
	})

	if err != nil {
		return trace.Wrap(err)
	}
	defer watch.Close()

	fmt.Println("Starting the watcher job")

	for {
		select {
		case e := <-watch.Events():
			if err := p.EventHandler.HandleEvent(ctx, e); err != nil {
				return trace.Wrap(err)
			}
		case <-watch.Done():
			fmt.Println("The watcher job is finished")
			return nil
		}
	}
}
