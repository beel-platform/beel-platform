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
  var project interface{}
  var toml_location string = "blue.toml"
  // Validate and Load TOML file
  if _, err := os.Stat(toml_location); os.IsNotExist(err) {
    msg_error("TOML file not found")
  }
  tree, err := toml.LoadFile(toml_location)
  if err != nil {
    fmt.Println("[TOML ERROR]", err.Error())
  } else {
    // Validate group project
    if !tree.Has("project") {
      msg_error("[CONFIG] Project group it's required")
    }
    g_project := tree.Get("project").(*toml.TomlTree)
    // Validate repository location
    if !g_project.Has("location") {
      msg_error("[CONFIG] Project's code location it's required")
    }
    location := g_project.Get("location").(string)
    if _, err := os.Stat(location); os.IsNotExist(err) {
      msg_error("[CONFIG] Project's code location not found")
    }
    // Validate schemas
    if tree.Has("schema") {
      schemas := tree.Get("schema").(*toml.TomlTree)
      for _, schema_key := range schemas.Keys() {
        schema := schemas.Get(schema_key).(*toml.TomlTree)
        fmt.Println(schema.Get("hostname"))
        // Validate SQL dump
        if schema.Has("sql_dump") {
          sqldump := schema.Get("sql_dump").(string)
          if _, err := os.Stat(sqldump); os.IsNotExist(err) {
            msg_error("[CONFIG] SQL dump not found")
          }
        }
      }
    }

    // Print out test result
    if g_project.Has("name") {
      project = g_project.Get("name").(string)
    }
    fmt.Println("The repo path for project:", project, "is", location)
  }
}

func main() {
  load_config()
}
