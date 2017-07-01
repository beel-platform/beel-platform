package main

import (
  "fmt"
  "context"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
)

func main() {
  cli, err := client.NewEnvClient()
  if err != nil {
    panic(err)
  }
  // searchResult, err := cli.ImageSearch(context.Background(), "bluespark", types.ImageSearchOptions{Limit:10})
  // if err != nil {
  //   fmt.Printf("%s", err.Error())
  // }
  // fmt.Printf("%s", searchResult)
  // os.Exit(0)

  var Dockerfilefoo string = "/Users/bbh/Projects/demo/docker/myproject/"
  options := types.ImageBuildOptions{Dockerfile: Dockerfilefoo}
  buildResponse, err := cli.ImageBuild(context.Background(), nil, options)
  if err != nil {
    fmt.Printf("%s", err.Error())
  }
  fmt.Printf("%s", buildResponse.OSType)
}
