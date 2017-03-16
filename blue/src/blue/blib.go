package main

import (
  "io"
  "os"
  "fmt"
  "os/exec"
  "strings"
  "net/http"
  "archive/tar"
  "archive/zip"
  "path/filepath"
  "compress/gzip"
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
  defer out.Close()
  if err != nil  {
    return err
  }
  // Get the data
  resp, err := http.Get(url)
  defer resp.Body.Close()
  if err != nil {
    return err
  }
  // Writer the body to file
  _, err = io.Copy(out, resp.Body)
  if err != nil  {
    return err
  }
  return nil
}

func get_abs(path string) string {
  var result string
  if strings.Contains(path, "~") {
    result = filepath.Join(get_home(), strings.Replace(path, "~", "", 1))
  } else {
    result = path
  }
  return result
}

func get_home() string {
  home := os.Getenv("HOME")
  if home == "" {
    msg_error("Cannot find HOME environment variable")
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

func unzip(src, dest string) error {
  r, err := zip.OpenReader(src)
  defer r.Close()
  if err != nil {
    return err
  }
  for _, f := range r.File {
    rc, err := f.Open()
    defer rc.Close()
    if err != nil {
      return err
    }
    fpath := filepath.Join(dest, f.Name)
    if f.FileInfo().IsDir() {
      os.MkdirAll(fpath, f.Mode())
    } else {
      var fdir string
      if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
        fdir = fpath[:lastIndex]
      }
      err = os.MkdirAll(fdir, f.Mode())
      if err != nil {
        return err
      }
      f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
      defer f.Close()
      if err != nil {
        return err
      }
      _, err = io.Copy(f, rc)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func printCommand(cmd *exec.Cmd) {
  fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
  if err != nil {
    os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
  }
}

func printOutput(outs []byte) {
  if len(outs) > 0 {
    fmt.Printf("==> Output: %s\n", string(outs))
  }
}
