package main

import (
	"flag"

	"github.com/jamesqin-cn/ldap_auth/controllers"
	"github.com/jamesqin-cn/ldap_auth/models"
	"github.com/jamesqin-cn/ldap_auth/utils"
)

var (
	configFile = flag.String("config", "./conf/config.yml", "config file path")
)

func init() {
	flag.Parse()
	utils.Inject(models.GetDefaultDataProxy())

	cfg := utils.GetConfig(*configFile)
	if err := models.GetDefaultDataProxy().LdapModel.SetLdapConf(&cfg.LdapConf); err != nil {
		panic(err)
	}
}

func main() {
	cfg := utils.GetConfig(*configFile)
	server := utils.NewHttpServer(cfg.ServerConf.Listen)

	apiCtl := &controllers.ApiController{}
	server.POST("/api/auth", utils.HttpHandlerWrapper(apiCtl.Auth))
	server.POST("/api/list", utils.HttpHandlerWrapper(apiCtl.List))

	server.RunHttpServer()
}
