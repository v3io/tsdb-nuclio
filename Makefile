TSDB_TAG ?= unstable
TSDB_DOCKER_REPO ?= iguazio/
NUCLIO_BUILD_OFFLINE ?= false

.PHONY: tsdb
build: ingest query
	@echo Done

.PHONY: ingest
ingest:
	cd functions/ingest && docker build --build-arg NUCLIO_BUILD_OFFLINE=$(NUCLIO_BUILD_OFFLINE) -t ${TSDB_DOCKER_REPO}tsdb-ingest:$(TSDB_TAG) .

.PHONY: query
query:
	cd functions/query && docker build --build-arg NUCLIO_BUILD_OFFLINE=$(NUCLIO_BUILD_OFFLINE) -t ${TSDB_DOCKER_REPO}tsdb-query:$(TSDB_TAG) .

.PHONY: promrw
promrw:
	cd functions/promrw && docker build --build-arg NUCLIO_BUILD_OFFLINE=$(NUCLIO_BUILD_OFFLINE) -t ${TSDB_DOCKER_REPO}tsdb-promrw:$(TSDB_TAG) .

.PHONY: push
push:
	docker push $(TSDB_DOCKER_REPO)/tsdb-ingest:$(TSDB_TAG)
	docker push $(TSDB_DOCKER_REPO)/tsdb-query:$(TSDB_TAG)
