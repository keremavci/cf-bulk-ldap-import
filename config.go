package main

import (
	"os"
)

type Config struct {
	LdapHost          string
	LdapPassword     string
	LdapBindDN       string
	LdapBaseDN string
	CFApiURL       string
	CFUsername     string
	CFPassword     string
	CFOrg          string
	CFDefaultSpace string
	ConfigFile		string
}

func DefaultConfig() *Config {
	return &Config{
		LdapHost:          os.Getenv("LdapHost"),
		LdapPassword:     os.Getenv("LdapPassword"),
		LdapBindDN: os.Getenv("LdapBindDN"),
		LdapBaseDN: os.Getenv("LdapBaseDN"),
		CFApiURL:       os.Getenv("CFApiURL"),
		CFUsername:     os.Getenv("CFUsername"),
		CFPassword:     os.Getenv("CFPassword"),
		CFOrg:          getEnv("CFOrg","system"),
		CFDefaultSpace: getEnv("CFDefaultSpace", "system"),
		ConfigFile:		os.Getenv("ConfigFile"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
