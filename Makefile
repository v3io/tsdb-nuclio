# default tag to "latest"
TSDB_TAG := $(if $(TSDB_TAG),$(TSDB_TAG),latest)

# default NUCLIO_BUILD_OFFLINE to "false"
NUCLIO_BUILD_OFFLINE := $(if $(NUCLIO_BUILD_OFFLINE),$(NUCLIO_BUILD_OFFLINE),latest)

.PHONY: ingest
ingest:
	cd functions/ingest && docker build --build-arg NUCLIO_BUILD_OFFLINE=$(NUCLIO_BUILD_OFFLINE) -t tsdb-ingest:$(TSDB_TAG) .

.PHONY: query
query:
	cd functions/query && docker build --build-arg NUCLIO_BUILD_OFFLINE=$(NUCLIO_BUILD_OFFLINE) -t tsdb-query:$(TSDB_TAG) .

.PHONY: tsdb
build: ingest query
	@echo Done

.PHONY: push
push:
	docker tag tsdb-ingest:$(TSDB_TAG) $(TSDB_DOCKER_REPO)/tsdb-ingest:$(TSDB_TAG)
	docker push $(TSDB_DOCKER_REPO)/tsdb-ingest:$(TSDB_TAG)
	docker tag tsdb-query:$(TSDB_TAG) $(TSDB_DOCKER_REPO)/tsdb-query:$(TSDB_TAG)
	docker push $(TSDB_DOCKER_REPO)/tsdb-query:$(TSDB_TAG)
