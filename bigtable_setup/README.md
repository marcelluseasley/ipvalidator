# BigTable Tables Setup

When deciding on how to best implement the IP Validator API, I needed a place to store the tables from [Geolite 2](https://dev.maxmind.com/geoip/geoip2/geolite2/) website. At first, I was going to store the necessary csv file data in a MySQL database, but I recall Randy telling me that BigTable is one of the technologies used in the stack. So I chose BigTable to store the data.

## ip_country Database

The name of the database where I store the tables is called `ip_country` and contains two tables: `Geoname-Country` and `IP-Country`.

```
$ cbt ls
Geoname-Country
IP-Country
```



### IP-Country table
Table `IP-Country` contains a family name of `ipGeoID`, which stores the geoname IDs:
```
$ cbt ls IP-Country
Family Name     GC Policy
-----------     ---------
ipGeoID         <never>
```
Looking up a network mask, returns the corresponding Geoname ID, which is stored in a column named `geoname_id`:
```
$ cbt lookup IP-Country 172.58.0.0/17
----------------------------------------
172.58.0.0/17
  ipGeoID:geoname_id                       @ 2019/04/06-13:05:28.665000
    "6252001"
```

### Geoname-Country table
Table `Geoname-Country` contains a family name of `geoIDCountry`, which stores the countries:
```
$ cbt ls Geoname-Country
Family Name     GC Policy
-----------     ---------
geoIDCountry    <never>

```
The IP validation APi takes the geoname ID from the `IP-Country` table and uses that value to lookup the relative country from the `Geoname-Country` table:
```
$ cbt lookup Geoname-Country 6252001
----------------------------------------
6252001
  geoIDCountry:country                     @ 2019/04/06-01:20:28.760000
    "United States"
```

If this country is in the list of approved countries, the API endpoint returns `true`. 