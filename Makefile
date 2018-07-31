.PHONY: ingest
ingest:
	cd functions/ingest && docker build -t tsdb-ingest:latest .

.PHONY: query
query:
	cd functions/query && docker build -t tsdb-query:latest .

.PHONY: tsdb
tsdb: ingest query
	@echo Done

.PHONY: push
push:
	docker tag tsdb-ingest:latest levrado1/tsdb-ingest:latest
	docker push levrado1/tsdb-ingest:latest
	docker tag tsdb-query:latest levrado1/tsdb-query:latest
	docker push levrado1/tsdb-query:latest
