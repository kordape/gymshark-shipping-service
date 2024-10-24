# Gymshark Shipping Service

Shipping Service is a HTTP web server exposing REST API to handle calculation of packs needed to send in order to fulfill the order.

Uses a list of `packSizes` (stored in memory) which represent the available choice of packs for calculations.

Consists of 2 Endpoints:

* `POST /packs` - used to set pack sizes
* `POST /calculate` - based on the order calculates the optimal pack combination to fullfill the order

### Live Demo

You can access the live demo of the app [here](http://gymshark-shipping-calculator.s3-website.eu-central-1.amazonaws.com/)

If you have some trouble accessing it please contact: petar.korda@gmail.com

### Local Development

You can run unit tests by:

```
make test
```

---
You can run the service locally by:
```
make local-run
```

This will spin up a container with a web server listening on port 80 (can be changed by changing `docker-compose.yml` file)

---

Example calculate request:

```
curl -v --header "Content-Type: application/json" --request POST --data '{"items_count":250}' http://127.0.0.1/calculate 
```


Example set pack sizes request:

```
curl -v --header "Content-Type: application/json" --request POST --data '{"sizes": [31, 23, 53]}' http://127.0.0.1/packs
```
