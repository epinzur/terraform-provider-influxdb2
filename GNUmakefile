default: testacc 
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS)
#HOSTNAME=hashicorp.com
NAME=influxdb2
NAMESPACE=rltvty
VERSION=0.0.1

BINARY=terraform-provider-${NAME}

INSTALL_PATH=~/.local/share/terraform/plugins/localhost/providers/${NAMESPACE}/${NAME}/${VERSION}/linux_$(GOARCH)

ifeq ($(GOOS), darwin)
	INSTALL_PATH=~/Library/Application\ Support/io.terraform/plugins/localhost/providers/${NAMESPACE}/${NAME}/${VERSION}/darwin_$(GOARCH)
endif
ifeq ($(GOOS), "windows")
	INSTALL_PATH=%APPDATA%/HashiCorp/Terraform/plugins/localhost/providers/${NAMESPACE}/${NAME}/${VERSION}/windows_$(GOARCH)
endif

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

dev:
	mkdir -p $(INSTALL_PATH)	
	go build -o $(INSTALL_PATH)/$(BINARY) main.go

generate:
	terraform fmt -recursive ./examples/
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

.PHONY: testacc docs
