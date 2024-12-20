package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gravitational/teleport/api/types"
	"github.com/gtuk/discordwebhook"
)

func (d *discordClient) sendWebhookNotification(message discordwebhook.Message) error {
	// Send the webhook
	err := discordwebhook.SendMessage(d.url, message)
	if err != nil {
		return err
	}

	return nil
}

func ptrString(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}

func (d *discordClient) handleNewAccessRequest(r types.AccessRequest) error {
	fields := []discordwebhook.Field{
		{
			Name:   ptrString("Request ID"),
			Value:  ptrString(r.GetName()),
			Inline: ptrBool(false),
		},
		{
			Name:   ptrString("User"),
			Value:  ptrString(r.GetUser()),
			Inline: ptrBool(true),
		},
		{
			Name:   ptrString("Roles"),
			Value:  ptrString(strings.Join(r.GetRoles(), ", ")),
			Inline: ptrBool(true),
		},
		{
			Name:   ptrString("Session TTL"),
			Value:  ptrString(r.GetSessionTLL().Sub(r.GetCreationTime()).Round(time.Second).String()),
			Inline: ptrBool(true),
		},
	}

	if r.GetRequestReason() != "" {
		fields = append(fields, discordwebhook.Field{
			Name:   ptrString("Request Reason"),
			Value:  ptrString(r.GetRequestReason()),
			Inline: ptrBool(false),
		})
	}

	embed := discordwebhook.Embed{
		Title:       ptrString("New Access Request"),
		Description: ptrString(fmt.Sprintf("Approve request by running command %s: ```tctl requests approve %s```", d.cfg.watcherList, r.GetName())),
		Color:       ptrString("3093206"),
		Fields:      &fields,
	}

	message := discordwebhook.Message{
		Embeds: &[]discordwebhook.Embed{embed},
	}

	err := d.sendWebhookNotification(message)
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
		return err
	}
	fmt.Println("Notification sent successfully!")

	return nil
}

func (d *discordClient) handleApproveAccessRequest(r types.AccessRequest) error {
	embed := discordwebhook.Embed{
		Title: ptrString("Access Request Approved"),
		Color: ptrString("2021216"),
		Fields: &[]discordwebhook.Field{
			{
				Name:   ptrString("Request ID"),
				Value:  ptrString(r.GetName()),
				Inline: ptrBool(false),
			},
			{
				Name:   ptrString("User"),
				Value:  ptrString(r.GetUser()),
				Inline: ptrBool(true),
			},
			{
				Name:   ptrString("Roles"),
				Value:  ptrString(strings.Join(r.GetRoles(), ", ")),
				Inline: ptrBool(true),
			},
		},
	}

	message := discordwebhook.Message{
		Embeds: &[]discordwebhook.Embed{embed},
	}

	err := d.sendWebhookNotification(message)
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
		return err
	}
	fmt.Println("Notification sent successfully!")

	return nil
}

func (d *discordClient) handleDenyAccessRequest(r types.AccessRequest) error {
	embed := discordwebhook.Embed{
		Title: ptrString("Access Request Denied"),
		Color: ptrString("15158332"),
		Fields: &[]discordwebhook.Field{
			{
				Name:   ptrString("Request ID"),
				Value:  ptrString(r.GetName()),
				Inline: ptrBool(false),
			},
			{
				Name:   ptrString("User"),
				Value:  ptrString(r.GetUser()),
				Inline: ptrBool(true),
			},
			{
				Name:   ptrString("Roles"),
				Value:  ptrString(strings.Join(r.GetRoles(), ", ")),
				Inline: ptrBool(true),
			},
			{
				Name:   ptrString("Reason"),
				Value:  ptrString(r.GetResolveReason()),
				Inline: ptrBool(false),
			},
		},
	}

	message := discordwebhook.Message{
		Embeds: &[]discordwebhook.Embed{embed},
	}

	err := d.sendWebhookNotification(message)
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
		return err
	}
	fmt.Println("Notification sent successfully!")

	return nil
}
