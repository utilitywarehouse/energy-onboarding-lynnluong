# energy-onboarding-lynnluong

[![CI](https://github.com/utilitywarehouse/energy-onboarding-lynnluong/actions/workflows/ci.yml/badge.svg)](https://github.com/utilitywarehouse/energy-onboarding-lynnluong/actions/workflows/ci.yml)

Give a high level explanation of what your service does.

Explain the business requirement it solves, explaining more about the process that end users might follow with your service.

And a brief overview of how an end user would interact with your service.

## Setup

To install, run go get, then follow build instructions from the new directory.

```bash
go get github.com/utilitywarehouse/energy-onboarding-lynnluong
cd  $GOPATH/src/github.com/utilitywarehouse/energy-onboarding-lynnluong
```

### Dependencies

* A DB
* A Kafka topic
* Another Service

### Usage

This should cover more technical details of how your service is actually run, this includes CLI options, Endpoints for an API, and Prometheus metrics.

#### Options

```bash
  --port=8080         The port to listen on for HTTP connections ($HTTP_PORT)
  --opsPort=8081      The port to listen on for HTTP connections for Ops ($OPS_PORT)
```

#### Endpoints

##### Ports

What if any ports are exposed, and used to interact with your app. These should be the external ports that kubernets exposes.

`80/443`

#### Routes

Is this an API that has endpoints?  If so list endpoints, broken by their httpMethod. For Post methods describe what the route returns, for routes with variable endings, give an example. For Post/Put methods give a curl example.

##### Get Methods

* `/api/v1/get`

Returns wonderful data

* `/api/v1/get/{something}`

Returns something
Eg. `/api/v1/get/thisIsSomething`

##### Post Methods

* `/api/v1/post`

Post some data to this route

```bash
curl -X "POST" "http://localhost:8080/api/v1/post" \
     -H "Content-Type: application/json" \
     -d $'{
  "data": "some data"
}'
```

#### Metrics

What are the Prometheus metrics that have been implemented? Explaining which metrics are performance vs error metrics, to help others diagnosis your service.

##### Performance

`TotalJobs`
Total number of jobs run

##### Errors

`ErrCount`
The total number of errors

## Building

```bash
go get -u -t ./...
go build cmd/energy-onboarding-lynnluong/main.go
```

## Testing

Are there any steps that need to be taken before a user is able to run your tests?

For example dependant containers to spin up via the docker-compose file, or environmental variables that need to be exported.

Set DB_CONNECTION_STRING to localhost

```bash
export DB_CONNECTION_STRING=localhost
go test ./...
```
