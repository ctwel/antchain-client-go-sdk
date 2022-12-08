package utils

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSign(t *testing.T) {
	_, err := Sign("hello", os.Getenv("GOPATH") + "/src/github.com/ctwel/antchain-client-go-sdk/test/access.key")
	require.Truef(t,err == nil,"sign text failed,err:%+v",err)
}
