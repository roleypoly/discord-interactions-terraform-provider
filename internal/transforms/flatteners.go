package transforms

import (
	"github.com/roleypoly/terraform-provider-discord-interactions/internal/client"
)

func FlattenCommand(command *client.InteractionCommand) map[string]interface{} {
	if command == nil {
		return nil
	}

	commandItem := make(map[string]interface{})

	commandItem["name"] = command.Name
	commandItem["id"] = command.ID
	commandItem["description"] = command.Description
	commandItem["default_permission"] = command.DefaultPermission
	commandItem["option"] = FlattenOptions(command.Options)

	if command.GuildID != "" {
		commandItem["guild_id"] = command.GuildID
	}

	return commandItem
}

func FlattenOptions(options []client.InteractionCommandOption) []interface{} {
	items := make([]interface{}, len(options))

	for i, option := range options {
		optionItem := make(map[string]interface{})

		optionItem["type"] = option.Type
		optionItem["name"] = option.Name
		optionItem["description"] = option.Description
		optionItem["required"] = option.Required
		optionItem["choice"] = FlattenChoices(option.Choices)

		items[i] = optionItem
	}

	return items
}

func FlattenChoices(choices []client.InteractionCommandOptionChoice) []interface{} {
	items := make([]interface{}, len(choices))

	for i, choice := range choices {
		choiceItem := make(map[string]interface{})

		choiceItem["name"] = choice.Name

		switch value := choice.Value.(type) {
		case int, int32, int64:
			choiceItem["int_value"] = value
		case string:
			choiceItem["string_value"] = value
		case float32, float64:
			choiceItem["float_value"] = value
		}

		items[i] = choiceItem
	}

	return items
}
