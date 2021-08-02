package provider

import (
	"log"
	"os"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/roleypoly/terraform-provider-discord-interactions/internal/client"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"discord-interactions": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

var (
	testGuildID          = os.Getenv("TEST_GUILD_ID")
	discordBotToken      = os.Getenv("DISCORD_BOT_TOKEN")
	discordApplicationID = os.Getenv("DISCORD_APPLICATION_ID")
)

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	if discordBotToken == "" {
		t.Fatalf("environment variable DISCORD_BOT_TOKEN is not set")
	}

	if discordApplicationID == "" {
		t.Fatalf("environment variable DISCORD_APPLICATION_ID is not set")
	}

	if testGuildID == "" {
		t.Fatalf("environment variable TEST_GUILD_ID is not set")
	}
}

func getClient() *client.InteractionsClient {
	c, err := client.NewInteractionsClient(client.ClientConfig{
		ApplicationID: discordApplicationID,
		BotToken:      discordBotToken,
		APIRoot:       "https://discord.com/api/v9",
	})

	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	return c
}

func getName() string {
	petname.NonDeterministicMode()
	return petname.Name()
}
