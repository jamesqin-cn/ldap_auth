package models

import "github.com/jamesqin-cn/ldap_auth/utils"

type DataProxy struct {
	LdapModel *utils.Ldap `inject:""`
}

var (
	DefaultDataProxy *DataProxy
)

func GetDefaultDataProxy() *DataProxy {
	if DefaultDataProxy == nil {
		DefaultDataProxy = new(DataProxy)
	}

	return DefaultDataProxy
}
