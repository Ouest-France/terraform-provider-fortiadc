package fortiadc

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"

	"github.com/Ouest-France/gofortiadc"

	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "FortiADC address",
			},
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "FortiADC username",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "FortiADC password",
			},
			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable TLS Verify",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"fortiadc_loadbalance_real_server":                 resourceFortiadcLoadbalanceRealServer(),
			"fortiadc_loadbalance_pool":                        resourceFortiadcLoadbalancePool(),
			"fortiadc_loadbalance_pool_member":                 resourceFortiadcLoadbalancePoolMember(),
			"fortiadc_loadbalance_virtual_server":              resourceFortiadcLoadbalanceVirtualServer(),
			"fortiadc_loadbalance_content_routing":             resourceFortiadcLoadbalanceContentRouting(),
			"fortiadc_loadbalance_content_routing_condition":   resourceFortiadcLoadbalanceContentRoutingCondition(),
			"fortiadc_loadbalance_content_rewriting":           resourceFortiadcLoadbalanceContentRewriting(),
			"fortiadc_loadbalance_content_rewriting_condition": resourceFortiadcLoadbalanceContentRewritingCondition(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	cookieJar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: d.Get("insecure").(bool)},
		},
	}

	client := &gofortiadc.Client{
		Client:   httpClient,
		Address:  d.Get("address").(string),
		Username: d.Get("user").(string),
		Password: d.Get("password").(string),
	}

	err := client.Login()
	if err != nil {
		return nil, err
	}

	return client, nil
}
