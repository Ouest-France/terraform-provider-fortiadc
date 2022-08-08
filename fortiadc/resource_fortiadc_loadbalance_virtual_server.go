package fortiadc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFortiadcLoadbalanceVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalanceVirtualServerCreate,
		Read:   resourceFortiadcLoadbalanceVirtualServerRead,
		Update: resourceFortiadcLoadbalanceVirtualServerUpdate,
		Delete: resourceFortiadcLoadbalanceVirtualServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enable",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "l4-load-balance",
			},
			"address_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ipv4",
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"packet_forward_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NAT",
			},
			"source_pool_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"connection_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"content_routing_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"content_routing_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"content_rewriting_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"content_rewriting_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connection_rate_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"error_page": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"error_msg": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Server-unavailable!",
			},
			"interface": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "port1",
			},
			"profile": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LB_PROF_TCP",
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LB_METHOD_ROUND_ROBIN",
			},
			"pool": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_ssl_profile": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"http_to_https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"persistence": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"traffic_log": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"transaction_rate_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
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

	spList := []string{}
	if raw, ok := d.GetOk("source_pool_list"); ok {
		for _, v := range raw.([]interface{}) {
			spList = append(spList, v.(string))
		}
	}

	// Content routing
	if d.Get("content_routing_enable").(bool) && len(crList) == 0 {
		return errors.New("content_routing_list cannot be empty when content_routing_enable is set to true")
	}

	if !d.Get("content_routing_enable").(bool) && len(crList) > 0 {
		return errors.New("content_routing_list must be empty when content_routing_enable is set to false")
	}

	// Content rewriting
	if d.Get("content_rewriting_enable").(bool) && len(rwList) == 0 {
		return errors.New("content_routing_list cannot be empty when content_rewriting_enable is set to true")
	}

	if !d.Get("content_rewriting_enable").(bool) && len(rwList) > 0 {
		return errors.New("content_rewriting_list must be empty when content_rewriting_enable is set to false")
	}

	http2https := "disable"
	if d.Get("http_to_https").(bool) {
		http2https = "enable"
	}
	if len(d.Get("client_ssl_profile").(string)) == 0 {
		http2https = ""
	}

	req := gofortiadc.LoadbalanceVirtualServer{
		Status:               d.Get("status").(string),
		Type:                 d.Get("type").(string),
		AddrType:             d.Get("address_type").(string),
		Address:              d.Get("address").(string),
		Address6:             "::",
		PacketFwdMethod:      d.Get("packet_forward_method").(string),
		SourcePoolList:       strings.Join(spList, " "),
		Port:                 fmt.Sprintf("%d", d.Get("port").(int)),
		ConnectionLimit:      fmt.Sprintf("%d", d.Get("connection_limit").(int)),
		ContentRouting:       boolToEnable(d.Get("content_routing_enable").(bool)),
		ContentRoutingList:   strings.Join(crList, " "),
		ContentRewriting:     boolToEnable(d.Get("content_rewriting_enable").(bool)),
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
		Persistence:          d.Get("persistence").(string),
		ErrorMsg:             d.Get("error_msg").(string),
		ErrorPage:            d.Get("error_page").(string),
		TrafficLog:           boolToEnable(d.Get("traffic_log").(bool)),
		TransRateLimit:       fmt.Sprintf("%d", d.Get("transaction_rate_limit").(int)),
		Comments:             d.Get("comments").(string),
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
	if errors.Is(err, gofortiadc.ErrNotFound) {
		// If not found, remove from state
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	contentRouting := enableToBool(rs.ContentRouting)
	crList := strings.Split(rs.ContentRoutingList, " ")
	if len(crList) > 1 {
		crList = crList[:len(crList)-1]
	}
	if !contentRouting {
		crList = []string{}
	}

	contentRewriting := enableToBool(rs.ContentRewriting)
	rwList := strings.Split(rs.ContentRewritingList, " ")
	if len(rwList) > 1 {
		rwList = rwList[:len(rwList)-1]
	}
	if !contentRewriting {
		rwList = []string{}
	}

	spList := []string{}
	if len(rs.SourcePoolList) != 0 {
		spList = strings.Split(rs.SourcePoolList, " ")
	}
	if len(spList) > 1 {
		spList = spList[:len(spList)-1]
	}

	arguments := map[string]interface{}{
		"status":                   rs.Status,
		"type":                     rs.Type,
		"address_type":             rs.AddrType,
		"address":                  rs.Address,
		"packet_forward_method":    rs.PacketFwdMethod,
		"source_pool_list":         spList,
		"content_routing_enable":   contentRouting,
		"content_routing_list":     crList,
		"content_rewriting_enable": contentRewriting,
		"content_rewriting_list":   rwList,
		"interface":                rs.Interface,
		"profile":                  rs.Profile,
		"method":                   rs.Method,
		"pool":                     rs.Pool,
		"client_ssl_profile":       rs.ClientSSLProfile,
		"http_to_https":            enableToBool(rs.HTTP2HTTPS),
		"persistence":              rs.Persistence,
		"error_msg":                rs.ErrorMsg,
		"error_page":               rs.ErrorPage,
		"traffic_log":              enableToBool(rs.TrafficLog),
		"comments":                 rs.Comments,
	}

	port, err := strconv.ParseInt(strings.TrimSpace(rs.Port), 10, 64)
	if err != nil {
		return err
	}
	arguments["port"] = port

	connectionLimit, err := strconv.ParseInt(rs.ConnectionLimit, 10, 64)
	if err != nil {
		return err
	}
	arguments["connection_limit"] = connectionLimit

	transactionRateLimit, err := strconv.ParseInt(rs.TransRateLimit, 10, 64)
	if err != nil {
		return err
	}
	arguments["transaction_rate_limit"] = transactionRateLimit

	connectionRateLimit, err := strconv.ParseInt(rs.ConnectionRateLimit, 10, 64)
	if err != nil {
		return err
	}
	arguments["connection_rate_limit"] = connectionRateLimit

	for arg, value := range arguments {
		err = d.Set(arg, value)
		if err != nil {
			return err
		}
	}

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

	spList := []string{}
	if raw, ok := d.GetOk("source_pool_list"); ok {
		for _, v := range raw.([]interface{}) {
			spList = append(spList, v.(string))
		}
	}

	// Content routing
	if d.Get("content_routing_enable").(bool) && len(crList) == 0 {
		return errors.New("content_routing_list cannot be empty when content_routing_enable is set to true")
	}

	if !d.Get("content_routing_enable").(bool) && len(crList) > 0 {
		return errors.New("content_routing_list must be empty when content_routing_enable is set to false")
	}

	// Content rewriting
	if d.Get("content_rewriting_enable").(bool) && len(rwList) == 0 {
		return errors.New("content_rewriting_list cannot be empty when content_rewriting_enable is set to true")
	}

	if !d.Get("content_rewriting_enable").(bool) && len(rwList) > 0 {
		return errors.New("content_rewriting_list must be empty when content_rewriting_enable is set to false")
	}

	http2https := boolToEnable(d.Get("http_to_https").(bool))
	if len(d.Get("client_ssl_profile").(string)) == 0 {
		http2https = ""
	}

	req := gofortiadc.LoadbalanceVirtualServer{
		Status:               d.Get("status").(string),
		Type:                 d.Get("type").(string),
		AddrType:             d.Get("address_type").(string),
		Address:              d.Get("address").(string),
		Address6:             "::",
		PacketFwdMethod:      d.Get("packet_forward_method").(string),
		SourcePoolList:       strings.Join(spList, " "),
		Port:                 fmt.Sprintf("%d", d.Get("port").(int)),
		ConnectionLimit:      fmt.Sprintf("%d", d.Get("connection_limit").(int)),
		ContentRouting:       boolToEnable(d.Get("content_routing_enable").(bool)),
		ContentRoutingList:   strings.Join(crList, " "),
		ContentRewriting:     boolToEnable(d.Get("content_rewriting_enable").(bool)),
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
		Persistence:          d.Get("persistence").(string),
		ErrorMsg:             d.Get("error_msg").(string),
		ErrorPage:            d.Get("error_page").(string),
		TrafficLog:           boolToEnable(d.Get("traffic_log").(bool)),
		TransRateLimit:       fmt.Sprintf("%d", d.Get("transaction_rate_limit").(int)),
		Comments:             d.Get("comments").(string),
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
