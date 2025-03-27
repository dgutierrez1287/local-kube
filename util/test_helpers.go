package util

import (
  "os"
  "bufio"
)

/*
  Helper functions for testing, these will be
  used in many places
*/

var MockAppDir = "./mock/appdir"
var MockAnsibleRoleDir = "./mock/appdir/ansible-roles"

func MockAppDirSetup() error {
  err := os.MkdirAll(MockAppDir, 0755)
  if err != nil {
    return err
  }

  err = os.Mkdir(MockAnsibleRoleDir, 0755)
  if err != nil {
    return err
  }
  return nil
}

func MockAppDirCleanup() error {
  err := os.RemoveAll("./mock")
  if err != nil {
    return err
  }
  return nil 
}

func ReadFileToStringArray(path string) ([]string, error) {
  var fileLines []string

  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    fileLines = append(fileLines, scanner.Text())
  }

  if err := scanner.Err(); err != nil {
    return nil, err
  }
  return fileLines, nil
}
