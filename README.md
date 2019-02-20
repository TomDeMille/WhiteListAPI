# Go REST API Example
A RESTfull API example utilizing Golang and MaxMind GEO DB

## Installation & Run
```go get github.com/TomDeMille/WhiteListAPI```

or download using SSH/HTTPS from Git  

```Git Clone github.com/TomDeMille/WhiteListAPI ```


```bash
# Build and Run
cd ~/go/src/whiteListApi
go get
go build
./whiteListApi

# Check This : http://127.0.0.1:8080/v1/api/ping
```

## Structure
```
├── whiteListApi
│   ├── main.go                         // main... nuff said
│   ├── routes                          // Our API core handlers
│   │   └── country    
│   │           └── country.go          //country based requests
│   └── DB
│       ├── db.go                       // DB instantiation and DAL
│       └── GeoLite2-Cty.mmdb           // MaxMind DB included

```

## API

NOTE : See included whiteListApi.postman_collection.json for example requests, 
including the POST format

 **( /v1/api/... )**

#### /ping
* `GET` : Simple heartbeat request returns 200 if alive

#### /country/whitelistrequest
* **`POST`** : Post a JSON WhiteListRequest, returns a WhiteListResponse as JSON

#### /country/latlngbyip/{ipAddress} `
* **`GET`** : Given the IP returns a JSON object with Latitude and Longitude
 
#### /country/timezonebyip/{ipAddress} `
* **`GET`** : Given the IP returns a JSON object with the timezone

#### /country/namebyip/{ipAddress} `
* **`GET`** : Given the IP returns a JSON object with the country name for that IP
