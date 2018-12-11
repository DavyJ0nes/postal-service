# Postal Service

## Description

A CLI tool that wraps running numerous [Postman](https://www.getpostman.com)
collections against an API that is started using [docker-compose](https://docs.docker.com/compose)


## Install

To install simply run `go install` within the root of the repository.
You will need to have Go installed, see [here](https://golang.org/doc/install) for instructions

## Usage

Please refer the help of the command for how to use it:

```shell
$ postal-service --help
```

## Example

Within [examples directory](./examples) there are a few examples of how the tool can be used.
For example within the [basic-http](./examples/basic-http) directory you can run the tool and see
the following output:

```shell
$ postal-service run          
Running Postal Service
-- Initialising API

-- Running Postman Collection
-----------------------------
--- Running: basic-http.postman_collection.json
-----------------------------
-- Killing API
-----------------------------
basic-http

Iteration 1/2

→ Get Index
  GET http://web/ [200 OK, 850B, 120ms]

Iteration 2/2

→ Get Index
  GET http://web/ [200 OK, 850B, 5ms]

┌─────────────────────────┬──────────┬──────────┐
│                         │ executed │   failed │
├─────────────────────────┼──────────┼──────────┤
│              iterations │        2 │        0 │
├─────────────────────────┼──────────┼──────────┤
│                requests │        2 │        0 │
├─────────────────────────┼──────────┼──────────┤
│            test-scripts │        2 │        0 │
├─────────────────────────┼──────────┼──────────┤
│      prerequest-scripts │        0 │        0 │
├─────────────────────────┼──────────┼──────────┤
│              assertions │        4 │        0 │
├─────────────────────────┴──────────┴──────────┤
│ total run duration: 255ms                     │
├───────────────────────────────────────────────┤
│ total data received: 1.2KB (approx)           │
├───────────────────────────────────────────────┤
│ average response time: 62ms                   │
└───────────────────────────────────────────────┘
```


