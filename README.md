# Reference TSDB functions for Nuclio

### Prerequisites
1. Docker
2. `nuctl` (https://github.com/nuclio/nuclio/releases)
3. `kubeconfig` under `~/.kube/config` pointing to proper Kubernetes cluster with proper privileges (e.g., `cluster-admin`)
4. A data container HTTP URL, username/password for it and a TSDB table under it
5. The URL and credentials of the Docker registry Nuclio was configured to work with

To build without internet connectivity, you will need the Golang onbuild image (`nuclio/handler-builder-golang-onbuild:0.5.11-amd64-alpine`) and base image (`alpine:3.7`) present locally.

### Building / deploying the functions

Clone this repository and `cd` into it:
```sh
mkdir tsdb-nuclio && \
    git clone https://github.com/v3io/tsdb-nuclio.git tsdb-nuclio/src/github.com/v3io/tsdb-nuclio && \
    cd tsdb-nuclio/src/github.com/v3io/tsdb-nuclio
```

> Note: The `make` commands here use Docker. If you need root privileges to run Docker commands you need to make with `sudo` or as root

Build the Nuclio functions:
```sh
make build
```

This will build local Docker images - `tsdb-ingest:latest` and `tsdb-query:latest`

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
nuctl create project tsdb \
    --namespace NUCLIO_NAMESPACE \
    --display-name 'Time-series reference'
```

Where:
- `NUCLIO_NAMESPACE`: The namespace to which the function will be deployed

Deploy the ingest and query functions to your cluster:
```sh
nuctl deploy \
    --namespace <NUCLIO_NAMESPACE> \
    --run-image <TSDB_DOCKER_REPO>/tsdb-ingest:latest \
    --runtime golang \
    --handler main:Ingest \
    --project-name tsdb \
    --readiness-timeout 10 \
    --data-bindings '{"db0": {"class": "v3io", "url": "<TSDB_CONTAINER_URL>", "secret": "<TSDB_CONTAINER_USERNAME>:<TSDB_CONTAINER_PASSWORD>"}}' \
    --env INGEST_V3IO_TSDB_PATH=<TSDB_TABLE_PATH> \
    --env INPUT_FORMAT=<INPUT_FORMAT> \
    tsdb-ingest

nuctl deploy \
    --namespace <NUCLIO_NAMESPACE> \
    --run-image <TSDB_DOCKER_REPO>/tsdb-query:latest \
    --runtime golang \
    --handler main:Query \
    --project-name tsdb \
    --readiness-timeout 10 \
    --data-bindings '{"db0": {"class": "v3io", "url": "<TSDB_CONTAINER_URL>", "secret": "<TSDB_CONTAINER_USERNAME>:<TSDB_CONTAINER_PASSWORD>"}}' \
    --env QUERY_V3IO_TSDB_PATH=<TSDB_TABLE_PATH> \
    tsdb-query
```

Where:
- `TSDB_DOCKER_REPO`: The Docker registry to which the function image shall be pushed to
- `TSDB_CONTAINER_URL`: The Iguazio container URL (e.g., `http://10.0.0.1:8081/bigdata`)
- `TSDB_CONTAINER_USERNAME`: The Iguazio container username
- `TSDB_CONTAINER_PASSWORD`: The Iguazio container password
- `TSDB_TABLE_PATH`: The TSDB table name (e.g., `mytsdb`)
- `NUCLIO_NAMESPACE`: The namespace to which the function will be deployed
- `INPUT_FORMAT`: The input format that this ingest function should expect. Valid options are `DEFAULT` or `TCOLLECTOR`. if this variable will not be set the function will assume `DEFAULT`

`nuctl` will report to which NodePort the function was bound to (31848 in this case):
```sh
nuctl (I) Function deploy complete {"httpPort": 31848}
```

Post a metric to the function with your favorite HTTP client:
```sh
echo '{
  "metric": "cpu",
  "labels": {
    "site_id": "0001",
    "device_id": "12"
  },
  "samples": [
    {
      "t": "1537724629000",
      "v": {
        "n": 95.2
      }
    }
  ]
}' | http http://<TSDB_APPNODE_IP>:<TSDB_INGEST_NODE_PORT>
```

Where:
- `TSDB_APPNODE_IP`: An IP address of one of the application nodes
- `TSDB_INGEST_NODE_PORT`: As printed by the previous step

You should receive a 200 OK with an empty body in response. Now execute a query through the query function:
```sh
echo '{
    "metric": "cpu",
    "step": "1m",
    "start_time": "1537724600000",
    "end_time": "1537724730000"
}' | http http://<TSDB_APPNODE_IP>:<TSDB_QUERY_NODE_PORT>
```

Where:
- `TSDB_APPNODE_IP`: An IP address of one of the application nodes
- `TSDB_QUERY_NODE_PORT`: As printed by the previous step

The response should be a 200 OK with the following body:
```json
[
    {
        "datapoints": [
            [
                95.2,
                1537724629000
            ]
        ],
        "target": "cpu{device_id=12,site_id=0001}"
    }
]
```
