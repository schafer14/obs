# Observations

[![version](https://img.shields.io/badge/version-v1.0.0-success)](https://img.shields.io/badge/version-v1.0.0-success)

This is an observation service loosly based on the Observation & Measurement specification and the CSIRO implementation. It makes significant deviations from both the original specification and the CSIRO implementation. 

## Usage

You can either run a local instance of the server or use a publically available demo.

### Running Locally

This project makes heavy use of Google Cloud Platform (because it's inexpensive). To do so you will need a Google Cloud Platform project and a Firestore instance. You will also need go >=1.13 installed.

```console
go run ./cmd/api -firestore-project=$GCP_PROJECT_ID
```

### Public Demo

The public demo comes with no guarentees regarding data longevity. It is meant to explore not API not to record observations.

The public demo is available at: https://linked-data-land.appspot.com/v1/observatoins

## API Docs

API docs are available in swagger form: [api docs](./swagger.yaml)

## Run Unit Tests

```console
go test ./... -v --short
```

## Run Unit and Integration Tests

Integration tests require Google Firestore account in a Google Project. Make sure you have credentials [setup](https://developers.google.com/accounts/docs/application-default-credentials).

```console
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/creds.json
go test ./... -v -project $MY_GCP_PROJECT
```

## TODO

- [x] Better validation errors observations
- [x] Better validation and errors for filters
- [x] Observation 404 errors
- [] Authboss authentication
- [] Authorization
- [x] API docs and clients
- [] Digital object observations
- [] Firestore emulator for improved integration testing
- [] Publish to a message queue when observations are created
