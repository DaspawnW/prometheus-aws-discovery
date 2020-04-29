package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/output"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/outputfile"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/outputkubernetes"
	"github.com/google/logger"
	log "github.com/sirupsen/logrus"
)

var Version = "development"

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	tagPrefix := flag.String("tagprefix", "prom/scrape", "Prefix used for tag key to filter for exporter")
	outputType := flag.String("output", "kubernetes", "Allowed Values {kubernetes|file}")
	filePath := flag.String("file-path", "", "Target file path to write file to")
	kubeconfig := flag.String("kube-config", "", "Path to a kubeconfig file")
	namespace := flag.String("kube-namespace", "", "Namespace where to create or update configmap (If in cluster and no namespace is provided it tries to detect namespace from incluster config otherwise it uses 'default' namespace)")
	configmapName := flag.String("kube-configmap-name", "", "Name of the configmap to create or update with discovery output")
	configmapKey := flag.String("kube-configmap-key", "", "Name of configmap key to set discovery output to")
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

	validateArg("outputtype", *outputType, []string{"kubernetes", "file"})

	awsSession := session.New()
	awsConfig := &aws.Config{}
	ec2Client := ec2.New(awsSession, awsConfig)

	log.Info("Start discovery of ec2 instances")
	d := discovery.NewDiscovery(ec2Client, *tagPrefix)
	instances, err := d.Discover()
	if err != nil {
		panic(err)
	}

	var o output.Output
	if *outputType == "kubernetes" {
		log.Info("Configured kubernetes as output target")
		k8s, err := outputkubernetes.NewOutputKubernetes(*kubeconfig, *namespace, *configmapName, *configmapKey)
		if err != nil {
			panic(err)
		}

		o = k8s
	} else {
		log.Info("Configured file as output target")
		o = outputfile.OutputFile{FilePath: *filePath}
	}
	e := o.Write(*instances)
	if e != nil {
		panic(e)
	}

	log.Info("Wrote output to target")
	log.Info("Completed successfully")
}

func validateArg(field string, arg string, allowedValues []string) {
	if !sliceContains(arg, allowedValues) {
		logger.Error(fmt.Sprintf("Field %v has allowed values %v but got %s", field, allowedValues, arg))
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
