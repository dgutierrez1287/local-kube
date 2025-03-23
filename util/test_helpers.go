package util

import (
  "os"
)

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
