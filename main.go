package main

import (
	"net"
	"strconv"

	"github.com/go-ldap/ldap/v3"
)

const (
	addr     = "127.0.0.1"
	port     = 389
	protocol = "tcp"
	base     = "dc=example,dc=co,dc=th"
	username = "cn=admin,dc=example,dc=co,dc=th"
	password = "P@ssw0rd"
)

type client struct {
	*ldap.Conn
}

func connect() (*client, error) {
	c, err := ldap.Dial(protocol, net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}

	if err = c.Bind(username, password); err != nil {
		return nil, err
	}

	return &client{c}, nil
}

func main() {
	c, err := connect()
	if err != nil {
		panic(err)
	}

	defer c.Close()

	var search = &ldap.SearchRequest{
		BaseDN:       base,
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldap.NeverDerefAliases,
		SizeLimit:    0,
		TimeLimit:    0,
		TypesOnly:    false,
		Filter:       "(&(objectClass=*))",
		Attributes:   []string{"*"},
		Controls:     nil,
	}

	results, err := c.Search(search)
	if err != nil {
		panic(err)
	}

	for _, ent := range results.Entries {
		println(ent.DN)
	}
}
