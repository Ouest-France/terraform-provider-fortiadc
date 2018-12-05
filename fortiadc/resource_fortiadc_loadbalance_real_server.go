package fortiadc

import (
	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFortiadcLoadbalanceRealServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalanceRealServerCreate,
		Read:   resourceFortiadcLoadbalanceRealServerRead,
		Update: resourceFortiadcLoadbalanceRealServerUpdate,
		Delete: resourceFortiadcLoadbalanceRealServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"address6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "::",
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enable",
			},
		},
	}
}

func resourceFortiadcLoadbalanceRealServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalanceRealServer{
		Mkey:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		Address:  d.Get("address").(string),
		Address6: d.Get("address6").(string),
	}

	err := client.LoadbalanceCreateRealServer(req)
	if err != nil {
		return err
	}

	d.SetId(req.Mkey)

	return resourceFortiadcLoadbalanceRealServerRead(d, m)
}

func resourceFortiadcLoadbalanceRealServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	rs, err := client.LoadbalanceGetRealServer(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", rs.Mkey)
	d.Set("address", rs.Address)
	d.Set("address6", rs.Address6)
	d.Set("status", rs.Status)

	return nil
}

func resourceFortiadcLoadbalanceRealServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalanceRealServer{
		Mkey:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		Address:  d.Get("address").(string),
		Address6: d.Get("address6").(string),
	}

	err := client.LoadbalanceUpdateRealServer(req)
	if err != nil {
		return err
	}

	return resourceFortiadcLoadbalanceRealServerRead(d, m)
}

func resourceFortiadcLoadbalanceRealServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	return client.LoadbalanceDeleteRealServer(d.Get("name").(string))
}
