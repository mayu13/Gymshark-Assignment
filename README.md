# Gymshark Assignment 

A HTTP web server hosts a REST API for calculating the number of packs required to fulfill an order. The server is designed to handle various pack sizes, which are stored in memory. These pack sizes represent the available options for packaging orders.

Consists of 2 Endpoints:

* `POST /packs` - used to set pack sizes
* `POST /calculate` - based on the order calculates the optimal pack combination to fullfill the order


### Local

You can run unit tests by:

```
make test
```

---
You can run the service locally by:
```
make run
```
---

```
To run in local (browser):

http://localhost:8009

```

Example calculate request:

```
curl -v --header "Content-Type: application/json" --request POST --data '{"items_count":250}' http://127.0.0.1:9000/calculate 
```


Example set pack sizes request:

```
curl -v --header "Content-Type: application/json" --request POST --data '{"sizes": [31, 23, 53]}' http://127.0.0.1:9000/packs

```

