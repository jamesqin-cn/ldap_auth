package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	ldap "gopkg.in/ldap.v2"
)

type LdapConf struct {
	Addr       string   `yaml:"addr"`
	BaseDn     string   `yaml:"base_dn"`
	Filter     string   `yaml:"filter"`
	Attributes []string `yaml:"attributes,flow"`
	BindDn     string   `yaml:"bind_dn"`
	BindPasswd string   `yaml:"bind_passwd"`
}

type LdapResult struct {
	DN         string
	Attributes map[string][]string
}

type Ldap struct {
	conf *LdapConf `inject:""`
}

func (l *Ldap) SetLdapConf(conf *LdapConf) error {
	if err := VerifyLdapConfig(conf); err != nil {
		return err
	}
	l.conf = conf
	return nil
}

func VerifyLdapConfig(conf *LdapConf) (err error) {
	if len(conf.Addr) == 0 || !strings.Contains(conf.Addr, ":") {
		return errors.New("field [ ldap.addr ] missing or wrong format, expect x.x.x.x:389")
	}

	if len(conf.BaseDn) == 0 || !strings.Contains(conf.BaseDn, "DC=") {
		return errors.New("field [ ldap.base_dn ] missing or wrong format, expect OU=xxx,DC=example,DC=com")
	}

	if len(conf.Filter) == 0 {
		return errors.New("field [ ldap.filter ] missing or wrong format")
	}

	if len(conf.Attributes) == 0 {
		return errors.New("field [ ldap.attributes ] missing or wrong format, attributes must be an array")
	}

	if len(conf.BindDn) == 0 || !strings.Contains(conf.BindDn, "DC=") {
		return errors.New("field [ ldap.bind_dn ] missing or wrong format, expect OU=xxx,DC=example,DC=com")
	}

	if len(conf.BindPasswd) == 0 {
		return errors.New("field [ ldap.bind_passwd ] missing or wrong format")
	}

	return nil
}

func (m *Ldap) Auth(username string, password string) (info *LdapResult, err error) {
	userInfo := &LdapResult{}

	ldap.DefaultTimeout = 5 * time.Second

	l, err := ldap.Dial("tcp", m.conf.Addr)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(m.conf.BindDn, m.conf.BindPasswd)
	if err != nil {
		return nil, err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		m.conf.BaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(m.conf.Filter, username),
		m.conf.Attributes,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) != 1 {
		return nil, errors.New("User does not exist or too many entries returned")
	}

	attributes := make(map[string][]string)
	for _, item := range sr.Entries[0].Attributes {
		attributes[item.Name] = item.Values
	}

	userInfo.DN = sr.Entries[0].DN
	userInfo.Attributes = attributes

	// Bind as the user to verify their password
	err = l.Bind(userInfo.DN, password)
	if err != nil {
		return nil, errors.New("wrong password of user " + username)
	}

	// Rebind as the read only user for any further queries
	err = l.Bind(m.conf.BindDn, m.conf.BindPasswd)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (m *Ldap) ListUsers() (list []LdapResult, err error) {
	userList := make([]LdapResult, 0)

	ldap.DefaultTimeout = 5 * time.Second

	l, err := ldap.Dial("tcp", m.conf.Addr)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(m.conf.BindDn, m.conf.BindPasswd)
	if err != nil {
		return nil, err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		m.conf.BaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(m.conf.Filter, "*"),
		m.conf.Attributes,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) == 0 {
		return nil, errors.New("User does not exist")
	}

	for _, userEntry := range sr.Entries {
		userInfo := LdapResult{}
		attributes := make(map[string][]string)
		for _, item := range userEntry.Attributes {
			attributes[item.Name] = item.Values
		}

		userInfo.DN = userEntry.DN
		userInfo.Attributes = attributes
		userList = append(userList, userInfo)
	}

	// Rebind as the read only user for any further queries
	err = l.Bind(m.conf.BindDn, m.conf.BindPasswd)
	if err != nil {
		return nil, err
	}

	return userList, nil
}
