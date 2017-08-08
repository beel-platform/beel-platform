package main

import (
  "os"
  "fmt"
  "regexp"
  "syscall"
  "runtime"
  "strconv"
  "reflect"
  "os/exec"
  "path/filepath"
  "github.com/pelletier/go-toml"
  "github.com/pelletier/go-toml/query"
)

func hab_install(tree *toml.Tree) string {
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

func hab_pkg_install(hab_pkg string, hab_bin string) {
  msg_info("Installing habitat package: "+hab_pkg)
  // Prepare command
  cmd := exec.Command("/usr/bin/sudo", hab_bin, "pkg", "install", hab_pkg)
  // Execute command
  printCommand(cmd)
  var waitStatus syscall.WaitStatus
  if err := cmd.Run(); err != nil {
    printError(err)
    // Did the command fail because of an unsuccessful exit code
    if exitError, ok := err.(*exec.ExitError); ok {
      waitStatus = exitError.Sys().(syscall.WaitStatus)
      printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
    }
  } else {
    // Command was successful
    waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
    printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
  }
}

func load_config(toml_file string) {
  // Load TOML tree
  if tree, err := toml.LoadFile(toml_file); err != nil {
    msg_error(err.Error())
  } else {
    // Validate group project
    if !tree.Has("project") {
      msg_error("[CONF] Group project it's required")
    }
    g_project := tree.Get("project").(*toml.Tree)
    // Validate repository location
    if !g_project.Has("path") {
      msg_error("[CONF] Project's code location it's required")
    }
    location := get_abs(g_project.Get("path").(string))
    if _, err := os.Stat(location); os.IsNotExist(err) {
      msg_error("[CONF] Directory "+ location +" not found")
    }
    // Validate schemas
    if tree.Has("database") {
      schemas := tree.Get("database").(*toml.Tree)
      for _, schema_key := range schemas.Keys() {
        schema := schemas.Get(schema_key).(*toml.Tree)
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
    // Validate habitat packages
    if !tree.Has("habitat") {
      msg_error("[CONF] Group stack it's required")
    }
    // TODO Validate if array it's empty
    if !tree.Has("habitat.packages") {
      msg_error("[CONF] Stack's packages required")
    }
    rt := reflect.TypeOf(tree.Get("habitat.packages"))
    if rt.Kind() != reflect.Slice {
      msg_error("[CONF] Habitat packages must be an array")
    }
    // Load primary results to add extra slice (workaround)
    // pres,_ := tree.Query("$.habitat.packages[0:-1]")
    pres,_ := query.CompileAndExecute("$.habitat.packages[0:-1]", tree)
    items := len(pres.Values())+1 // Workaround to get all the nodes
    // results,_ := tree.Query("$.habitat.packages[0:"+strconv.Itoa(items)+"]")
    results,_ := query.CompileAndExecute("$.habitat.packages[0:"+strconv.Itoa(items)+"]", tree)
    // Validate habitat and install
    hab_bin := hab_install(tree)
    // Iterate packages
    for _,hab_pkg := range results.Values() {
      hab_pkg_install(hab_pkg.(string), hab_bin)
    }
  }
}

func main() {
  var toml_file string = "blue.toml"
  // Validate argument
  if len(os.Args) > 1 {
    args := os.Args
    matched, _ := regexp.MatchString("[a-z0-9].toml", args[1])
    if matched {
      toml_file = args[1]
    }
  }
  // Validate and Load TOML file
  if _, err := os.Stat(toml_file); os.IsNotExist(err) {
    msg_error("TOML file: "+toml_file+" not found")
  }
  // Load TOML configuration
  load_config(toml_file)
}
