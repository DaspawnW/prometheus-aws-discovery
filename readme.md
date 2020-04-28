# Prometheus AWS Discovery

Prometheus AWS Discovery provides a more flexible way to discover exporters running on ec2 instances.

The normal Prometheus ec2_sd_configs object is very limited and doesn't support extended discovery. This limits it for example to the specified list of ports inside the prometheus configuration file.

To overcome this issue this project allows discovery based on EC2 instance tags in a format that the file_sd_config discovery function can read.

A tag must have the following format:

```
Key: <prefix>:<port>/<path>
Value: <name>
```

* prefix: Can be configured as command line argument, by default it is "prom/scrape"
* port: Describes the port where the application or exporter are listening on
* path: Path to the metrics provided by an exporter or an application (e.g. /metrics)
* name: Is a name that is added as "label" to the metrics

Here a simple example:

```
Key: prom/scrape:9100/metrics
Value: node_exporter
```

This will produce the following file_sd_config:

```json
[
  {
    "targets": [
      "<private-ip-address>:<port>"
    ],
    "labels": {
      "__address__": "<private-ip-address>:<port>",
      "__metrics_path__": "/<path>",
      "instancename": "<Name-Tag>",
      "name": "<name>",
      "...": "..."
    }
  }
]
```

where "..." is the list of tags associated with the EC2 Instance, all AWS specific tags starting with "aws:" are removed from the list of tags.