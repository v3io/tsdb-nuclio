# Reference TSDB functions for Nuclio

### Prerequisites
1. Docker
2. `nuctl` (https://github.com/nuclio/nuclio/releases)
3. `kubeconfig` under `~/.kube/config` pointing to proper Kubernetes cluster with proper privileges (e.g. `cluster-admin`)
4. A data container HTTP URL, username/password for it and a TSDB table under it
5. The URL and credentials of the Docker registry Nuclio was configured to work with

To build without internet connectivity, you will need the Golang onbuild image (`nuclio/handler-builder-golang-onbuild:0.5.8-amd64-alpine`) and base image (`alpine:3.7`) present locally.

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

> Note: To prevent the build process from accessing the internet, set the `NUCLIO_BUILD_OFFLINE` environment variable to `true`

Push the images to a Docker registry:
```sh
TSDB_DOCKER_REPO=<TSDB_DOCKER_REPO> make push
```

> Note:
> 1. You must be logged into this Docker registry (`docker login`)
> 2. `TSDB_DOCKER_REPO` can be something like `tsdbtest.azurecr.io` (private registry in ACR) or `iguaziodocker` (public, in Docker hub)

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

`nuctl` will report to which NodePort the function was bound to (31848 in this case):
```sh
nuctl (I) Function deploy complete {"httpPort": 31848}
```

Post a metric to the function with your favorite HTTP client:
```sh
echo '{
  "metric": "cpu",
  "labels": {
    "dc": "7",
    "hostname": "mybesthost"
  },
  "samples": [
    {
      "time": "now",
      "value": {
        "N": 95.2
      }
    },
    {
      "time": "now",
      "value": {
        "n": 86.8
      }
    }
  ]
}' | http http://<TSDB_APPNODE_IP>:<TSDB_INGEST_NODE_PORT>

```

Where:
- `TSDB_APPNODE_IP`: An IP address of one of the application nodes
- `TSDB_INGEST_NODE_PORT`: As printed by the previous step

