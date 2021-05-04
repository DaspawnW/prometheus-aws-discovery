package main

import (
	"encoding/csv"
	"flag"
	"os"
	"strings"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/awsdiscovery"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/azurediscovery"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/output"

	log "github.com/sirupsen/logrus"
)

var Version = "development"

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	var outputType string
	var iaasCSV string
	tagPrefix := false
	var tag string
	var filePath string
	var kubeconfig string
	var namespace string
	var configmapName string
	var configmapKey string
	var subscrID string

	flag.StringVar(&tag, "tagprefix", "prom/scrape", "Tag-Key-Prefix used for exporter config (V1-Tag)")
	flag.StringVar(&tag, "tag", "", "Tag-Key used to look for exporter data (V2-Tag)")
	flag.StringVar(&outputType, "output", "kubernetes", "Allowed Values {kubernetes|file}")
	flag.StringVar(&filePath, "file-path", "", "Target file path for sd_config file output")
	flag.StringVar(&kubeconfig, "kube-config", "", "Path to a kubeconfig file")
	flag.StringVar(&namespace, "kube-namespace", "", "Namespace where to create or update configmap (If in cluster and no namespace is provided it tries to detect namespace from incluster config otherwise it uses 'default' namespace)")
	flag.StringVar(&configmapName, "kube-configmap-name", "", "Name of the configmap to create or update with discovery output")
	flag.StringVar(&configmapKey, "kube-configmap-key", "", "Name of configmap key to set discovery output to")
	flag.StringVar(&iaasCSV, "iaas", "aws", "CSV of Clouds to check [aws/azure] e.g. aws,azure ")
	flag.StringVar(&subscrID, "azure-subscr", "", "Azure Subscription ID to look for VMs")
	verbose := flag.Bool("verbose", false, "Print verbose log messages")
	printVersion := flag.Bool("version", false, "Print version")
	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "tagprefix" {
			tagPrefix = true
		}
	})
	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	if *printVersion {
		log.Info("Version: " + Version)
		os.Exit(0)
	}

	log.Info("Infra " + iaasCSV)
	infraReader := csv.NewReader(strings.NewReader(iaasCSV))
	records, err := infraReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}
	var o output.Output
	switch outputType {
	case "stdout":
		log.Info("Configured stdout as output target")
		o = output.OutputStdOut{}
	case "file":
		log.Info("Configured file as output target")
		o = output.OutputFile{FilePath: filePath}
	default:
		log.Info("Configured kubernetes as output target")
		k8s, err := output.NewOutputKubernetes(kubeconfig, namespace, configmapName, configmapKey)
		if err != nil {
			panic(err)
		}
		o = k8s
	}
	validateArg("outputtype", outputType, []string{"kubernetes", "file", "stdout"})

	for _, runInfra := range records[0] {
		log.Info(runInfra)
		clients := []discovery.DiscoveryClient{}
		switch runInfra {
		case "aws":
			log.Info("starting aws discovery")
			clients = append(clients, &awsdiscovery.DiscoveryClientAWS{
				TagPrefix: tagPrefix,
				Tag:       tag,
			})

		case "azure":
			log.Info("starting azure discovery")
			if subscrID == "" {
				log.Errorf("Azure set as target but no subscription provided. Use --azure-subscr")
				os.Exit(1)
			}
			clients = append(clients, azurediscovery.DiscoveryClientAZURE{
				TagPrefix:    tagPrefix,
				Tag:          tag,
				Subscription: subscrID,
			})

		}
		getOutput(clients, o)
	}

}
func validateArg(field string, arg string, allowedValues []string) {
	if !sliceContains(arg, allowedValues) {
		log.Errorf("Field %v has allowed values %v but got %s", field, allowedValues, arg)
		os.Exit(1)
	}
}

func sliceContains(arg string, allowedValues []string) bool {
	for _, v := range allowedValues {
		if v == arg {
			return true
		}
	}
	return false
}
func getOutput(clients []discovery.DiscoveryClient, output output.Output) {
	outputInstances := []discovery.Instance{}
	for _, d := range clients {
		instances, err := d.GetInstances()
		if err != nil {
			log.Error(err)
		}
		outputInstances = append(outputInstances, instances...)
	}
	log.Debug("Writing output\n")
	err := output.Write(outputInstances)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Wrote output to target")
	log.Info("Completed successfully")

}
