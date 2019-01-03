# Data Validation Service

## Setup

`go get -u github.com/gorilla/mux`  
`go build`  
Run executable (ie `./validation-svc-go`)  
  
Locally on `localhost:8000`

## API

**POST** - `/api/validate`

Body
```
{
  "dataSetOne": [{..},{..},...],
  "dataSetTwo": [{..},{..},...]
}
```