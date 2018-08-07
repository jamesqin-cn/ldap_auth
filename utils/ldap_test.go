package utils

import (
	"fmt"
	"log"
	"testing"
)

var ldapConf = &LdapConf{
	Addr:       "10.255.8.254:389",
	BaseDn:     "OU=company,DC=example,DC=net",
	Filter:     "(sAMAccountName=%s)",
	Attributes: []string{"sAMAccountName", "mail", "telephoneNumber", "name", "memberOf", "userAccountControl"},
	BindDn:     "CN=administrator,CN=Users,DC=example,DC=net",
	BindPasswd: "company&2018",
}

func TestAuth(t *testing.T) {
	ldap := Ldap{}
	if err := ldap.SetLdapConf(ldapConf); err != nil {
		log.Fatal(err)
	}

	info, err := ldap.Auth("not_exist_user", "no_passwd")
	if err == nil {
		t.Error("test Auth() failed, expect user not exist")
	}

	info, err = ldap.Auth("meeting", "abcd")
	if err == nil {
		t.Error("test Auth() failed, expect password not match")
	}

	info, err = ldap.Auth("meeting", "abcd.1234")
	if err != nil {
		t.Error("test Auth() failed, expect userinfo from ldap server")
	}
	fmt.Println(info)
}

func TestListUsers(t *testing.T) {
	ldap := Ldap{}
	if err := ldap.SetLdapConf(ldapConf); err != nil {
		log.Fatal(err)
	}

	list, err := ldap.ListUsers()
	if err != nil {
		t.Errorf("test ListUsers() failed, ", err)
	}
	fmt.Println(list)
}
