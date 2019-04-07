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

---

## Getting data into BigTable

The necessary CSV files are hosted as ZIP files at https://dev.maxmind.com/geoip/geoip2/geolite2/ as https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country-CSV.zip.

Extracting the zip file returns a few zip files. The two I needed are : 
* GeoLite2-Country-Blocks-IPv4.csv
```
network,geoname_id,registered_country_geoname_id,represented_country_geoname_id,is_anonymous_proxy,is_satellite_provider
1.0.0.0/24,2077456,2077456,,0,0
1.0.1.0/24,1814991,1814991,,0,0
1.0.2.0/23,1814991,1814991,,0,0
```
* GeoLite2-Country-Locations-en.csv
```
geoname_id,locale_code,continent_code,continent_name,country_iso_code,country_name,is_in_european_union
49518,en,AF,Africa,RW,Rwanda,0
51537,en,AF,Africa,SO,Somalia,0
69543,en,AS,Asia,YE,Yemen,0
99237,en,AS,Asia,IQ,Iraq,0
```
Using the BigTable Go API, I added the records to their relative tables, using Apply and ApplyBulk.