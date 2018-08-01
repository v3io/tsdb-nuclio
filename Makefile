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
	docker tag tsdb-ingest:latest iguaziodocker/tsdb-ingest:1.9.0
	docker push iguaziodcoker/tsdb-ingest:1.9.0
	docker tag tsdb-query:latest iguaziodocker/tsdb-query:1.9.0
	docker push iguaziodocker/tsdb-query:1.9.0
