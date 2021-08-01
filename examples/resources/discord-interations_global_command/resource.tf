resource "discord-interactions_global_command" "example" {
  name        = "hello-world"
  description = "An example guild-specific command"

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
