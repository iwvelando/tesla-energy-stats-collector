# tesla-energy-stats-collector

## Purpose

tesla-energy-stats-collector is a data ingest tool for Tesla Energy systems (Powerwall and Solar).
This is based on the API documentation in
[vloschiavo/powerwall2](https://github.com/vloschiavo/powerwall2). This is designed to poll a
majority of the data available from the API with a focus on quantitative, non-duplicate metrics. The
complete list of endpoints polled is found in the [Connect.GetAll](/connect/connect.go) function
whose purpose is to perform all individual data gathering tasks.

The Tesla gateway is polled at a frequency set by the configuration. The data is written into
InfluxDB 1.x or 2.x asynchronously, and error handling behavior is defined by the configuration in
which the operator may choose to let an external system such as systemd handle restart behavior.

At this time this was written for my personal use, but I'm open to contributions or feedback if
someone wants to expand the functionality in a backwards-compatible manner.

## Schema

In InfluxDB this code writes to the following measurements:

* energy_configuration
* energy_devices
* energy_faults
* energy_inverters
* energy_meters
* energy_network
* energy_powerwalls

These can all be prefixed based on the configuration.

## References

| Reference | Description |
| --- | --- |
| [vloschiavo/powerwall2](https://github.com/vloschiavo/powerwall2) | Primary API reference |
| [pypowerwall vitals poller](https://github.com/jasonacox/pypowerwall/tree/main/examples/vitals) | Device vitals protobuf handling |
