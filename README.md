# Nature Remo API Client for Go

## Install

```
$ go get -u github.com/tenntenn/natureremo
```

## Supported API

### Global API

http://swagger.nature.global

|                 Endpoint                | Method |     Status       |
|-----------------------------------------|:------:|:----------------:|
|/1/users/me                              | GET    |:heavy_check_mark:|
|/1/users/me                              | POST   |:heavy_check_mark:|
|/1/devices                               | GET    |:heavy_check_mark:|
|/1/detectappliance                       | POST   |:heavy_check_mark:|
|/1/appliances                            | GET    |:heavy_check_mark:|
|/1/appliances                            | POST   |:heavy_check_mark:|
|/1/appliance_orders                      | POST   |:heavy_check_mark:|
|/1/appliances/{appliance}/delete         | POST   |:heavy_check_mark:|
|/1/appliances/{appliance}                | POST   |:heavy_check_mark:|
|/1/appliances/{appliance}/aircon_settings| POST   |:heavy_check_mark:|
|/1/appliances/{appliance}/signals        | GET    |:heavy_check_mark:|
|/1/appliances/{appliance}/signals        | POST   |:heavy_check_mark:|
|/1/appliances/{appliance}/signal_orders  | POST   |:heavy_check_mark:|
|/1/signals/{signal}                      | POST   |:heavy_check_mark:|
|/1/signals/{signal}/delete               | POST   |:heavy_check_mark:|
|/1/signals/{signal}/send                 | POST   |:heavy_check_mark:|
|/1/devices/{device}                      | POST   |:heavy_check_mark:|
|/1/devices/{device}/delete               | POST   |:heavy_check_mark:|
|/1/devices/{device}/temperature_offset   | POST   |:heavy_check_mark:|
|/1/devices/{device}/humidity_offset      | POST   |:heavy_check_mark:|
