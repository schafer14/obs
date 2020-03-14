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

TODO:

- [] Better validation errors observations
- [] Better validation and errors for filters
- [] Authboss stuff
- [] File uploader
- [] Observation 404 errors
