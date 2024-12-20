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
	"os"

	"google.golang.org/grpc"

	"github.com/gravitational/teleport/api/client"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg := Config{
		proxyAddr:         os.Getenv("PROXY_ADDR"),
		discordwebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		watcherList:       os.Getenv("WATCHER_LIST"),
	}

	// Create a new Teleport client
	ctx := context.Background()
	creds := client.LoadIdentityFile("auth.pem")
	teleport, err := client.New(ctx, client.Config{
		Addrs:       []string{cfg.proxyAddr},
		Credentials: []client.Credentials{creds},
		DialOpts: []grpc.DialOption{
			grpc.WithReturnConnectionError(),
		},
	})
	if err != nil {
		panic(err)
	}
	defer teleport.Close()

	// Create a new Discord client
	discordClient := discordClient{
		url: cfg.discordwebhookURL,
		cfg: cfg,
	}

	// Create a new AccessRequestPlugin
	plugin := AccessRequestPlugin{
		TeleportClient: teleport,
		EventHandler:   &discordClient,
	}

	if err := plugin.Run(); err != nil {
		panic(err)
	}
}
