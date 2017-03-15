package main

import (
  "os"
  "fmt"
  "runtime"
  "strconv"
  "reflect"
  "path/filepath"
  "github.com/pelletier/go-toml"
)

func hab_install(tree *toml.TomlTree) string {
  // Validate bsp_path
  if !tree.Has("habitat.bsp_path") {
    msg_error("[CONF] Group habitat it's required")
  }
  bsp_path := filepath.Join(get_home(), "/", tree.Get("habitat.bsp_path").(string))
  if _, err := os.Stat(bsp_path); os.IsNotExist(err) {
    // Prepare structure
    hab_bin := filepath.Join(bsp_path, "/bin")
    hab_tmp := filepath.Join(bsp_path, "/tmp")
    if _, err := os.Stat(hab_bin); os.IsNotExist(err) {
      os.MkdirAll(hab_bin, 0755)
    }
    if _, err := os.Stat(hab_tmp); os.IsNotExist(err) {
      os.MkdirAll(hab_tmp, 0755)
    }
    // Download habitat
    switch ostype := runtime.GOOS; ostype {
    case "darwin":
      if !tree.Has("habitat.download_url_macos") {
        msg_error("[CONF] Habitat download URL it's required")
      }
      download_url := tree.Get("habitat.download_url_macos").(string)
      zip_path := filepath.Join(hab_tmp, "/habitat.zip")
      if _, err := os.Stat(zip_path); os.IsNotExist(err) {
        msg_info("Downloading habitat")
        download_file(zip_path, download_url)
      }
      // Decompress tarball
      msg_info("Decompressig habitat")
      unzip(zip_path, hab_tmp)
      if err := os.Remove(zip_path); err != nil {
        msg_info("Cannot remove "+zip_path)
        fmt.Println(err)
      }
    case "linux":
      if !tree.Has("habitat.download_url_linux") {
        msg_error("[CONF] Habitat download URL it's required")
      }
      download_url := tree.Get("habitat.download_url_linux").(string)
      tar_path := filepath.Join(hab_tmp, "/habitat.tar.gz")
      if _, err := os.Stat(tar_path); os.IsNotExist(err) {
        msg_info("Downloading habitat")
        download_file(tar_path, download_url)
      }
      // Decompress tarball
      msg_info("Decompressig habitat")
      untar(tar_path, hab_tmp)
      if err := os.Remove(tar_path); err != nil {
        msg_info("Cannot remove "+tar_path)
        fmt.Println(err)
      }
    default: // freebsd, openbsd, windows...
      msg_error("Operating System: "+ostype+" not supported")
    }
    // Move, change permissions and delete temporal
    if tmp_dir,err := filepath.Glob(filepath.Join(hab_tmp, "/hab-*")); err != nil {
      fmt.Println(err)
    } else {
      msg_info("Installing habitat")
      if err := os.Rename(filepath.Join(tmp_dir[0],"/hab"), filepath.Join(hab_bin,"/hab")); err != nil {
        msg_error("Cannot move habitat binary to "+hab_bin)
      }
      if err := os.Chmod(filepath.Join(hab_bin,"/hab"),0755); err != nil {
        msg_error("Cannot change habitat's permissions")
      }
      if err := os.Remove(tmp_dir[0]); err != nil {
        msg_info("Cannot remove temporal directory "+tmp_dir[0])
        fmt.Println(err)
      }
    }
  }
  return filepath.Join(bsp_path, "/bin/hab")
}

func load_config() {
  var toml_location string = "blue.toml"
  // Validate and Load TOML file
  if _, err := os.Stat(toml_location); os.IsNotExist(err) {
    msg_error("TOML file not found")
  }
  if tree, err := toml.LoadFile(toml_location); err != nil {
    msg_error(err.Error())
  } else {
    // Validate group project
    if !tree.Has("project") {
      msg_error("[CONF] Group project it's required")
    }
    g_project := tree.Get("project").(*toml.TomlTree)
    // Validate repository location
    if !g_project.Has("path") {
      msg_error("[CONF] Project's code location it's required")
    }
    location := get_abs(g_project.Get("path").(string))
    if _, err := os.Stat(location); os.IsNotExist(err) {
      msg_error("[CONF] Directory "+ location +" not found")
    }
    // Validate schemas
    if tree.Has("schema") {
      schemas := tree.Get("schema").(*toml.TomlTree)
      for _, schema_key := range schemas.Keys() {
        schema := schemas.Get(schema_key).(*toml.TomlTree)
        // Validate SQL dump
        if schema.Has("sql_dump") {
          sqldump := schema.Get("sql_dump").(string)
          if _, err := os.Stat(sqldump); os.IsNotExist(err) {
            sqldump = filepath.Join(location, sqldump)
            if _, err := os.Stat(sqldump); os.IsNotExist(err) {
              msg_error("[CONF] File ("+ sqldump +") not found")
            }
          }
        }
      }
    }
    // Validate stack
    if !tree.Has("stack") {
      msg_error("[CONF] Group stack it's required")
    }
    // Validate packages
    // TODO Validate if array it's empty
    if !tree.Has("stack.packages") {
      msg_error("[CONF] Stack's packages required")
    }
    rt := reflect.TypeOf(tree.Get("stack.packages"))
    if rt.Kind() != reflect.Slice {
      msg_error("[CONF] Stack's packages attribute should be an array")
    }
    // Load primary results to add extra slice (workaround)
    pres,_ := tree.Query("$.stack.packages[0:-1]")
    items := len(pres.Values())+1 // Workaround to get all the nodes
    results,_ := tree.Query("$.stack.packages[0:"+strconv.Itoa(items)+"]")
    // Validate habitat and install
    hab_bin := hab_install(tree)
    fmt.Println(hab_bin)
    // Iterate packages
    for _,item := range results.Values() {
      fmt.Println(item)
    }
  }
}

func main() {
  load_config()
}
