# Nature Remo API Client for Go [![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godoc] [![Travis](https://img.shields.io/travis/tenntenn/natureremo.svg?style=flat-square)][travis] [![Go Report Card](https://goreportcard.com/badge/github.com/tenntenn/natureremo)](https://goreportcard.com/report/github.com/tenntenn/natureremo) [![codecov](https://codecov.io/gh/tenntenn/natureremo/branch/master/graph/badge.svg)](https://codecov.io/gh/tenntenn/natureremo)

[godoc]: http://godoc.org/github.com/tenntenn/natureremo
[travis]: https://travis-ci.org/tenntenn/natureremo

`tenntenn/natureremo` is [Nature Remo API](https://developer.nature.global/en/overview/) Client for Go.
[Nature Remo](https://nature.global/en/top) is a smart remote control that easily realizes smart home by connecting your appliances to the Internet.

## Install

```
$ go get -u github.com/tenntenn/natureremo
```

## Examples

See `_example` directory.

```go
func main() {
	cli := natureremo.NewClient(os.Args[1])
	ctx := context.Background()

	applianceName := os.Args[2]
	signalName := os.Args[3]

	as, err := cli.ApplianceService.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var target *natureremo.Appliance
	for _, a := range as {
		if a.Nickname == applianceName {
			target = a
			break
		}
	}

	if target == nil {
		log.Fatalf("%s not found", applianceName)
	}

	for _, s := range target.Signals {
		if s.Name == signalName {
			cli.SignalService.Send(ctx, s)
			break
		}
	}
}
```

## Supported API

### Cloud API

http://swagger.nature.global

|     Status       |                 Endpoint                | HTTP Method |    Service     |
|:----------------:|-----------------------------------------|:-----------:|----------------|
|:heavy_check_mark:|/1/users/me                              | GET         |[UserService](https://godoc.org/github.com/tenntenn/natureremo#UserService)     |
|:heavy_check_mark:|/1/users/me                              | POST        |[UserService](https://godoc.org/github.com/tenntenn/natureremo#UserService)     |
|:heavy_check_mark:|/1/devices                               | GET         |[DeviceService](https://godoc.org/github.com/tenntenn/natureremo#DeviceService)   |
|:heavy_check_mark:|/1/devices/{device}                      | POST        |[DeviceService](https://godoc.org/github.com/tenntenn/natureremo#DeviceService)   |
|:heavy_check_mark:|/1/devices/{device}/delete               | POST        |[DeviceService](https://godoc.org/github.com/tenntenn/natureremo#DeviceService)   |
|:heavy_check_mark:|/1/devices/{device}/temperature_offset   | POST        |[DeviceService](https://godoc.org/github.com/tenntenn/natureremo#DeviceService)   |
|:heavy_check_mark:|/1/devices/{device}/humidity_offset      | POST        |[DeviceService](https://godoc.org/github.com/tenntenn/natureremo#DeviceService)   |
|:heavy_check_mark:|/1/detectappliance                       | POST        |[ApplianceService](https://godoc.org/github.com/tenntenn/natureremo#ApplianceService)|
|:heavy_check_mark:|/1/appliances                            | GET         |[ApplianceService](https://godoc.org/github.com/tenntenn/natureremo#ApplianceService)|
|:heavy_check_mark:|/1/appliances                            | POST        |[ApplianceService](https://godoc.org/github.com/tenntenn/natureremo#ApplianceService)|
|:heavy_check_mark:|/1/appliance_orders                      | POST        |[ApplianceService](https://godoc.org/github.com/tenntenn/natureremo#ApplianceService)|
|:heavy_check_mark:|/1/appliances/{appliance}/delete         | POST        |[ApplianceService](https://godoc.org/github.com/tenntenn/natureremo#ApplianceService)|
|:heavy_check_mark:|/1/appliances/{appliance}                | POST        |[ApplianceService](https://godoc.org/github.com/tenntenn/natureremo#ApplianceService)|
|:heavy_check_mark:|/1/appliances/{appliance}/aircon_settings| POST        |[ApplianceService](https://godoc.org/github.com/tenntenn/natureremo#ApplianceService)|
|                  |/1/appliances/{appliance}/tv             | POST        |                                                                                     |
|                  |/1/appliances/{appliance}/light          | POST        |                                                                                     |
|:heavy_check_mark:|/1/appliances/{appliance}/signals        | GET         |[SignalService](https://godoc.org/github.com/tenntenn/natureremo#SignalService)   |
|:heavy_check_mark:|/1/appliances/{appliance}/signals        | POST        |[SignalService](https://godoc.org/github.com/tenntenn/natureremo#SignalService)   |
|:heavy_check_mark:|/1/appliances/{appliance}/signal_orders  | POST        |[SignalService](https://godoc.org/github.com/tenntenn/natureremo#SignalService)   |
|:heavy_check_mark:|/1/signals/{signal}                      | POST        |[SignalService](https://godoc.org/github.com/tenntenn/natureremo#SignalService)   |
|:heavy_check_mark:|/1/signals/{signal}/delete               | POST        |[SignalService](https://godoc.org/github.com/tenntenn/natureremo#SignalService)   |
|:heavy_check_mark:|/1/signals/{signal}/send                 | POST        |[SignalService](https://godoc.org/github.com/tenntenn/natureremo#SignalService)   |

### Local API

http://local.swagger.nature.global/

|     Status       |Endpoint | HTTP Method |LocalClient Method|
|:----------------:|---------|:-----------:|------------------|
|:heavy_check_mark:|/messages| GET         |[Fetch](https://godoc.org/github.com/tenntenn/natureremo#LocalClient.Fetch)|
|:heavy_check_mark:|/messages| POST        |[Emit](https://godoc.org/github.com/tenntenn/natureremo#LocalClient.Emit) |

