##### Azure policies assignments
####  * Enforces secure defaults for the resources deployed on Azure
###   * https://docs.microsoft.com/en-us/azure/governance/policy/overview
##    * https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/policy_definition

##### Network controls
####  * Services wrongly exposed
###   * too permissive security groups

resource "azurerm_management_group_policy_assignment" "nsg_deny_any_inbound" {
  name                     = "nsg deny any any rule"
  display_name             = "ANY:ANY rules on Network Security Group Rules are not permitted"
  policy_definition_id     = azurerm_policy_definition.nsg_deny_inbound_ports.id
  # management_group_id = azurerm_management_group.your_root_managementgroup.id # <- if you want to attach the policy at the management-group level
  parameters = <<PARAMETERS
  {
    "deniedPorts": {
      "value": [ "*" ]
    }
  }
  PARAMETERS
  non_compliance_message {
    content = "ANY:ANY rules on Network Security Group Rules are not permitted, restrict your rule to certain ports or source IP addresses"
  }
}

resource "azurerm_management_group_policy_assignment" "nsr_deny_any_inbound" {
  name                     = "nsr deny any any rule"
  display_name             = "ANY:ANY rules on Network Security Rules are not permitted"
  policy_definition_id     = azurerm_policy_definition.nsr_deny_inbound_ports.id
  # management_group_id = azurerm_management_group.your_root_managementgroup.id # <- if you want to attach the policy at the management-group level
  parameters = <<PARAMETERS
  {
    "deniedPorts": {
      "value": [ "*" ]
    }
  }
  PARAMETERS
  non_compliance_message {
    content = "ANY:ANY rules on Network Security Rules are not permitted, restrict your rule to certain ports or source IP addresses"
  }
}

resource "azurerm_management_group_policy_assignment" "nsg_deny_ssh_inbound" {
  name                     = "nsg deny any ssh rule"
  display_name             = "ANY:SSH rules on Network Security Group Rules are not permitted"
  policy_definition_id     = azurerm_policy_definition.nsg_deny_inbound_ports.id
  # management_group_id = azurerm_management_group.your_root_managementgroup.id # <- if you want to attach the policy at the management-group level
  parameters = <<PARAMETERS
  {
    "deniedPorts": {
      "value": [ "22" ]
    }
  }
  PARAMETERS
  non_compliance_message {
    content = "Exposing SSH to the Internet without filtering is not permitted, restrict your rule to your source IP address"
  }
}

resource "azurerm_management_group_policy_assignment" "nsr_deny_ssh_inbound" {
  name                     = "nsr deny any ssh rule"
  display_name             = "ANY:SSH rules on Network Security Rules are not permitted"
  policy_definition_id     = azurerm_policy_definition.nsr_deny_inbound_ports.id
  # management_group_id = azurerm_management_group.your_root_managementgroup.id # <- if you want to attach the policy at the management-group level
  parameters = <<PARAMETERS
  {
    "deniedPorts": {
      "value": [ "22" ]
    }
  }
  PARAMETERS
  non_compliance_message {
    content = "Exposing SSH to the Internet without filtering is not permitted, restrict your rule to your source IP address"
  }
}

resource "azurerm_management_group_policy_assignment" "nsg_deny_rdp_inbound" {
  name                     = "nsg deny any rdp rule"
  display_name             = "ANY:RDP rules on Network Security Group Rules are not permitted"
  policy_definition_id     = azurerm_policy_definition.nsg_deny_inbound_ports.id
  # management_group_id = azurerm_management_group.your_root_managementgroup.id # <- if you want to attach the policy at the management-group level
  parameters = <<PARAMETERS
  {
    "deniedPorts": {
      "value": [ "3389" ]
    }
  }
  PARAMETERS
  non_compliance_message {
    content = "Exposing RDP to the Internet without filtering is not permitted, restrict your rule to your source IP address"
  }
}

resource "azurerm_management_group_policy_assignment" "nsr_deny_rdp_inbound" {
  name                     = "nsr deny any rdp rule"
  display_name             = "ANY:RDP rules on Network Security Rules are not permitted"
  policy_definition_id     = azurerm_policy_definition.nsr_deny_inbound_ports.id
  # management_group_id = azurerm_management_group.your_root_managementgroup.id # <- if you want to attach the policy at the management-group level
  parameters = <<PARAMETERS
  {
    "deniedPorts": {
      "value": [ "3389" ]
    }
  }
  PARAMETERS
  non_compliance_message {
    content = "Exposing RDP to the Internet without filtering is not permitted, restrict your rule to your source IP address"
  }
}

##### Generic
####  * Regions allowed or disallowed to deploy resources
###   * Resources tags present

resource "azurerm_management_group_policy_assignment" "denied_regions" {
  name                     = "denied regions"
  display_name             = "Deploying resources in Brazil Southeast is not permitted"
  policy_definition_id     = azurerm_policy_definition.denied_regions.id
  # management_group_id = azurerm_management_group.your_root_managementgroup.id # <- if you want to attach the policy at the management-group level
  parameters = <<PARAMETERS
  {
    "deniedRegions": {
      "value": [ "Brazil Southeast" ]
    }
  }
  PARAMETERS
  non_compliance_message {
    content = "Deploying resources in Brazil Southeast is not permitted (due to costs reason)"
  }
}