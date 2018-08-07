package controllers

import (
	"context"

	"github.com/jamesqin-cn/ldap_auth/models"
	"github.com/jamesqin-cn/ldap_auth/proto"
)

type ApiController struct {
}

func (ctl *ApiController) Auth(ctx context.Context, req *proto.AuthRequest) (reply *proto.AuthReply) {
	info, err := models.GetDefaultDataProxy().LdapModel.Auth(req.Username, req.Password)
	if err != nil {
		return &proto.AuthReply{
			ErrCode: 2,
			ErrMsg:  err.Error(),
		}
	}

	return &proto.AuthReply{
		ErrCode: 0,
		ErrMsg:  "ok",
		Data:    info,
	}
}

func (ctl *ApiController) List(ctx context.Context, req *proto.ListRequest) (reply *proto.ListReply) {
	list, err := models.GetDefaultDataProxy().LdapModel.ListUsers()
	if err != nil {
		return &proto.ListReply{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}
	}

	return &proto.ListReply{
		ErrCode: 0,
		ErrMsg:  "ok",
		Count:   len(list),
		Data:    list,
	}
}
