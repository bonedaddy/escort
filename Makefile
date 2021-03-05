.PHONY: cli-garbled
cli-garbled:
	garble -tiny -literals -seed=random -debugdir=out build -o escort ./cmd
	upx --ultra-brute --best escort

.PHONY: cli
cli:
	go build -o escort ./cmd