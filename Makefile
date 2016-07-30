install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/bot_agent_auth/bot_agent_auth.go
test:
	GO15VENDOREXPERIMENT=1 go test `glide novendor`
runledis:
	ledis-server \
	-addr=localhost:5555 \
	-databases=1
runauth:
	auth_server \
	-loglevel=debug \
	-port=6666 \
	-ledisdb-address=localhost:5555 \
	-auth-application-password=test123
run:
	bot_agent_auth \
	-loglevel=debug \
	-port=7777 \
	-ledisdb-address=localhost:5555 \
	-prefix "/wiki" \
	-root "./files"
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/bborbe/auth/bin/auth_server
	go get -u github.com/siddontang/ledisdb/cmd/ledis-server
	glide install
update:
	glide up
clean:
	rm -rf var vendor