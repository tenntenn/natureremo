# Nature Remo API Client for Go [![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godoc] [![Travis](https://img.shields.io/travis/tenntenn/natureremo.svg?style=flat-square)][travis] [![Go Report Card](https://goreportcard.com/badge/github.com/tenntenn/natureremo)](https://goreportcard.com/report/github.com/tenntenn/natureremo) [![codecov](https://codecov.io/gh/tenntenn/natureremo/branch/master/graph/badge.svg)](https://codecov.io/gh/tenntenn/natureremo)

[godoc]: http://godoc.org/github.com/tenntenn/natureremo
[travis]: https://travis-ci.org/tenntenn/natureremo

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
