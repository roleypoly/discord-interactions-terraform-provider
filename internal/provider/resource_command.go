package provider

import (
	"context"

	"github.com/roleypoly/terraform-provider-discord-interactions/internal/client"
	"github.com/roleypoly/terraform-provider-discord-interactions/internal/transforms"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlobalCommand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCommandCreate,
		ReadContext:   resourceCommandRead,
		UpdateContext: resourceCommandUpdate,
		DeleteContext: resourceCommandDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "1-32 lowercase character name. Changing this will force recreation.",
				Required:     true,
				ValidateFunc: transforms.ValidateName,
				ForceNew:     true,
			},
			"description": {
				Type:         schema.TypeString,
				Description:  "1-100 character description",
				Required:     true,
				ValidateFunc: transforms.ValidateDescription,
			},
			"default_permission": {
				Type:        schema.TypeBool,
				Description: "whether the command is enabled by default when the app is added to a guild",
				Optional:    true,
				Default:     true,
			},
			"option": optionsSchema(true),
		},
	}
}

func optionsSchema(allowNesting bool) *schema.Schema {
	options := &schema.Schema{
		Type:        schema.TypeList,
		Description: "Parameters for the command. Note: As an implementation detail, making a subcommand group is supported, but only exactly one nesting down of subcommands.",
		Optional:    true,
		MaxItems:    25,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Type of the option, refer to documentation: https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type",
					Default:     3,
				},
				"name": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: transforms.ValidateName,
				},
				"description": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: transforms.ValidateDescription,
				},
				"required": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"choice": {
					Type:     schema.TypeList,
					MaxItems: 25,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:         schema.TypeString,
								Description:  "1-100 character choice name",
								ValidateFunc: transforms.ValidateDescription,
								Required:     true,
							},
							"string_value": {
								Type:         schema.TypeString,
								ValidateFunc: transforms.ValidateDescription,
								Optional:     true,
							},
							"int_value": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"float_value": {
								Type:     schema.TypeFloat,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}

	if allowNesting {
		options.Elem.(*schema.Resource).Schema["option"] = optionsSchema(false)
	}

	return options
}

// resourceGuildCommand adds guild_id to resourceGlobalCommand as that's the only difference
func resourceGuildCommand() *schema.Resource {
	resource := resourceGlobalCommand()

	resource.Schema["guild_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: transforms.ValidateSnowflake,
	}

	return resource
}

func resourceCommandCreate(ctx context.Context, resource *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.InteractionsClient)
	guildID, _ := resource.Get("guild_id").(string)

	command := &client.InteractionCommand{
		Name:              resource.Get("name").(string),
		Description:       resource.Get("description").(string),
		DefaultPermission: resource.Get("default_permission").(bool),
		Options:           transforms.ExpandOptions(resource.Get("option").([]interface{})),
	}

	command, err := c.UpsertInteractionCommand(guildID, command)
	if err != nil {
		return diag.FromErr(err)
	}

	resource.SetId(command.ID)

	return resourceCommandRead(ctx, resource, m)
}

func resourceCommandRead(ctx context.Context, resource *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*client.InteractionsClient)
	guildID, _ := resource.Get("guild_id").(string)

	command, err := c.GetInteractionCommand(guildID, resource.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	commandItem := transforms.FlattenCommand(command)
	for key, value := range commandItem {
		err := resource.Set(key, value)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceCommandUpdate(ctx context.Context, resource *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.InteractionsClient)
	guildID, _ := resource.Get("guild_id").(string)

	command := &client.InteractionCommand{
		ID:                resource.Id(),
		Name:              resource.Get("name").(string),
		Description:       resource.Get("description").(string),
		DefaultPermission: resource.Get("default_permission").(bool),
		Options:           transforms.ExpandOptions(resource.Get("option").([]interface{})),
	}

	command, err := c.UpsertInteractionCommand(guildID, command)
	if err != nil {
		return diag.FromErr(err)
	}

	resource.SetId(command.ID)

	return resourceCommandRead(ctx, resource, m)
}

func resourceCommandDelete(ctx context.Context, resource *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*client.InteractionsClient)
	guildID, _ := resource.Get("guild_id").(string)

	err := c.DeleteInteractionCommand(guildID, resource.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
