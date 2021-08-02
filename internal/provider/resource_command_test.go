package provider

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/roleypoly/terraform-provider-discord-interactions/internal/client"
)

func init() {
	resource.AddTestSweepers("discord-interactions_guild_command", &resource.Sweeper{
		Name: "discord-interactions_guild_command",
		F: func(guildID string) error {
			c := getClient()

			commands, err := c.GetInteractionCommands(guildID)
			if err != nil {
				return fmt.Errorf("error getting commands: %w", err)
			}
			for _, command := range commands {
				if strings.HasPrefix(command.Name, "test-acc") {
					err := c.DeleteInteractionCommand(guildID, command.ID)

					if err != nil {
						log.Printf("Error destroying %s during sweep: %v", command.Name, err)
					}
				}
			}
			return nil
		},
	})
}

func TestAccDiscordInteractionsGuildCommand_basic(t *testing.T) {
	path := "discord-interactions_guild_command.hello-world"
	name := "test-acc-hello-world"
	description := "A test command for terraform acceptance tests"

	var command client.InteractionCommand
	knownCommand := &client.InteractionCommand{
		Name:        name,
		Description: description,
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccGuildCommandDestroy(testGuildID, name),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGuildCommand(testGuildID, name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccGuildCommandCreated(path, testGuildID, &command),
					testAccCommandMatches(&command, knownCommand),
					resource.TestCheckResourceAttrSet(path, "id"),
					resource.TestCheckResourceAttr(path, "name", name),
					resource.TestCheckResourceAttr(path, "description", description),
				),
			},
		},
	})
}

func testAccGuildCommandCreated(path, guildID string, commandOut *client.InteractionCommand) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[path]
		id := r.Primary.Attributes["id"]

		c := getClient()

		command, err := c.GetInteractionCommand(guildID, id)
		if err != nil {
			return fmt.Errorf("failed to get command, %w", err)
		}

		if command == nil {
			return fmt.Errorf("command was nil after creation check")
		}

		*commandOut = *command

		return nil
	}
}

func testAccCommandMatches(incoming, control *client.InteractionCommand) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if incoming.Name != control.Name {
			return fmt.Errorf("name doesn't match, got: %s, wanted: %s", incoming.Name, control.Name)
		}

		if incoming.Description != control.Description {
			return fmt.Errorf("description doesn't match, got: %s, wanted: %s", incoming.Description, control.Description)
		}

		return nil
	}
}

func testAccGuildCommandDestroy(guildID, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := getClient()

		commands, err := c.GetInteractionCommands(guildID)
		if err != nil {
			return fmt.Errorf("failed to get all commands, %w", err)
		}

		for _, command := range commands {
			if command.Name == name {
				return fmt.Errorf("command still exists")
			}
		}

		return nil
	}
}

func testAccResourceGuildCommand(guildID, name, description string) string {
	return fmt.Sprintf(`
	resource "discord-interactions_guild_command" "hello-world" {
		guild_id = "%s"
		name = "%s"
		description = "%s"
	}
	`, guildID, name, description)
}
