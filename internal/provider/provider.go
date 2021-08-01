package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/roleypoly/terraform-provider-discord-interactions/internal/client"
	"github.com/roleypoly/terraform-provider-discord-interactions/internal/transforms"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"discord-interactions_global_command": resourceGlobalCommand(),
				"discord-interactions_guild_command":  resourceGuildCommand(),
				// TODO: add permissions
			},
			DataSourcesMap: map[string]*schema.Resource{},
			Schema: map[string]*schema.Schema{
				"application_id": {
					Type:         schema.TypeString,
					Description:  "Discord Application ID from https://discord.com/developers",
					Required:     true,
					DefaultFunc:  schema.EnvDefaultFunc("DISCORD_APPLICATION_ID", nil),
					ValidateFunc: transforms.ValidateSnowflake,
				},
				"bot_token": {
					Type:         schema.TypeString,
					Sensitive:    true,
					Optional:     true,
					ExactlyOneOf: []string{"client_credentials_token", "bot_token"},
					Description:  "Discord bot token from https://discord.com/developers",
					DefaultFunc:  schema.EnvDefaultFunc("DISCORD_BOT_TOKEN", nil),
				},
				"client_credentials_token": {
					Type:         schema.TypeString,
					Sensitive:    true,
					Optional:     true,
					ExactlyOneOf: []string{"client_credentials_token", "bot_token"},
					Description:  "Discord client credentials token. (must have applications.commands.update scope)",
					DefaultFunc:  schema.EnvDefaultFunc("DISCORD_CLIENT_TOKEN", nil),
				},
				"api_root": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Discord API root. Only useful for testing or version swaps, don't use this otherwise.",
					DefaultFunc: schema.EnvDefaultFunc("DISCORD_API_ROOT", "https://discord.com/api/v9"),
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						value := val.(string)

						if strings.HasSuffix(value, "/") {
							errs = append(errs, fmt.Errorf("api_root cannot end in /, got: `%s`", value))
						}

						if value != "https://discord.com/api/v9" {
							warns = append(warns, "api_root is not default, this is extremely unsupported behavior")
						}

						return
					},
				},
			},
		}

		p.ConfigureContextFunc = configureProvider(version, p)

		return p
	}
}

func configureProvider(version string, provider *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		userAgent := provider.UserAgent("terraform-provider-discord-interactions", version)

		var diags diag.Diagnostics

		applicationID := d.Get("application_id").(string)
		botToken := d.Get("bot_token").(string)
		clientCredentials := d.Get("client_credentials_token").(string)
		apiRoot := d.Get("api_root").(string)

		client, err := client.NewInteractionsClient(client.ClientConfig{
			ApplicationID:     applicationID,
			BotToken:          botToken,
			ClientCredentials: clientCredentials,
			APIRoot:           apiRoot,
			UserAgent:         userAgent,
		})

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}
}
