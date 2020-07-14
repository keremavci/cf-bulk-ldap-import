package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"gopkg.in/yaml.v2"
)

type Organizations struct {
	Organizations []Organization
}
type Organization struct {
	Name   string  `yaml:"name"`
	Spaces []Space `yaml:"spaces"`
}

type Space struct {
	Space            string   `yaml:"space"`
	Groups           []string `yaml:"groups"`
	Roles            []string `yaml:"roles"`
	ExcludeFromSpace []string `yam:"excludefromspace"`
}

func main() {

	var err error

	config := DefaultConfig()

	organizations := Organizations{}

	userData, err := ioutil.ReadFile(config.ConfigFile)
	if err != nil {
		log.Fatalf("Failed to load config file %v", err)
	}
	if err = yaml.Unmarshal(userData, &organizations); err != nil {
		log.Fatalf("Failed to parse config file %v", err)
	}

	ldapConn, err := ldap.Dial("tcp", config.LdapHost)
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer ldapConn.Close()

	if err = ldapConn.Bind(config.LdapBindDN, config.LdapPassword); err != nil {

		log.Fatalf("ldap: initial bind for user %q failed: %v", config.LdapBindDN, err)
	}

	loginCFApi(config)
	log.Println("Login to PCF")

	for _, org := range organizations.Organizations {
		for _, space := range org.Spaces {
			createPcfUserSetSpaceRole(org.Name, space, config, ldapConn)
			if len(space.ExcludeFromSpace) > 0 {
				unSetSpaceRole(org.Name, space, config, ldapConn)
			}
		}
	}

}

func getUsersByGroupName(groups []string, config *Config, ldapConn *ldap.Conn) []string {
	users := make([]string, 0)

	for _, g := range groups {
		filter := fmt.Sprintf("(&(objectCategory=user)(memberOf=cn=%s,OU=Groups,%s))", g, config.LdapBaseDN)
		result, _ := ldapConn.Search(ldap.NewSearchRequest(
			config.LdapBaseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			filter,
			[]string{"sAMAccountName"},
			nil,
		))
		for i := range result.Entries {
			users = append(users, result.Entries[i].Attributes[0].Values[0])
		}
	}

	return users

}

func createPcfUserSetSpaceRole(organization string, space Space, config *Config, ldapConn *ldap.Conn) {
	users := getUsersByGroupName(space.Groups, config, ldapConn)
	for _, user := range users {

		cmd := strings.Split(string(fmt.Sprintf("create-user %s --origin ldap", user)), " ")
		runCFCli(cmd)
		for _, role := range space.Roles {
			cmd := strings.Split(string(fmt.Sprintf("set-space-role %s %s %s %s", user, organization, space.Space, role)), " ")
			runCFCli(cmd)

		}

	}
}

func unSetSpaceRole(organization string, space Space, config *Config, ldapConn *ldap.Conn) {
	notPermittedUsers := getUsersByGroupName(space.ExcludeFromSpace, config, ldapConn)

	for _, user := range notPermittedUsers {
		for _, role := range space.Roles {
			cmd := strings.Split(string(fmt.Sprintf("unset-space-role %s %s %s %s", user, organization, space.Space, role)), " ")
			runCFCli(cmd)
		}
	}
}

func loginCFApi(config *Config) {
	runCFCli(strings.Split(string(fmt.Sprintf("login -a %s -u %s -p %s -o %s -s %s --skip-ssl-validation", config.CFApiURL, config.CFUsername, config.CFPassword, config.CFOrg, config.CFDefaultSpace)), " "))
}

func runCFCli(args []string) {

	c := strings.Join(args, " ")
	log.Println(fmt.Sprintf("Running Command: cf %s", c))
	out, err := exec.Command("cf", args...).Output()
	log.Println(string(out))
	if err != nil {
		log.Println(err)
	}

}
