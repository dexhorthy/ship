.PHONY: build deps prebuild build gen-schema merge-mutations validate done markdown-assets pipeline-strict pipeline copy-to-help-center copy-replicated-lint ship

SCHEMA_JSON := ./schema.json
MUTATIONS := ./mutations.json
HACK_DOCS := ./bin/hack-docs

.state/deps: package.json yarn.lock
	yarn
	@mkdir -p .state
	@touch .state/deps

deps: .state/deps

prebuild:
	rm -rf assets
	rm -rf lifecycle
	rm -rf config
	mkdir -p assets
	mkdir -p lifecycle
	mkdir -p config

build:
	mkdir -p bin
	go build  -o $(HACK_DOCS) ../../cmd/hack-docs

gen-schema:
	$(HACK_DOCS) > $(SCHEMA_JSON)

merge-mutations:
	`yarn bin`/ts-node src/index.ts merge-mutations -f $(SCHEMA_JSON) -m $(MUTATIONS)  > $(SCHEMA_JSON).tmp
	mv $(SCHEMA_JSON).tmp $(SCHEMA_JSON)

validate:
	`yarn bin`/ts-node --inspect=0.0.0.0 src/index.ts validate -f $(SCHEMA_JSON)

done:
   : ✓ All Clear!

markdown-assets:
	mkdir -p build
	`yarn bin`/ts-node --inspect=0.0.0.0 src/index.ts markdown-assets -f $(SCHEMA_JSON) -o ./assets


pipeline-strict: deps prebuild build gen-schema merge-mutations validate markdown-assets
pipeline: deps prebuild build gen-schema merge-mutations 
	
# require validation to generate docs, pipeline only does schema
pipeline-docs: pipeline-strict

copy-to-help-center:
	: To copy to the Ship documentation site, run the following command:
	:
	:    make pipeline
	:    cp assets/* [path/to/ship-www]/content/reference/assets/
	:    cp lifecycle/* [path/to/ship-www]/content/reference/lifecycle/
	:    cp config/* [path/to/ship-www]/content/reference/config/

copy-to-help-center-gopath:
	cp assets/* $$GOPATH/src/github.com/replicatedhq/ship-www/content/reference/assets/
	cp lifecycle/* $$GOPATH/src/github.com/replicatedhq/ship-www/content/reference/lifecycle/
	cp config/* $$GOPATH/src/github.com/replicatedhq/ship-www/content/reference/config/

copy-replicated-lint:
	: To Copy to replicated lint, run the following command:
	:
	:    make pipeline
	:    cp schema.json  [path/to/replicated-lint]/projects/replicated-ship/
