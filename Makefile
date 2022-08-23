# Copyright 2018 Iguazio
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
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

.PHONY: push
push:
	docker push $(TSDB_DOCKER_REPO)/tsdb-ingest:$(TSDB_TAG)
	docker push $(TSDB_DOCKER_REPO)/tsdb-query:$(TSDB_TAG)
