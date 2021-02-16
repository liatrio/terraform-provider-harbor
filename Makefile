LOCALDIR?=$$(pwd)

.PHONY: build install lint local testacc

build:
	go build -o terraform-provider-harbor

install: build
	mkdir -p ~/.terraform.d/plugins/terraform.local/liatrio/harbor/0.0.1/darwin_amd64
	cp terraform-provider-harbor ~/.terraform.d/plugins/terraform.local/liatrio/harbor/0.0.1/darwin_amd64/

lint:
	docker run --rm -v $(LOCALDIR):/app -w /app golangci/golangci-lint:v1.24.0 golangci-lint run -v --timeout 1h

local:
	./scripts/helmRelease.sh

testacc: install local
	./scripts/runAccTests.sh
