package utils

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSign(t *testing.T) {
	_, err := Sign("hello", os.Getenv("GOPATH") + "/src/gitlab.alipay-inc.com/antchain/restclient-go-sdk/test/access.key")
	require.Truef(t,err == nil,"sign text failed,err:%+v",err)
}
