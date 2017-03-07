package main

import (
  "fmt"
  "os"
  "github.com/pelletier/go-toml"
)

func msg_error(error_message string) {
  fmt.Println("[ERROR]", error_message)
  os.Exit(1)
}

func load_config() {
  // Declare variables
  var toml_location string = "blue.toml"
  var repository string = ""
  var project interface{}
  // Validate TOML file location
  if _, err := os.Stat(toml_location); os.IsNotExist(err) {
    msg_error("TOML file not found")
  }
  // Load TOML
  tree, err := toml.LoadFile("blue.toml")
  if err != nil {
    fmt.Println("ERROR ", err.Error())
  } else {
    // Validate blue group
    if !tree.Has("blue") {
      msg_error("[CONFIG] Group blue it's required")
    }
    config := tree.Get("blue").(*toml.TomlTree)
    // Validate repository
    if !config.Has("repo") {
      msg_error("[CONFIG] Repository path it's required")
    }
    repository = config.Get("repo").(string)
    if _, err := os.Stat(repository); os.IsNotExist(err) {
      msg_error("[CONFIG] Repository path not found")
    }

    // Print out test result
    if config.Has("name") {
      project = config.Get("name").(string)
    }
    fmt.Println("The repo path for project:", project, "is", repository)
  }
}

func main() {
  load_config()
}
