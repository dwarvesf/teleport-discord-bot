package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/dwarvesf/teleport-discord-bot/internal/config"
	repo "github.com/dwarvesf/teleport-discord-bot/internal/repository"

	"github.com/gravitational/teleport/api/types"
	"github.com/gtuk/discordwebhook"
)

// Client handles Discord webhook notifications
type Client struct {
	url string
	cfg *config.Config
}

// NewClient creates a new Discord webhook client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		url: cfg.DiscordWebhookURL,
		cfg: cfg,
	}
}

// sendWebhookNotification sends a message to the configured Discord webhook
func (d *Client) sendWebhookNotification(message discordwebhook.Message) error {
	return discordwebhook.SendMessage(d.url, message)
}

// ptrString is a helper to convert string to pointer
func ptrString(s string) *string {
	return &s
}

// ptrBool is a helper to convert bool to pointer
func ptrBool(b bool) *bool {
	return &b
}

// HandleNewAccessRequest creates a notification for a new access request
func (d *Client) HandleNewAccessRequest(r types.AccessRequest) error {
	approverIDs := ""
	approvers, err := repo.GetApprovers()
	if err != nil {
		fmt.Printf("Failed to get approvers: %v\n", err)
	} else {
		for _, a := range approvers {
			if a.TlpUsername != r.GetUser() {
				approverIDs += fmt.Sprintf("<@%s> ", a.DiscordID)
			}
		}
	}

	requester := r.GetUser()

	dbRequester, err := repo.GetUserByTlpUsername(requester)
	if err == nil {
		requester = fmt.Sprintf("<@%s>", dbRequester.DiscordID)
	}

	fields := []discordwebhook.Field{
		{
			Name:   ptrString("Request ID"),
			Value:  ptrString(r.GetName()),
			Inline: ptrBool(false),
		},
		{
			Name:   ptrString("Requester"),
			Value:  ptrString(requester),
			Inline: ptrBool(true),
		},
		{
			Name:   ptrString("Roles"),
			Value:  ptrString(strings.Join(r.GetRoles(), ", ")),
			Inline: ptrBool(true),
		},
		{
			Name:   ptrString("Duration"),
			Value:  ptrString(r.GetSessionTLL().Sub(r.GetCreationTime()).Round(time.Second).String()),
			Inline: ptrBool(true),
		},
	}

	if r.GetRequestReason() != "" {
		fields = append(fields, discordwebhook.Field{
			Name:   ptrString("Reason"),
			Value:  ptrString(r.GetRequestReason()),
			Inline: ptrBool(false),
		})
	}

	embed := discordwebhook.Embed{
		Title:       ptrString("New access request"),
		Description: ptrString(fmt.Sprintf("Approve request by running command: ```tctl requests approve %s```", r.GetName())),
		Color:       ptrString("3093206"),
		Fields:      &fields,
	}

	message := discordwebhook.Message{
		Embeds:  &[]discordwebhook.Embed{embed},
		Content: ptrString(approverIDs),
	}

	if err := d.sendWebhookNotification(message); err != nil {
		return fmt.Errorf("failed to send new access request notification: %w", err)
	}

	return nil
}

// HandleApproveAccessRequest creates a notification for an approved access request
func (d *Client) HandleApproveAccessRequest(r types.AccessRequest) error {
	requester := r.GetUser()

	dbRequester, err := repo.GetUserByTlpUsername(requester)
	if err == nil {
		requester = fmt.Sprintf("<@%s>", dbRequester.DiscordID)
	}

	embed := discordwebhook.Embed{
		Title:       ptrString("Access request approved"),
		Description: ptrString(fmt.Sprintf("Locking user by running command: ```tctl lock --user=%s --message=\"block reason\" --ttl=1h```", r.GetUser())),
		Color:       ptrString("2021216"),
		Fields: &[]discordwebhook.Field{
			{
				Name:   ptrString("Request ID"),
				Value:  ptrString(r.GetName()),
				Inline: ptrBool(false),
			},
			{
				Name:   ptrString("Requester"),
				Value:  ptrString(requester),
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

	if err := d.sendWebhookNotification(message); err != nil {
		return fmt.Errorf("failed to send access request approval notification: %w", err)
	}

	return nil
}

// HandleDenyAccessRequest creates a notification for a denied access request
func (d *Client) HandleDenyAccessRequest(r types.AccessRequest) error {
	requester := r.GetUser()

	dbRequester, err := repo.GetUserByTlpUsername(requester)
	if err == nil {
		requester = fmt.Sprintf("<@%s>", dbRequester.DiscordID)
	}

	embed := discordwebhook.Embed{
		Title: ptrString("Access request denied"),
		Color: ptrString("15158332"),
		Fields: &[]discordwebhook.Field{
			{
				Name:   ptrString("Request ID"),
				Value:  ptrString(r.GetName()),
				Inline: ptrBool(false),
			},
			{
				Name:   ptrString("Requester"),
				Value:  ptrString(requester),
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

	if err := d.sendWebhookNotification(message); err != nil {
		return fmt.Errorf("failed to send access request denial notification: %w", err)
	}

	return nil
}
