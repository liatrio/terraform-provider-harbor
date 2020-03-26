LOCALDIR?=$$(pwd)

build:
	go build -o terraform-provider-harbor

install: build
	cp terraform-provider-harbor ~/.terraform.d/plugins/terraform-provider-harbor_v0.0.0

lint:
	docker run --rm -v $(LOCALDIR):/app -w /app golangci/golangci-lint:v1.24.0 golangci-lint run -v

local:
	./scripts/helmRelease.sh

testacc: install local
	./scripts/runAccTests.sh
