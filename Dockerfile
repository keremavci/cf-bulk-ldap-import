FROM golang:1.12-alpine3.9 as builder
RUN mkdir -p $GOPATH/src/github.com/keremavci/cf-bulk-ldap-insert
ADD . $GOPATH/src/github.com/keremavci/cf-bulk-ldap-insert/
WORKDIR $GOPATH/src/github.com/keremavci/cf-bulk-ldap-insert
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cf-bulk-ldap-insert .


FROM govau/cf-cli
RUN mkdir -p /opt/cf-bulk-ldap-insert
COPY --from=builder /go/src/github.com/keremavci/cf-bulk-ldap-insert/cf-bulk-ldap-insert /opt/cf-bulk-ldap-insert
RUN chmod +x /opt/cf-bulk-ldap-insert/cf-bulk-ldap-insert
ENTRYPOINT ["/opt/cf-bulk-ldap-insert/cf-bulk-ldap-insert"]
