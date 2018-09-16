# default tag to "latest"
TSDB_TAG := $(if $(TSDB_TAG),$(TSDB_TAG),latest)

# default NUCLIO_BUILD_OFFLINE to "false"
NUCLIO_BUILD_OFFLINE := $(if $(NUCLIO_BUILD_OFFLINE),$(NUCLIO_BUILD_OFFLINE),latest)

.PHONY: ingest
ingest:
	cd functions/ingest && docker build --build-arg NUCLIO_BUILD_OFFLINE=$(NUCLIO_BUILD_OFFLINE) -t tsdb-ingest:$(TSDB_TAG) .

.PHONY: tsdb
tsdb: ingest
	@echo Done

.PHONY: push
push:
	docker tag tsdb-ingest:$(TSDB_TAG) $(TSDB_DOCKER_REPO)/tsdb-ingest:$(TSDB_TAG)
	docker push $(TSDB_DOCKER_REPO)/tsdb-ingest:$(TSDB_TAG)
