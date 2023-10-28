# rest-geoip

[![Go Report Card](https://goreportcard.com/badge/github.com/TwistTheNeil/rest-geoip)](https://goreportcard.com/report/github.com/TwistTheNeil/rest-geoip)
[![Golangci-lint Status](https://github.com/TwistTheNeil/rest-geoip/workflows/golangci-lint/badge.svg)](https://github.com/TwistTheNeil/rest-geoip/actions?query=workflow%3Agolangci-lint)
[![Docker build Status](https://github.com/TwistTheNeil/rest-geoip/workflows/Docker%20Image%20CI/badge.svg)](https://github.com/TwistTheNeil/rest-geoip/actions?query=workflow%3A%22Docker+Image+CI%22)
[![Build Report](https://github.com/TwistTheNeil/rest-geoip/workflows/go-build/badge.svg)](https://github.com/TwistTheNeil/rest-geoip/actions?query=workflow%3Ago-build)

** For v0.4.1, the previous stable, look in the appropriate branch **

A self hosted geoip lookup application written in Go and Vue.js 3 which provides a client with information about their IP address or any other. It uses the [Maxmind](https://www.maxmind.com) GeoLite2-City database.

The webapp provides general geoip information. There is also an api available

```
GET  /                    : Return client IP Address (when used with curl or HTTPie)
GET  /api/geoip           : Return client Geoip information
GET  /api/geoip/:address  : Return Geoip information for "address"
PUT  /api/update          : Update the Maxmind database
```

The application doesn't provide a database. A `PUT` request to `/api/update` will update the database and will ideally be protected by an api key (header: `X-API-KEY`). If `API_KEY` env var is not set, then the application will set one on startup and notify via STDOUT

### Screenshots of optional webapp
![screenshot](docs/screen.png)
