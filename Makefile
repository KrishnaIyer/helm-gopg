.PHONY: init

init:
	@echo "Initialize repository..."

.PHONY: test

test:
	@echo "Run tests..."
	@helm create test/chart
	@helm package test/chart
	@go run main.go sign --package chart-0.1.0.tgz --signer.pgp.passphrase foobar --signer.pgp.private-key ./test/keys/priv.key
	@go run main.go verify --package chart-0.1.0.tgz --signer.pgp.public-key ./test/keys/pub.key
	@echo "Clean up..."
	@rm -rf test/chart
	@rm chart-0.1.0.tgz
	@rm chart-0.1.0.tgz.prov

