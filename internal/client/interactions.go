package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (i *InteractionsClient) GetInteractionCommands(guildID string) ([]*InteractionCommand, error) {
	url := `/commands`
	if guildID != "" {
		url = `/guilds/` + guildID + url
	}

	response, err := i.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("GET call to %s failed, %w", url, err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("response body unavilable, %w", err)
	}

	if response.StatusCode != 200 {
		return nil, i.ErrFromResponse(response)
	}

	commands := []*InteractionCommand{}
	err = json.Unmarshal(body, &commands)
	if err != nil {
		return nil, fmt.Errorf("JSON parse issue for %s: %w", url, i.ErrFromResponse(response))
	}
	return commands, nil
}

func (i *InteractionsClient) GetInteractionCommand(guildID string, commandID string) (*InteractionCommand, error) {
	url := `/commands/` + commandID
	if guildID != "" {
		url = `/guilds/` + guildID + url
	}

	response, err := i.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("GET call to %s failed, %w", url, err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("response body unavilable, %w", err)
	}

	if response.StatusCode != 200 {
		return nil, i.ErrFromResponse(response)
	}

	command := &InteractionCommand{}
	err = json.Unmarshal(body, command)
	if err != nil {
		return nil, fmt.Errorf("JSON parse issue for %s: %w", url, i.ErrFromResponse(response))
	}
	return command, nil
}

func (i *InteractionsClient) UpsertInteractionCommand(guildID string, command *InteractionCommand) (*InteractionCommand, error) {
	url := `/commands`
	if guildID != "" {
		url = `/guilds/` + guildID + url
	}
	// if command.ID != "" {
	// 	url = url + `/` + command.ID
	// }

	response, err := i.makeRequest("POST", url, command)
	if err != nil {
		return nil, fmt.Errorf("POST call to %s failed: %w", url, err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("response body unavilable, %w", err)
	}

	if response.StatusCode != 201 && response.StatusCode != 200 {
		return nil, i.ErrFromResponse(response, body)
	}

	commandResponse := &InteractionCommand{}
	err = json.Unmarshal(body, commandResponse)
	if err != nil {
		return nil, fmt.Errorf("JSON parse issue for %s: %w", url, i.ErrFromResponse(response))
	}
	return commandResponse, nil
}

func (i *InteractionsClient) DeleteInteractionCommand(guildID string, commandID string) error {
	url := `/commands/` + commandID
	if guildID != "" {
		url = `/guilds/` + guildID + url
	}

	response, err := i.makeRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("DELETE call to %s failed, %w", url, err)
	}

	if response.StatusCode != 204 {
		return i.ErrFromResponse(response)
	}

	return err
}
