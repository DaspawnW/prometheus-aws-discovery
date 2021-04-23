package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
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
	var tagPrefix string
	var filePath string
	var kubeconfig string
	var namespace string
	var configmapName string
	var configmapKey string
	var subscrID string

	flag.StringVar(&tagPrefix, "tagprefix", "prom/scrape", "Prefix used for tag key to filter for exporter")
	flag.StringVar(&outputType, "output", "kubernetes", "Allowed Values {kubernetes|file}")
	flag.StringVar(&filePath, "file-path", "", "Target file path to write file to")
	flag.StringVar(&kubeconfig, "kube-config", "", "Path to a kubeconfig file")
	flag.StringVar(&namespace, "kube-namespace", "", "Namespace where to create or update configmap (If in cluster and no namespace is provided it tries to detect namespace from incluster config otherwise it uses 'default' namespace)")
	flag.StringVar(&configmapName, "kube-configmap-name", "", "Name of the configmap to create or update with discovery output")
	flag.StringVar(&configmapKey, "kube-configmap-key", "", "Name of configmap key to set discovery output to")
	flag.StringVar(&iaasCSV, "iaas", "", "CSV of Clouds to check [aws/azure]")
	flag.StringVar(&subscrID, "azure-subsc", "", "Azure Subscription ID to look for VMs")
	verbose := flag.Bool("verbose", false, "Print verbose log messages")
	printVersion := flag.Bool("version", false, "Print version")
	flag.Parse()

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

	for _, runInfra := range records[0] {
		log.Info(runInfra)
		switch runInfra {
		case "aws":
			log.Info("starting aws discovery")
			validateArg("outputtype", outputType, []string{"kubernetes", "file", "stdout"})
			awsSession := session.New()
			awsConfig := &aws.Config{}
			ec2Client := ec2.New(awsSession, awsConfig)
			log.Info("Start discovery of ec2 instances")
			d := awsdiscovery.DiscoveryClientAWS{
				Ec2Client: ec2Client,
				TagPrefix: tagPrefix,
			}
			getOutput(d, o)
		case "azure":
			log.Info("starting azure discovery")

			d := azurediscovery.DiscoveryClientAZURE{
				TagPrefix:    tagPrefix,
				Subscription: subscrID,
			}
			getOutput(d, o)
		}
	}

}
func validateArg(field string, arg string, allowedValues []string) {
	if !sliceContains(arg, allowedValues) {
		log.Error(fmt.Sprintf("Field %v has allowed values %v but got %s", field, allowedValues, arg))
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
func getOutput(client discovery.DiscoveryClient, output output.Output) {
	instances, err := client.GetInstances()
	if err != nil {
		log.Error(err)
	}
	log.Debug("Writing output\n")
	e := output.Write(instances)
	if e != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Wrote output to target")
	log.Info("Completed successfully")

}
