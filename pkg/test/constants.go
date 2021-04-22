package test

var Targets = []map[string]interface{}{
	{
		"targets": []string{
			"127.0.0.1:9100",
		},
		"labels": map[string]string{
			"__address__":      "127.0.0.1:9100",
			"__metrics_path__": "/metrics",
			"__scheme__":       "http",
			"billingnumber":    "1111",
			"instancename":     "Testinstance1",
			"name":             "node_exporter",
		},
	},
	{
		"targets": []string{
			"127.0.0.1:8080",
		},
		"labels": map[string]string{
			"__address__":      "127.0.0.1:8080",
			"__metrics_path__": "/metrics",
			"__scheme__":       "https",
			"billingnumber":    "1111",
			"instancename":     "Testinstance1",
			"name":             "blackbox_exporter",
		},
	},
	{
		"targets": []string{
			"127.0.0.2:9100",
		},
		"labels": map[string]string{
			"__address__":      "127.0.0.2:9100",
			"__metrics_path__": "/metrics",
			"__scheme__":       "http",
			"billingnumber":    "2222",
			"instancename":     "Testinstance2",
			"name":             "node_exporter",
		},
	},
}
var Instances = []map[string]interface{}{
	{
		"InstanceType": "t2.medium",
		"PrivateIP":    "127.0.0.1",
		"Tags": map[string]string{
			"Name":          "Testinstance1",
			"billingnumber": "1111",
		},
		"Metrics": []map[string]interface{}{
			{
				"Name":   "node_exporter",
				"Port":   9100,
				"Path":   "/metrics",
				"Scheme": "http",
			}, {
				"Name":   "blackbox_exporter",
				"Path":   "/metrics",
				"Port":   8080,
				"Scheme": "https",
			},
		},
	}, {
		"InstanceType": "t2.small",
		"PrivateIP":    "127.0.0.2",
		"Tags": map[string]string{
			"Name":          "Testinstance2",
			"billingnumber": "2222",
		},
		"Metrics": []map[string]interface{}{
			{
				"Name":   "node_exporter",
				"Path":   "/metrics",
				"Port":   9100,
				"Scheme": "http",
			},
		},
	},
}
