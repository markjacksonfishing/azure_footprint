package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

func main() {
	// Create a new Azure credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Failed to create credential: %v", err)
	}

	ctx := context.Background()

	// List available subscriptions
	subscriptions, err := listSubscriptions(ctx, cred)
	if err != nil {
		log.Fatalf("Failed to list subscriptions: %v", err)
	}

	// Prompt user to choose a subscription
	fmt.Println("Available Subscriptions:")
	for i, sub := range subscriptions {
		fmt.Printf("[%d] %s (%s)\n", i, *sub.DisplayName, *sub.SubscriptionID)
	}
	fmt.Print("Select a subscription by number: ")
	var subIndex int
	fmt.Scanln(&subIndex)

	selectedSub := subscriptions[subIndex]
	fmt.Printf("Selected subscription: %s (%s)\n", *selectedSub.DisplayName, *selectedSub.SubscriptionID)

	// List resource groups in the selected subscription
	resourceGroups, err := listResourceGroups(ctx, cred, *selectedSub.SubscriptionID)
	if err != nil {
		log.Fatalf("Failed to list resource groups: %v", err)
	}

	fmt.Println("Available Resource Groups:")
	for i, rg := range resourceGroups {
		fmt.Printf("[%d] %s\n", i, *rg.Name)
	}
	fmt.Print("Select a resource group by number: ")
	var rgIndex int
	fmt.Scanln(&rgIndex)

	selectedResourceGroup := resourceGroups[rgIndex]
	fmt.Printf("Selected resource group: %s\n", *selectedResourceGroup.Name)

	// List VMs in the selected resource group
	vms, err := listVMs(ctx, cred, *selectedSub.SubscriptionID, *selectedResourceGroup.Name)
	if err != nil {
		log.Fatalf("Failed to list VMs: %v", err)
	}

	fmt.Println("Available VMs:")
	for i, vm := range vms {
		fmt.Printf("[%d] %s\n", i, *vm.Name)
	}
	fmt.Print("Select a VM by number: ")
	var vmIndex int
	fmt.Scanln(&vmIndex)

	selectedVM := vms[vmIndex]
	fmt.Printf("Selected VM: %s\n", *selectedVM.Name)

	// Get VM information
	vmInfo, err := getVMInfo(ctx, cred, *selectedSub.SubscriptionID, *selectedResourceGroup.Name, *selectedVM.Name)
	if err != nil {
		log.Fatalf("Failed to get VM information: %v", err)
	}

	fmt.Printf("\nVM Information:\nName: %s\nLocation: %s\nSize: %s\n",
		*vmInfo.Name, *vmInfo.Location, *vmInfo.Properties.HardwareProfile.VMSize)
}

func listSubscriptions(ctx context.Context, cred *azidentity.DefaultAzureCredential) ([]*armsubscriptions.Subscription, error) {
	client, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		return nil, err
	}

	pager := client.NewListPager(nil)
	var subscriptions []*armsubscriptions.Subscription
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, resp.Value...)
	}
	return subscriptions, nil
}

func listResourceGroups(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID string) ([]*armresources.ResourceGroup, error) {
	// Create a ResourceGroupsClient to list resource groups
	rgClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	pager := rgClient.NewListPager(nil)
	var resourceGroups []*armresources.ResourceGroup
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		resourceGroups = append(resourceGroups, resp.Value...)
	}
	return resourceGroups, nil
}

func listVMs(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID, resourceGroup string) ([]*armcompute.VirtualMachine, error) {
	// Create a VirtualMachinesClient to list VMs
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	pager := vmClient.NewListPager(resourceGroup, nil)
	var vms []*armcompute.VirtualMachine
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		vms = append(vms, resp.Value...)
	}
	return vms, nil
}

func getVMInfo(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID, resourceGroup, vmName string) (*armcompute.VirtualMachine, error) {
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	vm, err := vmClient.Get(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return nil, err
	}

	return &vm.VirtualMachine, nil
}
