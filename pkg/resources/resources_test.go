package resources

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func authenticate() (*Azure, context.Context, *azidentity.DefaultAzureCredential) {

	azure := &Azure{}

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

	return azure, ctx, cred
}

func TestRegionDeployment(t *testing.T) {

	azure, ctx, cred := authenticate()

	defer Cleanup(ctx, cred, *azure)

	testRegionIn := func(region string) func() error {
		return func() error {
			_, err := CreateResourceGroup(ctx, cred, *azure, region)
			return err
		}
	}

	testCases := map[string]struct {
		test             func() error
		shouldError      bool
		linkedPolicyName string
	}{
		"ResourceGroupCreation -> Brazil Southeast": {
			test:             testRegionIn("Brazil Southeast"),
			shouldError:      true,
			linkedPolicyName: "denied regions",
		},
		"ResourceGroupCreation -> Switzerland North": {
			test:        testRegionIn("Switzerland North"),
			shouldError: false,
		},
		"ResourceGroupCreation -> Switzerland West": {
			test:        testRegionIn("Switzerland West"),
			shouldError: false,
		},
		"ResourceGroupCreation -> Central US": {
			test:        testRegionIn("Central US"),
			shouldError: false,
		},
		"ResourceGroupCreation -> East US": {
			test:        testRegionIn("East US"),
			shouldError: false,
		},
		"ResourceGroupCreation -> East US 2": {
			test:        testRegionIn("East US 2"),
			shouldError: false,
		},
		"ResourceGroupCreation -> North Central US": {
			test:        testRegionIn("North Central US"),
			shouldError: false,
		},
		"ResourceGroupCreation -> South Central US": {
			test:        testRegionIn("South Central US"),
			shouldError: false,
		},
		"ResourceGroupCreation -> West Central US": {
			test:        testRegionIn("West Central US"),
			shouldError: false,
		},
		"ResourceGroupCreation -> West US": {
			test:        testRegionIn("West US"),
			shouldError: false,
		},
		"ResourceGroupCreation -> West US 2": {
			test:        testRegionIn("West US 2"),
			shouldError: false,
		},
		"ResourceGroupCreation -> West US 3": {
			test:        testRegionIn("West US 3"),
			shouldError: false,
		},
	}

	checkPolicy := func(err error, expectedPolicy string) error {
		if expectedPolicy == "" {
			return nil
		}

		azerr, ok := err.(*AzureError)

		if len(azerr.Response.AdditionalInfo) != 0 {

			errPolicy := azerr.Response.AdditionalInfo[0].Info.PolicyAssignmentName
			if !ok {
				return fmt.Errorf("Deployment was denied but we couldn't evaluate the policy violation: %s", err)
			}

			if errPolicy != expectedPolicy {
				return fmt.Errorf("Policy assignment failed but not with the right policy: have %s, want %s", errPolicy, expectedPolicy)
			}
		}
		return nil
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			Cleanup(ctx, cred, *azure)

			err := testCase.test()

			switch {
			case err == nil && testCase.shouldError:
				t.Fatalf("Deployment was allowed but it should have been denied: %s", err)

			case err != nil && testCase.shouldError:
				if err := checkPolicy(err, testCase.linkedPolicyName); err != nil {
					t.Fatal(err)
				}

			case err != nil && !testCase.shouldError:
				t.Fatalf("Deployment was denied but it should have been allowed: %s", err)
			}
		})
	}
}

func TestNetworkSecurityGroupPolicies(t *testing.T) {
	azure, ctx, cred := authenticate()

	region := "Switzerland West"

	defer Cleanup(ctx, cred, *azure)

	_, err := CreateResourceGroup(ctx, cred, *azure, region)
	if err != nil {
		log.Fatal(err)
	}

	_, err = CreateNetworkSecurityGroup(ctx, cred, *azure, "192.168.1.1/32", "443", region)
	if err != nil {
		log.Fatal()
	}

	testSecurityGroupOn := func(prefix string, port string) func() error {
		return func() error {
			_, err := CreateNetworkSecurityGroup(ctx, cred, *azure, prefix, port, region)
			return err
		}
	}

	testSecurityRuleOn := func(prefix string, port string) func() error {
		return func() error {
			_, err := CreateNetworkSecurityRule(ctx, cred, *azure, prefix, port)
			return err
		}
	}

	testCases := map[string]struct {
		test             func() error
		shouldError      bool
		linkedPolicyName string
	}{
		"NetworkSecurityGroupRule SourceAddressPrefix -> * on port -> *": {
			test:             testSecurityGroupOn("*", "*"),
			shouldError:      true,
			linkedPolicyName: "nsg deny any any rule",
		},
		"NetworkSecurityGroupRule SourceAddressPrefix -> * on port -> 22": {
			test:             testSecurityGroupOn("*", "22"),
			shouldError:      true,
			linkedPolicyName: "nsg deny any ssh rule",
		},
		"NetworkSecurityGroupRule SourceAddressPrefix -> * on port -> 3389": {
			test:             testSecurityGroupOn("*", "3389"),
			shouldError:      true,
			linkedPolicyName: "nsg deny any rdp rule",
		},
		"NetworkSecurityGroupRule SourceAddressPrefix -> internet on port -> 22": {
			test:             testSecurityGroupOn("internet", "22"),
			shouldError:      true,
			linkedPolicyName: "nsg deny any ssh rule",
		},
		"NetworkSecurityGroupRule SourceAddressPrefix -> internet on port -> 3389": {
			test:             testSecurityGroupOn("internet", "3389"),
			shouldError:      true,
			linkedPolicyName: "nsg deny any rdp rule",
		},
		"NetworkSecurityGroupRule SourceAddressPrefix -> 0.0.0.0/0 on port -> 22": {
			test:             testSecurityGroupOn("0.0.0.0/0", "22"),
			shouldError:      true,
			linkedPolicyName: "nsg deny any ssh rule",
		},
		"NetworkSecurityGroupRule SourceAddressPrefix -> 0.0.0.0/0 on port -> 3389": {
			test:             testSecurityGroupOn("0.0.0.0/0", "3389"),
			shouldError:      true,
			linkedPolicyName: "nsg deny any rdp rule",
		},
		"NetworkSecurityGroupRule SourceAddressPrefix -> 192.168.0.0/24 on port -> 22": {
			test:        testSecurityGroupOn("192.168.0.0/24", "22"),
			shouldError: false,
		},
		"NetworkSecurityGroupRule SourceAddressPrefix -> 192.168.0.0/24 on port -> 3389": {
			test:        testSecurityGroupOn("192.168.0.0/24", "3389"),
			shouldError: false,
		},
		"NetworkSecurityRule SourceAddressPrefix -> * on port -> *": {
			test:             testSecurityRuleOn("*", "*"),
			shouldError:      true,
			linkedPolicyName: "nsr deny any any rule",
		},
		"NetworkSecurityRule SourceAddressPrefix -> * on port -> 22": {
			test:             testSecurityRuleOn("*", "22"),
			shouldError:      true,
			linkedPolicyName: "nsr deny any ssh rule",
		},
		"NetworkSecurityRule SourceAddressPrefix -> * on port -> 3389": {
			test:             testSecurityRuleOn("*", "3389"),
			shouldError:      true,
			linkedPolicyName: "nsr deny any rdp rule",
		},
		"NetworkSecurityRule SourceAddressPrefix -> internet on port -> 22": {
			test:             testSecurityRuleOn("internet", "22"),
			shouldError:      true,
			linkedPolicyName: "nsr deny any ssh rule",
		},
		"NetworkSecurityRule SourceAddressPrefix -> internet on port -> 3389": {
			test:             testSecurityRuleOn("internet", "3389"),
			shouldError:      true,
			linkedPolicyName: "nsr deny any rdp rule",
		},
		"NetworkSecurityRule SourceAddressPrefix -> 0.0.0.0/0 on port -> 22": {
			test:             testSecurityRuleOn("0.0.0.0/0", "22"),
			shouldError:      true,
			linkedPolicyName: "nsr deny any ssh rule",
		},
		"NetworkSecurityRule SourceAddressPrefix -> 0.0.0.0/0 on port -> 3389": {
			test:             testSecurityRuleOn("0.0.0.0/0", "3389"),
			shouldError:      true,
			linkedPolicyName: "nsr deny any rdp rule",
		},
		"NetworkSecurityRule SourceAddressPrefix -> 192.168.0.0/24 on port -> 22": {
			test:        testSecurityRuleOn("192.168.0.0/24", "22"),
			shouldError: false,
		},
		"NetworkSecurityRule SourceAddressPrefix -> 192.168.0.0/24 on port -> 3389": {
			test:        testSecurityRuleOn("192.168.0.0/24", "3389"),
			shouldError: false,
		},
	}

	checkPolicy := func(err error, expectedPolicy string) error {
		if expectedPolicy == "" {
			return nil
		}

		azerr, ok := err.(*AzureError)

		if len(azerr.Response.AdditionalInfo) != 0 {
			errPolicy := azerr.Response.AdditionalInfo[0].Info.PolicyAssignmentName

			if !ok {
				return fmt.Errorf("Deployment was denied but we couldn't evaluate policy violation: %s", err)
			}

			if errPolicy != expectedPolicy {
				return fmt.Errorf("Policy assignment failed but not with the right policy: have %s, want %s", errPolicy, expectedPolicy)
			}
		}
		return nil
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := testCase.test()

			switch {
			case err == nil && testCase.shouldError:
				t.Fatalf("Deployment was allowed but it should have been denied: %s", err)

			case err != nil && testCase.shouldError:
				if err := checkPolicy(err, testCase.linkedPolicyName); err != nil {
					t.Fatal(err)
				}

			case err != nil && !testCase.shouldError:
				t.Fatalf("Deployment was denied but it should have been allowed: %s", err)
			}
		})
	}
}
