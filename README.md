# Reference TSDB functions for Nuclio

### Prerequisites
1. Docker
2. `nuctl` (https://github.com/nuclio/nuclio/releases)
3. `kubeconfig` under `~/.kube/config` pointing to proper Kubernetes cluster with proper privileges (e.g. `cluster-admin`)
4. A data container HTTP URL, username/password for it and a TSDB table created under `$TSDB_TSDB_TABLE_NAME`
5. The URL and credentials of the Docker registry Nuclio was configured to work with

### Building / deploying the functions

Clone this repository and `cd` into it:
```sh
mkdir tsdb-nuclio && \
    git clone git@github.com:v3io/tsdb-nuclio.git tsdb-nuclio/src/github.com/v3io/tsdb-nuclio && \
    cd tsdb-nuclio/src/github.com/v3io/tsdb-nuclio
```

Build the Nuclio functions:
```sh
make tsdb
```

This will build a local Docker image - `tsdb-ingest:latest`

> Note: To prevent the build process from accessing the internet, set `NUCLIO_BUILD_OFFLINE` to `true`. In such a case, you must make sure that the image `nuclio/handler-builder-golang-onbuild:0.5.8-amd64-alpine` is present locally.

Push the images to a Docker registry:
```sh
TSDB_DOCKER_REPO=<TSDB_DOCKER_REPO> make push
```

> Note: You must be logged into this Docker registry (`docker login`)

Create a project for the functions:
```sh
nuctl create project tsdb --display-name 'Time-series reference'
```

Deploy the ingest function to your cluster:
```sh
nuctl deploy \
  --namespace <NUCLIO_NAMESPACE> \
  --run-image <TSDB_DOCKER_REPO>/tsdb-ingest:latest \
  --runtime golang \
  --handler main:Ingest \
  --project-name tsdb \
  --readiness-timeout 10 \
  --data-bindings '{"db0": {"class": "v3io", "url": "<TSDB_CONTAINER_URL>", "secret": "<TSDB_CONTAINER_USERNAME>:<TSDB_CONTAINER_PASSWORD>"}}' \
  --env INGEST_V3IO_TSDB_PATH=<TSDB_TSDB_TABLE_NAME> \
  tsdb-ingest
```

Where:
- `TSDB_DOCKER_REPO`: The Docker registry to which the function image shall be pushed to
- `TSDB_CONTAINER_URL`: The Iguazio container URL (e.g. `http://10.0.0.1:8081/bigdata`)
- `TSDB_CONTAINER_USERNAME`: The Iguazio container username
- `TSDB_CONTAINER_PASSWORD`: The Iguazio container password
- `TSDB_TSDB_TABLE_NAME`: The TSDB table name (e.g. `mytsdb`)
- `NUCLIO_NAMESPACE`: The namespace to which the the function will be deployed
