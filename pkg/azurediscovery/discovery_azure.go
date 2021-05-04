package azurediscovery

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-03-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"
	log "github.com/sirupsen/logrus"
)

type DiscoveryClientAZURE struct {
	TagPrefix    bool
	Tag          string
	Subscription string
}

func (d DiscoveryClientAZURE) GetInstances() ([]discovery.Instance, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	vmssinstanceClient := compute.NewVirtualMachineScaleSetVMsClient(d.Subscription)
	vmssinstanceClient.Authorizer = authorizer
	vmssClient := compute.NewVirtualMachineScaleSetsClient(d.Subscription)
	vmssClient.Authorizer = authorizer
	interfacesClient := network.NewInterfacesClient(d.Subscription)
	interfacesClient.Authorizer = authorizer
	interfaceIPConfigurationClient := network.NewInterfaceIPConfigurationsClient(d.Subscription)
	interfaceIPConfigurationClient.Authorizer = authorizer
	targetInstances := []discovery.Instance{}

	for vmssList, err := vmssClient.ListAll(context.Background()); vmssList.NotDone(); err = vmssList.Next() {
		if err != nil {
			log.Error(err)
			return nil, err
		}

		for _, vmss := range vmssList.Values() {
			var metrics = []discovery.InstanceMetrics{}
			if val, ok := vmss.Tags[d.Tag]; ok {
				if !ok {
					continue
				}
				json.Unmarshal([]byte(*val), &metrics)
			}

			log.Debugf("Found VMSS: %v\n", *vmss.Name)
			rg := strings.Split(*vmss.ID, "/")[4]
			log.Debugf("VMSS RG: %v\n", rg)
			//InterfaceIPConfiguration InterfaceIPConfigurationPropertiesFormat
			for interfaceList, err := interfacesClient.ListVirtualMachineScaleSetNetworkInterfaces(context.Background(), rg, *vmss.Name); interfaceList.NotDone(); err = interfaceList.Next() {
				if err != nil {
					log.Error(err)
					return nil, err
				}

				for _, interfaceValue := range interfaceList.Values() {
					//IPConfigurationPropertiesFormat
					log.Debugf("Interface IPConfigs %v", &interfaceValue)
					tagsBytes, err := json.Marshal(&vmss.Tags)
					if err != nil {
						log.Error(err)
						return nil, err
					}
					var tags map[string]string
					err = json.Unmarshal(tagsBytes, &tags)
					if err != nil {
						log.Error(err)
						return nil, err
					}
					log.Debugf("Tags: %v", tags)

					for _, ipConfigValue := range *interfaceValue.IPConfigurations {
						log.Debugf("ipConfiguration: %v", *ipConfigValue.PrivateIPAddress)
						targetInstance := discovery.Instance{
							PrivateIP:    *ipConfigValue.PrivateIPAddress,
							InstanceType: *vmss.Name,
							Tags:         tags,
							Metrics:      metrics,
						}
						targetInstances = append(targetInstances, targetInstance)
					}

				}

			}
		}
	}

	return targetInstances, nil
}
