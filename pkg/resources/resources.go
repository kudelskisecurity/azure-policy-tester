package resources

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type Azure struct {
	SubscriptionID    string
	ResourceGroupName string
	SecurityGroupName string
}

func CreateNetworkSecurityGroup(ctx context.Context, cred azcore.TokenCredential, azure Azure, sourceAddressPrefix string, destinationPortRange string, region string) (*armnetwork.SecurityGroup, error) {
	networkSecurityGroupClient, err := armnetwork.NewSecurityGroupsClient(azure.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	pollerResp, err := networkSecurityGroupClient.BeginCreateOrUpdate(
		ctx,
		azure.ResourceGroupName,
		azure.SecurityGroupName,
		armnetwork.SecurityGroup{
			Location: &region,
			Properties: &armnetwork.SecurityGroupPropertiesFormat{
				SecurityRules: []*armnetwork.SecurityRule{
					{
						Name: to.Ptr("allow_ssh"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
							SourceAddressPrefix:      to.Ptr(sourceAddressPrefix),
							SourcePortRange:          to.Ptr("1-65535"),
							DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
							DestinationPortRange:     to.Ptr(destinationPortRange),
							Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
							Direction:                to.Ptr(armnetwork.SecurityRuleDirectionInbound),
							Priority:                 to.Ptr[int32](100),
						},
					},
				},
			},
		},
		nil)
	if err != nil {
		return nil, asAzureError(err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.SecurityGroup, nil
}

func CreateNetworkSecurityRule(ctx context.Context, cred azcore.TokenCredential, azure Azure, sourceAddressPrefix string, destinationPortRange string) (*armnetwork.SecurityRule, error) {
	networkSecurityRulesClient, err := armnetwork.NewSecurityRulesClient(azure.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	pollerResp, err := networkSecurityRulesClient.BeginCreateOrUpdate(
		ctx,
		azure.ResourceGroupName,
		azure.SecurityGroupName,
		"allow_ssh",
		armnetwork.SecurityRule{
			Name: to.Ptr("allow_ssh"),
			Properties: &armnetwork.SecurityRulePropertiesFormat{
				Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
				SourceAddressPrefix:      to.Ptr(sourceAddressPrefix),
				SourcePortRange:          to.Ptr("1-65535"),
				DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
				DestinationPortRange:     to.Ptr(destinationPortRange),
				Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
				Direction:                to.Ptr(armnetwork.SecurityRuleDirectionInbound),
				Priority:                 to.Ptr[int32](100),
			},
		}, nil)
	if err != nil {
		return nil, asAzureError(err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.SecurityRule, nil
}

func CreateResourceGroup(ctx context.Context, cred azcore.TokenCredential, azure Azure, region string) (*armresources.ResourceGroup, error) {
	resourceGroupClient, err := armresources.NewResourceGroupsClient(azure.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		azure.ResourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(region),
		},
		nil)
	if err != nil {
		return nil, asAzureError(err)
	}
	return &resourceGroupResp.ResourceGroup, nil
}

func Cleanup(ctx context.Context, cred azcore.TokenCredential, azure Azure) error {
	resourceGroupClient, err := armresources.NewResourceGroupsClient(azure.SubscriptionID, cred, nil)
	if err != nil {
		return err
	}

	pollerResp, err := resourceGroupClient.BeginDelete(ctx, azure.ResourceGroupName, nil)
	if err != nil {
		return err
	}
	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
