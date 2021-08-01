resource "discord-interactions_guild_command" "example" {
  name        = "hello-world"
  description = "An example global command"
  guild_id    = "386659935687147521"

  option {
    type        = 6
    name        = "user"
    description = "Tell this person hello!"
  }

  option {
    type        = 3
    name        = "message"
    description = "What message do I send?"
  }
}
