##### Azure policies
####  * Enforces secure defaults for the resources deployed on Azure
###   * https://docs.microsoft.com/en-us/azure/governance/policy/overview
##    * https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/policy_definition

resource "azurerm_policy_definition" "nsg_deny_inbound_ports" {
  name                = "NSG deny inbound on specific ports"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "Deny Network Security Group Rules where all inbound traffic is allowed on certain ports"
  management_group_id = azurerm_management_group.your_root_managementgroup.id

  metadata = <<METADATA
    {
      "category": "Network"
    }
  METADATA

  parameters = <<PARAMETERS
    {
      "deniedPorts": {
        "type": "Array",
        "metadata": {
          "displayName": "Ports to block",
          "description": "The inbound ports that should be blocked"
        }
      }
    }
  PARAMETERS

  policy_rule = <<POLICY_RULE
    {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.Network/networkSecurityGroups"
          },
          {
            "count": {
              "field": "Microsoft.Network/networkSecurityGroups/securityRules[*]",
              "where": {
                "allOf": [
                  {
                    "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].access",
                    "equals": "Allow"
                  },
                  {
                    "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].direction",
                    "equals": "Inbound"
                  },
                  {
                    "anyOf": [
                      {
                        "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].destinationPortRange",
                        "in": "[parameters('deniedPorts')]"
                      },
                      {
                        "not": {
                          "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].destinationPortRanges[*]",
                          "notIn": "[parameters('deniedPorts')]"
                        }
                      }
                    ]
                  },
                  {
                    "anyOf": [
                      {
                        "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].sourceAddressPrefix",
                        "in": [
                          "*",
                          "0.0.0.0/0",
                          "Internet"
                        ]
                      },
                      {
                        "not": {
                          "field": "Microsoft.Network/networkSecurityGroups/securityRules[*].sourceAddressPrefixes[*]",
                          "notIn": [
                            "*",
                            "0.0.0.0/0",
                            "Internet"
                          ]
                        }
                      }
                    ]
                  }
                ]
              }
            },
            "greaterOrEquals": 1
          }
        ]
      },
      "then": {
        "effect": "deny"
      }
    }
  POLICY_RULE
}

resource "azurerm_policy_definition" "nsr_deny_inbound_ports" {
  name                = "NSR deny inbound on specific ports"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "Deny Network Security Rules where all inbound traffic is allowed on certain ports"
  # management_group_id = azurerm_management_group.your_root_managementgroup.id # <- if you want to attach the policy at the management-group level

  metadata = <<METADATA
    {
      "category": "Network"
    }
  METADATA

  parameters = <<PARAMETERS
    {
      "deniedPorts": {
        "type": "Array",
        "metadata": {
          "displayName": "Ports to block",
          "description": "The inbound ports that should be blocked"
        }
      }
    }
  PARAMETERS

  policy_rule = <<POLICY_RULE
    {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.Network/networkSecurityGroups/securityRules"
          },
          {
            "allOf": [
              {
                "field": "Microsoft.Network/networkSecurityGroups/securityRules/access",
                "equals": "Allow"
              },
              {
                "field": "Microsoft.Network/networkSecurityGroups/securityRules/direction",
                "equals": "Inbound"
              },
              {
                "anyOf": [
                  {
                    "field": "Microsoft.Network/networkSecurityGroups/securityRules/destinationPortRange",
                    "in": "[parameters('deniedPorts')]"
                  },
                  {
                    "not": {
                      "field": "Microsoft.Network/networkSecurityGroups/securityRules/destinationPortRanges[*]",
                      "notIn": "[parameters('deniedPorts')]"
                    }
                  }
                ]
              },
              {
                "anyOf": [
                  {
                    "field": "Microsoft.Network/networkSecurityGroups/securityRules/sourceAddressPrefix",
                    "in": [
                      "*",
                      "0.0.0.0/0",
                      "Internet"
                    ]
                  },
                  {
                    "not": {
                      "field": "Microsoft.Network/networkSecurityGroups/securityRules/sourceAddressPrefixes[*]",
                      "notIn": [
                        "*",
                        "0.0.0.0/0",
                        "Internet"
                      ]
                    }
                  }
                ]
              }
            ]
          }
        ]
      },
      "then": {
        "effect": "deny"
      }
    }
  POLICY_RULE
}

resource "azurerm_policy_definition" "denied_regions" {
  name                = "Deny specific regions"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "Deny resources creation on certain specific regions"
  # management_group_id = azurerm_management_group.your_root_managementgroup.id # <- if you want to attach the policy at the management-group level

  metadata = <<METADATA
    {
      "category": "General"
    }
  METADATA

  parameters = <<PARAMETERS
    {
      "deniedRegions": {
        "type": "Array",
        "metadata": {
          "displayName": "Regions to block",
          "description": "The Regions where it is not allowed to deploy resources"
        }
      }
    }
  PARAMETERS

  policy_rule = <<POLICY_RULE
    {
      "if": {
        "field": "location",
        "in": "[parameters('deniedRegions')]"
      },
      "then": {
        "effect": "deny"
      }
    }
  POLICY_RULE
}
