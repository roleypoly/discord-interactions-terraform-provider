package transforms

import "github.com/roleypoly/discord-interactions-terraform-provider/internal/client"

func ExpandOptions(optionItems []interface{}) []client.InteractionCommandOption {
	options := make([]client.InteractionCommandOption, len(optionItems))

	for i, itemIntf := range optionItems {
		item := itemIntf.(map[string]interface{})
		option := client.InteractionCommandOption{
			Type:        item["type"].(int),
			Name:        item["name"].(string),
			Description: item["description"].(string),
			Required:    item["required"].(bool),
			Choices:     ExpandChoices(item["choice"].([]interface{})),
		}

		options[i] = option
	}

	// TODO: validate nesting?

	return SortRequiredOptions(options)
}

func SortRequiredOptions(options []client.InteractionCommandOption) []client.InteractionCommandOption {
	required := []client.InteractionCommandOption{}
	notRequired := []client.InteractionCommandOption{}

	for _, option := range options {
		if option.Required {
			required = append(required, option)
		} else {
			notRequired = append(notRequired, option)
		}
	}

	return append(required, notRequired...)
}

func ExpandChoices(choiceItems []interface{}) []client.InteractionCommandOptionChoice {
	choices := make([]client.InteractionCommandOptionChoice, len(choiceItems))

	for i, itemIntf := range choiceItems {
		item := itemIntf.(map[string]interface{})
		choice := client.InteractionCommandOptionChoice{
			Name: item["name"].(string),
		}

		if item["string_value"] != nil {
			choice.Value = item["string_value"].(string)
		}

		if item["int_value"] != nil {
			choice.Value = item["int_value"].(int)
		}

		if item["float_value"] != nil {
			choice.Value = item["float_value"].(float64)
		}

		choices[i] = choice
	}

	return choices
}
