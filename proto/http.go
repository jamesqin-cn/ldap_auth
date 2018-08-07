package proto

import "github.com/jamesqin-cn/ldap_auth/utils"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthReply struct {
	ErrCode int               `json:"err_code"`
	ErrMsg  string            `json:"err_msg"`
	Data    *utils.LdapResult `json:"data"`
}

type ListRequest struct {
}

type ListReply struct {
	ErrCode int                `json:"err_code"`
	ErrMsg  string             `json:"err_msg"`
	Count   int                `json:"count"`
	Data    []utils.LdapResult `json:"data"`
}
