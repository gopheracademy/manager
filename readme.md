# ShowRunner

Conference management platform, bespoke for GopherCon, open source for all.

## How to Use

### Requirements

* Go 1.15
* Docker

### Build

To generate service files, build the front end and compile the app, run:

```shell
$ make manager
```

### Run
Start Jaeger and Postgres: 
```shell
$ cd docker && docker-compose up -d
```

Jaeger must be running for the server to start.

Start the manager server:
```shell
$ make run
```

### View

To view the web app visit [here](http://127.0.0.1:8000/).

To visit Jaeger, go [here](http://127.0.0.1:16686/).


### Clean up

To stop Jaeger and Postgres, run:
```shell
$ cd docker && docker-compose down
```

## Coming Soon
There isn't much to share yet.  Come back soon to see how you can participate.
