package settings

import (
	"testing"

	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
)

/*
      Tests for PreflightCheck
*/
func TestPreflightPass(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  CreateDefaultSettingsFile(util.MockAppDir)

  result := PreflightCheck(util.MockAppDir)

  assert.True(t, result)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestPreflightFail(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  result := PreflightCheck(util.MockAppDir)

  assert.False(t, result)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}


