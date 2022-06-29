package main

import (
	"context"
	"log"
	"os"

	"azpolicy-checker/pkg/resources"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {

	azure := &resources.Azure{}
	location := "West Europe"

	azure.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	azure.ResourceGroupName = "daily_policies_tests"
	azure.SecurityGroupName = "daily_policies_tests"

	if len(azure.SubscriptionID) == 0 {
		log.Fatal("AZURE_SUBSCRIPTION_ID is not set.")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	resourceGroup, err := resources.CreateResourceGroup(ctx, cred, *azure, location)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("resources group:", *resourceGroup.ID)
	_, err = resources.CreateNetworkSecurityGroup(ctx, cred, *azure, "internet", "22", location)
	if err != nil {
		log.Fatal(err)
	}

	keepResource := os.Getenv("KEEP_RESOURCE")
	if len(keepResource) == 0 {
		err = resources.Cleanup(ctx, cred, *azure)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("cleaned up successfully.")
	}
}
