# rest-geoip

This application replies with information about a client using the [Maxmind](https://www.maxmind.com) GeoLite2-City database.

The webapp provides general geoip information. There is also an api available

```
GET  /api/ip     : Return client IP Address
GET  /api/geoip  : Return client Geoip information
POST /api/update : Update the Maxmind database
```

The application doesn't provide a database. POSTing to `/api/update` will update the database. Ideally, this endpoint should be protected via Nginx or whatever is available.
