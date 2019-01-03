TSDB_TAG ?= unstable
NUCLIO_BUILD_OFFLINE ?= false
TSDB_DOCKER_REPO ?= iguazio

.PHONY: tsdb
build: ingest query
	@echo Done

.PHONY: ingest
ingest:
	cd functions/ingest && docker build --build-arg NUCLIO_BUILD_OFFLINE=$(NUCLIO_BUILD_OFFLINE) -t ${TSDB_DOCKER_REPO}/tsdb-ingest:$(TSDB_TAG) .

.PHONY: query
query:
	cd functions/query && docker build --build-arg NUCLIO_BUILD_OFFLINE=$(NUCLIO_BUILD_OFFLINE) -t ${TSDB_DOCKER_REPO}/tsdb-query:$(TSDB_TAG) .

.PHONY: push
push:
	docker push $(TSDB_DOCKER_REPO)/tsdb-ingest:$(TSDB_TAG)
	docker push $(TSDB_DOCKER_REPO)/tsdb-query:$(TSDB_TAG)
