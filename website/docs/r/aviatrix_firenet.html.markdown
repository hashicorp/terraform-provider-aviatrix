---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firenet"
description: |-
  Creates and manages Aviatrix FireNets
---

# aviatrix_firewall_instance

The aviatrix_firenet resource allows the creation and management of Aviatrix FireNets.

## Example Usage

```hcl
# Create an Aviatrix FireNet associated to a Firewall Instance
resource "aviatrix_firewall" "test_firewall" {
  vpc_id             = "vpc-032005cc371"
  inspection_enabled = true
  egress_enabled     = false
  
  firewall_instance_association {
    gw_name              = "avx_firenet_gw"
    instance_id          = "i-09dc118db6a1eb901"
    firewall_name        = "avx_firewall_instance"
    attached             = true
    lan_interface        = "eni-0a34b1827bf222353"
    management_interface = "eni-030e53176c7f7d34a"
    egress_interface     = "eni-03b8dd53a1a731481"
  }
}

# Create an Aviatrix FireNet associated to an FQDN Gateway
resource "aviatrix_firewall" "test_firewall" {
  vpc_id             = "vpc-032005cc371"
  inspection_enabled = true
  egress_enabled     = false
  
  firewall_instance_association {
    gw_name       = "avx_firenet_gw"
    firewall_name = "avx_fqdn_gateway"
    vendor_type   = "fqdn_gateway"
    attached      = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) ID of the Security VPC.
* `inspection_enabled` - (Optional) Enable/Disable traffic inspection. Valid values: true, false. Default value: true.
* `egress_enabled` - (Optional) Enable/Disable egress through firewall. Valid values: true, false. Default value: false.
* `firewall_instance_association` - (Optional) List of firewall instances to be associated with fireNet.
  * `gw_name` - (Required) Name of the primary FireNet gateway.
  * `vendor_type` - (Required) Type of the firewall. Valid values: "firewall_instance", "fqdn_gateway". Default value: "firewall_instance".  
  * `firewall_name` - (Required) Firewall instance name. If associating FQDN gateway to fireNet, it is FQDN gateway's gw_name.
  * `instance_id` - (Required) ID of Firewall instance, required if it is a firewall instance.
  * `lan_interface`- (Optional) Lan interface ID, required if it is a firewall instance.
  * `management_interface` - (Optional) Management interface ID, required if it is a firewall instance.
  * `egress_interface`- (Optional) Egress interface ID, required if it is a firewall instance.
  * `attached`- (Optional) Switch to attach/detach firewall instance to/from fireNet. Valid values: true, false. Default value: false.
                                                                      
## Import

Instance firenet can be imported using the vpc_id, e.g.

```
$ terraform import aviatrix_firenet.test vpc_id
```