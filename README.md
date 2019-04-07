# IP Validator API

Prerequisite to API working: lookup data in BigTable; see [BigTable Setup](https://github.com/marcelluseasley/ipvalidator/tree/master/bigtable_setup)

The purpose of this API is to receive data sent from the API gateway, which contains the requesting IP address and a list of white-listed countries. It wasn't clear how the data would be supplied to the API endpoint.

This API endpoint assumes that the data is received via a POST request of JSON data, which will look like this:

* ip validation: `true`
```json
{
   "ip":"172.58.7.17",
   "approved_countries":[
      "Russia",
      "Germany",
      "Mexico",
      "United States"
   ]
}
```
* ip validation: `false`
```json
{
   "ip":"172.58.7.17",
   "approved_countries":[
      "Russia",
      "Germany",
      "Mexico"
   ]
}
```

The endpoint response will be JSON and either contain a `valid_status` of `true` or 'false`:
```json
{
    "valid_status": "true"
}
```
```json
{
    "valid_status": "false"
}
```
Note: see `bigtable_search.go` for details on the IP/Country lookup logic.