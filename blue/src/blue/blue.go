package main

import (
  "io"
  "os"
  "fmt"
  "regexp"
  "strconv"
  "strings"
  "reflect"
  "net/http"
  "archive/tar"
  "path/filepath"
  "compress/gzip"
  "github.com/pelletier/go-toml"
)

func msg_error(message string) {
  fmt.Println("[ERROR]", message)
  os.Exit(1)
}

func msg_info(message string) {
  fmt.Println("[INFO]", message)
}

func download_file(filepath string, url string) (err error) {
  // Create the file
  out, err := os.Create(filepath)
  if err != nil  {
    return err
  }
  defer out.Close()
  // Get the data
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  // Writer the body to file
  _, err = io.Copy(out, resp.Body)
  if err != nil  {
    return err
  }
  return nil
}

func get_home() string {
  var home string
  for _, env_var := range os.Environ() {
    match, err := regexp.MatchString("HOME=(.*)", env_var)
    if (err != nil) {
      msg_error("Cannot determine environment variable HOME")
    } else if match {
      home = strings.Split(env_var,"=")[1]
      break
    }
  }
  return home
}

func overwrite(file_path string) (*os.File, error) {
  f, err := os.OpenFile(file_path, os.O_RDWR|os.O_TRUNC, 0777)
  if err != nil {
    f, err = os.Create(file_path)
    if err != nil {
      return f, err
    }
  }
  return f, nil
}

func read(file_path string) (*os.File, error) {
  f, err := os.OpenFile(file_path, os.O_RDONLY, 0444)
  if err != nil {
    return f, err
  }
  return f, nil
}

func untar(tar_path string, dest_path string) {
  fr, err := read(tar_path)
  defer fr.Close()
  if err != nil {
    panic(err)
  }
  gr, err := gzip.NewReader(fr)
  defer gr.Close()
  if err != nil {
    panic(err)
  }
  tr := tar.NewReader(gr)
  for {
    hdr, err := tr.Next()
    if err == io.EOF {
      // end of tar archive
      break
    }
    if err != nil {
      panic(err)
    }
    path := dest_path+"/"+hdr.Name
    switch hdr.Typeflag {
    case tar.TypeDir:
      if err := os.MkdirAll(path, os.FileMode(hdr.Mode)); err != nil {
        panic(err)
      }
    case tar.TypeReg:
      ow, err := overwrite(path)
      defer ow.Close()
      if err != nil {
        panic(err)
      }
      if _, err := io.Copy(ow, tr); err != nil {
        panic(err)
      }
    default:
      fmt.Printf("Can't: %c, %s\n", hdr.Typeflag, path)
    }
  }
}

func hab_install(tree *toml.TomlTree) {
  if !tree.Has("habitat.bin_path") {
    msg_error("[CONF] Group habitat it's required")
  }
  location := tree.Get("habitat.bin_path").(string)
  if _, err := os.Stat(location); os.IsNotExist(err) {
    // Download habitat
    if !tree.Has("habitat.download_url") {
      msg_error("[CONF] Habitat download URL it's required")
    }
    download_url := tree.Get("habitat.download_url").(string)
    hab_path := filepath.Join(get_home(), "/", location)
    hab_bin := filepath.Join(hab_path, "/bin")
    hab_tmp := filepath.Join(hab_path, "/tmp")
    if _, err := os.Stat(hab_bin); os.IsNotExist(err) {
      os.MkdirAll(hab_bin, 0755)
    }
    if _, err := os.Stat(hab_tmp); os.IsNotExist(err) {
      os.MkdirAll(hab_tmp, 0755)
    }
    tar_path := filepath.Join(hab_tmp, "/habitat.tar.gz")
    if _, err := os.Stat(tar_path); os.IsNotExist(err) {
      msg_info("Downloading habitat")
      download_file(tar_path, download_url)
    }
    // Decompress habitat's tarball
    msg_info("Decompressig habitat")
    untar(tar_path, hab_tmp)
    if err := os.Remove(tar_path); err != nil {
      msg_info("Cannot remove "+tar_path)
      fmt.Println(err)
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
}

func load_config() {
  // Declare variables
  // var project interface{}
  var toml_location string = "blue.toml"
  // Validate and Load TOML file
  if _, err := os.Stat(toml_location); os.IsNotExist(err) {
    msg_error("TOML file not found")
  }
  tree, err := toml.LoadFile(toml_location)
  if err != nil {
    msg_error(err.Error())
  } else {
    // Validate group project
    if !tree.Has("project") {
      msg_error("[CONF] Group project it's required")
    }
    g_project := tree.Get("project").(*toml.TomlTree)
    // Validate repository location
    if !g_project.Has("location") {
      msg_error("[CONF] Project's code location it's required")
    }
    location := g_project.Get("location").(string)
    if _, err := os.Stat(location); os.IsNotExist(err) {
      msg_error("[CONF] Project dir "+ location +" not found")
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
            msg_error("[CONF] File ("+ sqldump +") not found")
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
    hab_install(tree)
    // Iterate packages
    for _,item := range results.Values() {
      fmt.Println(item)
    }
  }
}

func main() {
  load_config()
}
