## tl;tr
Find ldap user via groups and add Cloud Foundry

## Build
```
docker build -t cf-bulk-ldap-import .
```
## Run
```
docker run -it \
-v ${PWD}/config/cfuserrole.yml:/tmp/cfuserrole.yml \
-e LdapHost="Ldap:389" \
-e LdapPassword="$LdapPass*" \
-e LdapBindDN="CN=LDAP USER,OU=Resource,DC=keremavci,DC=dev" \
-e LdapBaseDN="OU=developers,DC=keremavci,DC=dev" \
-e CFApiURL="https://api.system.prod.pcf" \
-e CFUsername="admin" \
-e CFPassword="password" \
-e ConfigFile=/tmp/cfuserrole.yml \
 cf-bulk-ldap-import
```
