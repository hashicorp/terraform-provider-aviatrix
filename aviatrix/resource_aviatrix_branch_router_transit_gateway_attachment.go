package aviatrix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixBranchRouterTransitGatewayAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixBranchRouterTransitGatewayAttachmentCreate,
		Read:   resourceAviatrixBranchRouterTransitGatewayAttachmentRead,
		Delete: resourceAviatrixBranchRouterTransitGatewayAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"branch_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Branch router name.",
			},
			"transit_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Aviatrix Transit Gateway name.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Connection name.",
			},
			"transit_gateway_bgp_asn": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "BGP AS Number for transit gateway.",
			},
			"branch_router_bgp_asn": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "BGP AS Number for branch router.",
			},
			"phase1_authentication": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "SHA-256",
				Description: "Phase 1 authentication algorithm.",
			},
			"phase1_dh_groups": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     14,
				Description: "Phase 1 Diffie-Hellman groups.",
			},
			"phase1_encryption": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "AES-256-CBC",
				Description: "Phase 1 encryption algorithm.",
			},
			"phase2_authentication": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "HMAC-SHA-256",
				Description: "Phase 2 authentication algorithm.",
			},
			"phase2_dh_groups": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     14,
				Description: "Phase 2 Diffie-Hellman groups.",
			},
			"phase2_encryption": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "AES-256-CBC",
				Description: "Phase 2 encryption algorithm.",
			},
			"enable_global_accelerator": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable AWS Global Accelerator",
			},
			"pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Pre-shared Key.",
			},
			"local_tunnel_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Local tunnel IP",
			},
			"remote_tunnel_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote tunnel IP",
			},
		},
	}
}

func marshalBranchRouterTransitGatewayAttachmentInput(d *schema.ResourceData) *goaviatrix.BranchRouterTransitGatewayAttachment {
	brata := &goaviatrix.BranchRouterTransitGatewayAttachment{
		BranchName:              d.Get("branch_name").(string),
		TransitGatewayName:      d.Get("transit_gateway_name").(string),
		ConnectionName:          d.Get("connection_name").(string),
		RoutingProtocol:         "bgp",
		TransitGatewayBgpAsn:    strconv.Itoa(d.Get("transit_gateway_bgp_asn").(int)),
		BranchRouterBgpAsn:      strconv.Itoa(d.Get("branch_router_bgp_asn").(int)),
		Phase1Authentication:    d.Get("phase1_authentication").(string),
		Phase1DHGroups:          strconv.Itoa(d.Get("phase1_dh_groups").(int)),
		Phase1Encryption:        d.Get("phase1_encryption").(string),
		Phase2Authentication:    d.Get("phase2_authentication").(string),
		Phase2DHGroups:          strconv.Itoa(d.Get("phase2_dh_groups").(int)),
		Phase2Encryption:        d.Get("phase2_encryption").(string),
		EnableGlobalAccelerator: strconv.FormatBool(d.Get("enable_global_accelerator").(bool)),
		PreSharedKey:            d.Get("pre_shared_key").(string),
		LocalTunnelIP:           d.Get("local_tunnel_ip").(string),
		RemoteTunnelIP:          d.Get("remote_tunnel_ip").(string),
	}

	return brata
}

func resourceAviatrixBranchRouterTransitGatewayAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	brata := marshalBranchRouterTransitGatewayAttachmentInput(d)

	if err := client.CreateBranchRouterTransitGatewayAttachment(brata); err != nil {
		return err
	}

	d.SetId(brata.ConnectionName)
	return nil
}

func resourceAviatrixBranchRouterTransitGatewayAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	isImport := false
	if connectionName == "" {
		isImport = true
		id := d.Id()
		d.SetId(id)
		connectionName = id
		log.Printf("[DEBUG] Looks like an import, no branch_router_transit_gateway_attachment connection_name received. Import Id is %s", id)
	}

	brata := &goaviatrix.BranchRouterTransitGatewayAttachment{
		ConnectionName: connectionName,
	}

	brata, err := client.GetBranchRouterTransitGatewayAttachment(brata)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find branch_router_transit_gateway_attachment %s: %v", connectionName, err)
	}

	d.Set("branch_name", brata.BranchName)
	d.Set("transit_gateway_name", brata.TransitGatewayName)
	d.Set("connection_name", brata.ConnectionName)

	transitGatewayBgpAsn, err := strconv.Atoi(brata.TransitGatewayBgpAsn)
	if err != nil {
		return fmt.Errorf("could not convert transitGatewayBgpAsn to int: %v", err)
	}
	d.Set("transit_gateway_bgp_asn", transitGatewayBgpAsn)

	branchRouterBgpAsn, err := strconv.Atoi(brata.BranchRouterBgpAsn)
	if err != nil {
		return fmt.Errorf("could not convert branchRouterBgpAsn to int: %v", err)
	}
	d.Set("branch_router_bgp_asn", branchRouterBgpAsn)

	d.Set("phase1_authentication", brata.Phase1Authentication)

	phase1DhGroups, err := strconv.Atoi(brata.Phase1DHGroups)
	if err != nil {
		return fmt.Errorf("could not convert phase1DhGroups to int: %v", err)
	}
	d.Set("phase1_dh_groups", phase1DhGroups)

	d.Set("phase1_encryption", brata.Phase1Encryption)
	d.Set("phase2_authentication", brata.Phase2Authentication)

	phase2DhGroups, err := strconv.Atoi(brata.Phase2DHGroups)
	if err != nil {
		return fmt.Errorf("could not convert phase2DhGroups to int: %v", err)
	}
	d.Set("phase2_dh_groups", phase2DhGroups)

	d.Set("phase2_encryption", brata.Phase2Encryption)

	enableGlobalAccelerator, err := strconv.ParseBool(brata.EnableGlobalAccelerator)
	if err != nil {
		return fmt.Errorf("could not convert enableGlobalAccelerator to bool: %v", err)
	}
	d.Set("enable_global_accelerator", enableGlobalAccelerator)

	if isImport || d.Get("local_tunnel_ip") != "" {
		d.Set("local_tunnel_ip", brata.LocalTunnelIP)
	}
	if isImport || d.Get("remote_tunnel_ip") != "" {
		d.Set("remote_tunnel_ip", brata.RemoteTunnelIP)
	}

	d.SetId(brata.ConnectionName)
	return nil
}

func resourceAviatrixBranchRouterTransitGatewayAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	cn := d.Get("connection_name").(string)

	if err := client.DeleteBranchRouterAttachment(cn); err != nil {
		return err
	}

	return nil
}
