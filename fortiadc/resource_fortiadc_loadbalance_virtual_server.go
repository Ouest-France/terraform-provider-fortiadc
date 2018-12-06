package fortiadc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFortiadcLoadbalanceVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalanceVirtualServerCreate,
		Read:   resourceFortiadcLoadbalanceVirtualServerRead,
		Update: resourceFortiadcLoadbalanceVirtualServerUpdate,
		Delete: resourceFortiadcLoadbalanceVirtualServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enable",
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "l4-load-balance",
			},
			"address_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ipv4",
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"packet_forward_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NAT",
			},
			"nat_source_pool": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"connection_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"content_routing_enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"connection_rate_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"interface": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "port1",
			},
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LB_PROF_TCP",
			},
			"method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LB_METHOD_ROUND_ROBIN",
			},
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFortiadcLoadbalanceVirtualServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	contentRouting := "disable"
	if d.Get("content_routing_enable").(bool) {
		contentRouting = "enable"
	}

	if d.Get("packet_forward_method").(string) != "FullNAT" && len(d.Get("nat_source_pool").(string)) > 0 {
		return errors.New("nat_source_pool cannot be defined when packet_forward_method is not FullNAT")
	}

	req := gofortiadc.LoadbalanceVirtualServerReq{
		Status:              d.Get("status").(string),
		Type:                d.Get("type").(string),
		AddrType:            d.Get("address_type").(string),
		Address:             d.Get("address").(string),
		Address6:            "::",
		PacketFwdMethod:     d.Get("packet_forward_method").(string),
		SrcPool:             d.Get("nat_source_pool").(string),
		Port:                fmt.Sprintf("%d", d.Get("port").(int)),
		PortRange:           "0",
		ConnectionLimit:     fmt.Sprintf("%d", d.Get("connection_limit").(int)),
		ContentRouting:      contentRouting,
		ContentRewriting:    "disable",
		Warmup:              "0",
		Warmrate:            "10",
		ConnectionRateLimit: fmt.Sprintf("%d", d.Get("connection_rate_limit").(int)),
		Log:                 "enable",
		Alone:               "enable",
		Mkey:                d.Get("name").(string),
		Interface:           d.Get("interface").(string),
		Profile:             d.Get("profile").(string),
		Method:              d.Get("method").(string),
		Pool:                d.Get("pool").(string),
	}

	err := client.LoadbalanceCreateVirtualServer(req)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return resourceFortiadcLoadbalanceVirtualServerRead(d, m)
}

func resourceFortiadcLoadbalanceVirtualServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	rs, err := client.LoadbalanceGetVirtualServer(d.Id())
	if err != nil {
		return err
	}

	contentRouting := false
	if rs.ContentRouting == "enable" {
		contentRouting = true
	}

	d.Set("status", rs.Status)
	d.Set("type", rs.Type)
	d.Set("address_type", rs.AddrType)
	d.Set("address", rs.Address)
	d.Set("packet_forward_method", rs.PacketFwdMethod)
	d.Set("nat_source_pool", strings.TrimSpace(rs.SrcPool))
	d.Set("content_routing_enable", contentRouting)
	d.Set("interface", rs.Interface)
	d.Set("profile", rs.Profile)
	d.Set("method", rs.Method)
	d.Set("pool", rs.Pool)

	port, err := strconv.ParseInt(strings.TrimSpace(rs.Port), 10, 64)
	if err != nil {
		return err
	}
	d.Set("port", port)

	connectionLimit, err := strconv.ParseInt(rs.ConnectionLimit, 10, 64)
	if err != nil {
		return err
	}
	d.Set("connection_limit", connectionLimit)

	connectionRateLimit, err := strconv.ParseInt(rs.ConnectionRateLimit, 10, 64)
	if err != nil {
		return err
	}
	d.Set("connection_rate_limit", connectionRateLimit)

	return nil
}

func resourceFortiadcLoadbalanceVirtualServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	contentRouting := "disable"
	if d.Get("content_routing_enable").(bool) {
		contentRouting = "enable"
	}

	if d.Get("packet_forward_method").(string) != "FullNAT" && len(d.Get("nat_source_pool").(string)) > 0 {
		return errors.New("nat_source_pool cannot be defined when packet_forward_method is not FullNAT")
	}

	req := gofortiadc.LoadbalanceVirtualServerReq{
		Status:              d.Get("status").(string),
		Type:                d.Get("type").(string),
		AddrType:            d.Get("address_type").(string),
		Address:             d.Get("address").(string),
		Address6:            "::",
		PacketFwdMethod:     d.Get("packet_forward_method").(string),
		SrcPool:             d.Get("nat_source_pool").(string),
		Port:                fmt.Sprintf("%d", d.Get("port").(int)),
		PortRange:           "0",
		ConnectionLimit:     fmt.Sprintf("%d", d.Get("connection_limit").(int)),
		ContentRouting:      contentRouting,
		ContentRewriting:    "disable",
		Warmup:              "0",
		Warmrate:            "10",
		ConnectionRateLimit: fmt.Sprintf("%d", d.Get("connection_rate_limit").(int)),
		Log:                 "enable",
		Alone:               "enable",
		Mkey:                d.Get("name").(string),
		Interface:           d.Get("interface").(string),
		Profile:             d.Get("profile").(string),
		Method:              d.Get("method").(string),
		Pool:                d.Get("pool").(string),
	}

	err := client.LoadbalanceUpdateVirtualServer(req)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return resourceFortiadcLoadbalanceVirtualServerRead(d, m)
}

func resourceFortiadcLoadbalanceVirtualServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	return client.LoadbalanceDeleteVirtualServer(d.Get("name").(string))
}
