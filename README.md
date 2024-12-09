# Azure Footprint Collector Code Tour

This Go program allows you to explore Azure resources by listing available subscriptions, selecting resource groups within those subscriptions, and viewing details for virtual machines within the selected resource groups. The program utilizes the Azure SDK for Go and demonstrates an interactive CLI tool for navigating Azure resources.

## Prerequisites

- **Go**: Ensure you have Go installed. [Install Go](https://golang.org/doc/install) if you haven't.
- **Azure Subscription**: You need an Azure account with access to at least one subscription.
- **Azure CLI**: Although not required, you may want to have Azure CLI installed and authenticated for easier troubleshooting.

## Installation

Clone this repository and navigate to the project directory:

```bash
git clone <repository-url>
cd <project-directory>
```

Install the required Go packages:

```bash
go mod tidy
```

## Usage

Run the program:

```bash
go run azure_footprint.go
```

The program will guide you through the following steps:

1. **Select an Azure Subscription**: Lists available subscriptions and prompts you to select one by entering its corresponding number.
2. **Select a Resource Group**: Lists resource groups within the selected subscription and prompts for selection.
3. **Select a Virtual Machine**: Lists VMs in the selected resource group and prompts for selection.
4. **View VM Details**: Displays details of the selected VM, including its name, location, and size.

Example output:

```plaintext
Available Subscriptions:
[0] Azure subscription 1 (64d93d73-b769-41ef-879b-e587f5f86f6c)
Select a subscription by number: 0
Selected subscription: Azure subscription 1 (64d93d73-b769-41ef-879b-e587f5f86f6c)

Available Resource Groups:
[0] resourceGroup1
[1] resourceGroup2
Select a resource group by number: 0
Selected resource group: resourceGroup1

Available VMs:
[0] MyVM1
[1] MyVM2
Select a VM by number: 1
Selected VM: MyVM2

VM Information:
Name: MyVM2
Location: westus
Size: Standard_B2s
```

## Code Walkthrough

### Main Function

```go
func main() {
	// Initialize Azure credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Failed to create credential: %v", err)
	}
	ctx := context.Background()

	// List subscriptions and prompt for selection
	subscriptions, err := listSubscriptions(ctx, cred)
	...
}
```

The `main` function initializes the Azure credentials using `azidentity.NewDefaultAzureCredential`, which automatically handles authentication using environment variables, Azure CLI credentials, or managed identities.

### List Subscriptions

```go
func listSubscriptions(ctx context.Context, cred *azidentity.DefaultAzureCredential) ([]*armsubscriptions.Subscription, error) {
	client, err := armsubscriptions.NewClient(cred, nil)
	...
}
```

The `listSubscriptions` function uses the `armsubscriptions.NewClient` to create a subscription client and retrieve all subscriptions associated with the authenticated account. It returns a slice of `Subscription` pointers, which includes details like `DisplayName` and `SubscriptionID`.

### Select a Subscription

In `main`, the program prompts the user to select a subscription by displaying a numbered list. The chosen subscription is used for all subsequent resource queries.

```go
for i, sub := range subscriptions {
	fmt.Printf("[%d] %s (%s)\n", i, *sub.DisplayName, *sub.SubscriptionID)
}
fmt.Print("Select a subscription by number: ")
var subIndex int
fmt.Scanln(&subIndex)
```

### List Resource Groups

```go
func listResourceGroups(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID string) ([]*armresources.ResourceGroup, error) {
	rgClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	...
}
```

The `listResourceGroups` function creates a `ResourceGroupsClient` using the selected `subscriptionID`. It retrieves all resource groups within the subscription, and the main function then prompts the user to choose one.

### List Virtual Machines in the Resource Group

```go
func listVMs(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID, resourceGroup string) ([]*armcompute.VirtualMachine, error) {
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	...
}
```

The `listVMs` function creates a `VirtualMachinesClient` and retrieves all VMs within the specified resource group. The main function then lists these VMs for the user to select.

### Retrieve VM Details

```go
func getVMInfo(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID, resourceGroup, vmName string) (*armcompute.VirtualMachine, error) {
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	...
}
```

The `getVMInfo` function fetches detailed information about a selected VM using its name and resource group. It returns the `VirtualMachine` struct, containing attributes such as `Name`, `Location`, and `HardwareProfile.VMSize`.

### Error Handling

Each function handles errors by returning them to the caller. The `main` function logs errors and exits on failure, ensuring that errors are visible and informative.

### Dependencies

This project uses the Azure SDK for Go and its relevant packages:
- `azidentity` for authentication
- `armsubscriptions` for subscriptions management
- `armresources` for resource group management
- `armcompute` for virtual machine management


# Dockerfile

If you want to pull the already-built container, you can find it on Docker Hub at [anuclei/azure_footprint](https://hub.docker.com/r/anuclei/azure_footprint).

## Next Steps

This project provides a basic CLI tool for exploring Azure resources. You can extend it by adding more features, such as listing other types of resources, starting/stopping VMs, or deploying new resources.

## Resources

- [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```