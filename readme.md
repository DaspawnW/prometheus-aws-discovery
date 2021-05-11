# Prometheus AWS/Azure Tag based Discovery

Prometheus AWS/Azure Discovery provides a more flexible way to discover exporters running on ec2 instances or withing VMSS Instances (Azure VMs currently not covered by binary).

### AUTH

AWS Auth is done via the AWS go sdk. Meaning it shoud support ENV Vars & Instance Profiles by default [AWS Docs](https://docs.aws.amazon.com/sdk-for-go/api/aws/session/)

Azure Auth requires providing ENV Vars. For supported combinations see [Azure Docs](https://github.com/Azure/azure-sdk-for-go#more-authentication-details)

### Discovery Tag

The normal Prometheus ec2_sd_configs object is very limited and doesn't support extended discovery. This limits it for example to the specified list of ports inside the prometheus configuration file.

To overcome this issue this project allows discovery based on EC2/VMSS tags in a format that the file_sd_config discovery function can read.

IMPORTANT: V2 Tags for discovery 
if the flag --tag is provided, the binary will look for the V2 Tag Structure.
```
key: example
value: [{"name": "ref-test","port": 9100,"path": "metrics","scheme": "https"},{...},{...}]
```

IMPORTANT: The below applies to v1 tags. These will be deprecated soon.
if the flag `--tag-prefix` is provided, the binary will currently look for v1 tags with below structure. 
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
