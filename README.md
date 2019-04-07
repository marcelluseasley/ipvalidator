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

## Running API locally
Since during development, I used `gcloud auth application-default login` for the BigTable access, it was necessary to create a service account and generate a JSON key. This JSON key file `ipgeo-readrights.json` is called from the code, where the tables are accessed. This allows for just enough rights to read the tables, in order for the table lookups to occur.

### Docker
To run the API in a Docker container, you can use the included `Dockerfile`.
```
docker build -t ipvalidator .
```
Once the image has been built, you can run the container with:
```
docker run -d -p 8080:8080 ipvalidator
```