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
			"content_routing_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"content_rewriting_enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"content_rewriting_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"client_ssl_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"http_to_https": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceFortiadcLoadbalanceVirtualServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	crList := []string{}
	if raw, ok := d.GetOk("content_routing_list"); ok {
		for _, v := range raw.([]interface{}) {
			crList = append(crList, v.(string))
		}
	}

	rwList := []string{}
	if raw, ok := d.GetOk("content_rewriting_list"); ok {
		for _, v := range raw.([]interface{}) {
			rwList = append(rwList, v.(string))
		}
	}

	// Content routing
	if d.Get("content_routing_enable").(bool) && len(crList) == 0 {
		return errors.New("content_routing_list cannot be empty when content_routing_enable is set to true")
	}

	if !d.Get("content_routing_enable").(bool) && len(crList) > 0 {
		return errors.New("content_routing_list must be empty when content_routing_enable is set to false")
	}

	contentRouting := "disable"
	if d.Get("content_routing_enable").(bool) {
		contentRouting = "enable"
	}

	// Content rewriting
	if d.Get("content_rewriting_enable").(bool) && len(rwList) == 0 {
		return errors.New("content_routing_list cannot be empty when content_rewriting_enable is set to true")
	}

	if !d.Get("content_rewriting_enable").(bool) && len(rwList) > 0 {
		return errors.New("content_rewriting_list must be empty when content_rewriting_enable is set to false")
	}

	contentRewriting := "disable"
	if d.Get("content_rewriting_enable").(bool) {
		contentRewriting = "enable"
	}

	// Packet forward
	if d.Get("packet_forward_method").(string) != "FullNAT" && len(d.Get("nat_source_pool").(string)) > 0 {
		return errors.New("nat_source_pool cannot be defined when packet_forward_method is not FullNAT")
	}

	http2https := "disable"
	if d.Get("http_to_https").(bool) {
		http2https = "enable"
	}
	if len(d.Get("client_ssl_profile").(string)) == 0 {
		http2https = ""
	}

	req := gofortiadc.LoadbalanceVirtualServerReq{
		Status:               d.Get("status").(string),
		Type:                 d.Get("type").(string),
		AddrType:             d.Get("address_type").(string),
		Address:              d.Get("address").(string),
		Address6:             "::",
		PacketFwdMethod:      d.Get("packet_forward_method").(string),
		SrcPool:              d.Get("nat_source_pool").(string),
		Port:                 fmt.Sprintf("%d", d.Get("port").(int)),
		ConnectionLimit:      fmt.Sprintf("%d", d.Get("connection_limit").(int)),
		ContentRouting:       contentRouting,
		ContentRoutingList:   strings.Join(crList, " "),
		ContentRewriting:     contentRewriting,
		ContentRewritingList: strings.Join(rwList, " "),
		Warmup:               "0",
		Warmrate:             "10",
		ConnectionRateLimit:  fmt.Sprintf("%d", d.Get("connection_rate_limit").(int)),
		Alone:                "enable",
		Mkey:                 d.Get("name").(string),
		Interface:            d.Get("interface").(string),
		Profile:              d.Get("profile").(string),
		Method:               d.Get("method").(string),
		Pool:                 d.Get("pool").(string),
		ClientSSLProfile:     d.Get("client_ssl_profile").(string),
		HTTP2HTTPS:           http2https,
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

	contentRewriting := false
	if rs.ContentRewriting == "enable" {
		contentRewriting = true
	}

	http2https := false
	if rs.HTTP2HTTPS == "enable" {
		http2https = true
	}

	crList := strings.Split(rs.ContentRoutingList, " ")
	if len(crList) > 1 {
		crList = crList[:len(crList)-1]
	}
	if contentRouting == false {
		crList = []string{}
	}

	rwList := strings.Split(rs.ContentRewritingList, " ")
	if len(rwList) > 1 {
		rwList = rwList[:len(rwList)-1]
	}
	if contentRewriting == false {
		rwList = []string{}
	}

	d.Set("status", rs.Status)
	d.Set("type", rs.Type)
	d.Set("address_type", rs.AddrType)
	d.Set("address", rs.Address)
	d.Set("packet_forward_method", rs.PacketFwdMethod)
	d.Set("nat_source_pool", strings.TrimSpace(rs.SrcPool))
	d.Set("content_routing_enable", contentRouting)
	d.Set("content_routing_list", crList)
	d.Set("content_rewriting_enable", contentRewriting)
	d.Set("content_rewriting_list", rwList)
	d.Set("interface", rs.Interface)
	d.Set("profile", rs.Profile)
	d.Set("method", rs.Method)
	d.Set("pool", rs.Pool)
	d.Set("client_ssl_profile", rs.ClientSSLProfile)
	d.Set("http_to_https", http2https)

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

	crList := []string{}
	if raw, ok := d.GetOk("content_routing_list"); ok {
		for _, v := range raw.([]interface{}) {
			crList = append(crList, v.(string))
		}
	}

	rwList := []string{}
	if raw, ok := d.GetOk("content_rewriting_list"); ok {
		for _, v := range raw.([]interface{}) {
			rwList = append(rwList, v.(string))
		}
	}

	// Content routing
	if d.Get("content_routing_enable").(bool) && len(crList) == 0 {
		return errors.New("content_routing_list cannot be empty when content_routing_enable is set to true")
	}

	if !d.Get("content_routing_enable").(bool) && len(crList) > 0 {
		return errors.New("content_routing_list must be empty when content_routing_enable is set to false")
	}

	contentRouting := "disable"
	if d.Get("content_routing_enable").(bool) {
		contentRouting = "enable"
	}

	// Content rewriting
	if d.Get("content_rewriting_enable").(bool) && len(rwList) == 0 {
		return errors.New("content_rewriting_list cannot be empty when content_rewriting_enable is set to true")
	}

	if !d.Get("content_rewriting_enable").(bool) && len(rwList) > 0 {
		return errors.New("content_rewriting_list must be empty when content_rewriting_enable is set to false")
	}

	contentRewriting := "disable"
	if d.Get("content_rewriting_enable").(bool) {
		contentRewriting = "enable"
	}

	// Packet forward
	if d.Get("packet_forward_method").(string) != "FullNAT" && len(d.Get("nat_source_pool").(string)) > 0 {
		return errors.New("nat_source_pool cannot be defined when packet_forward_method is not FullNAT")
	}

	http2https := "disable"
	if d.Get("http_to_https").(bool) {
		http2https = "enable"
	}
	if len(d.Get("client_ssl_profile").(string)) == 0 {
		http2https = ""
	}

	req := gofortiadc.LoadbalanceVirtualServerReq{
		Status:               d.Get("status").(string),
		Type:                 d.Get("type").(string),
		AddrType:             d.Get("address_type").(string),
		Address:              d.Get("address").(string),
		Address6:             "::",
		PacketFwdMethod:      d.Get("packet_forward_method").(string),
		SrcPool:              d.Get("nat_source_pool").(string),
		Port:                 fmt.Sprintf("%d", d.Get("port").(int)),
		ConnectionLimit:      fmt.Sprintf("%d", d.Get("connection_limit").(int)),
		ContentRouting:       contentRouting,
		ContentRoutingList:   strings.Join(crList, " "),
		ContentRewriting:     contentRewriting,
		ContentRewritingList: strings.Join(rwList, " "),
		Warmup:               "0",
		Warmrate:             "10",
		ConnectionRateLimit:  fmt.Sprintf("%d", d.Get("connection_rate_limit").(int)),
		Alone:                "enable",
		Mkey:                 d.Get("name").(string),
		Interface:            d.Get("interface").(string),
		Profile:              d.Get("profile").(string),
		Method:               d.Get("method").(string),
		Pool:                 d.Get("pool").(string),
		ClientSSLProfile:     d.Get("client_ssl_profile").(string),
		HTTP2HTTPS:           http2https,
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
