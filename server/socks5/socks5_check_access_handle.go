package socks5

import (
	"context"
	"errors"
	"github.com/cnlh/nps/core"
	"net"
)

type CheckAccess struct {
	core.NpsPlugin
	clientConn     net.Conn
	clientUsername string
	clientPassword string
	configUsername string
	configPassword string
}

func (check *CheckAccess) GetConfigName() *core.NpsConfigs {
	c := core.NewNpsConfigs("socks5_simple_access_check", "need check the permission simply")
	c.Add("socks5_simple_access_username", "simple auth username")
	c.Add("socks5_simple_access_password", "simple auth password")
	return c
}

func (check *CheckAccess) Run(ctx context.Context, config map[string]string) error {
	check.clientConn = check.GetClientConn(ctx)
	check.configUsername = config["socks5_access_username"]
	check.configPassword = config["socks5_access_password"]

	return nil
}

func (check *CheckAccess) checkAuth(configUserName, configPassword string) error {
	if check.clientUsername == configUserName && check.clientPassword == configPassword {
		_, err := check.clientConn.Write([]byte{userAuthVersion, authSuccess})
		return err
	} else {
		_, err := check.clientConn.Write([]byte{userAuthVersion, authFailure})
		if err != nil {
			return err
		}
		return errors.New("auth check error,username or password does not match")
	}
}
