package fortiadc

import (
	"fmt"

	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFortiadcLoadbalancePoolMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalancePoolMemberCreate,
		Read:   resourceFortiadcLoadbalancePoolMemberRead,
		Update: resourceFortiadcLoadbalancePoolMemberUpdate,
		Delete: resourceFortiadcLoadbalancePoolMemberDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enable",
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"weight": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"conn_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"recover": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"warmup": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"warmrate": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"conn_rate_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceFortiadcLoadbalancePoolMemberCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalancePoolMember{
		Address:                  "0.0.0.0",
		Address6:                 "::",
		Backup:                   "disable",
		ConnectionRateLimit:      fmt.Sprintf("%d", d.Get("conn_rate_limit").(int)),
		Connlimit:                fmt.Sprintf("%d", d.Get("conn_limit").(int)),
		Cookie:                   "",
		HcStatus:                 "1",
		HealthCheckInherit:       "enable",
		MHealthCheck:             "disable",
		MHealthCheckRelationship: "AND",
		MHealthCheckList:         "",
		MysqlGroupID:             "0",
		MysqlReadOnly:            "disable",
		Port:                     fmt.Sprintf("%d", d.Get("port").(int)),
		RealServerID:             d.Get("name").(string),
		Recover:                  fmt.Sprintf("%d", d.Get("recover").(int)),
		RsProfileInherit:         "enable",
		Ssl:                      "disable",
		Status:                   d.Get("status").(string),
		Weight:                   fmt.Sprintf("%d", d.Get("weight").(int)),
		Warmup:                   fmt.Sprintf("%d", d.Get("warmup").(int)),
		Warmrate:                 fmt.Sprintf("%d", d.Get("warmrate").(int)),
	}

	err := client.LoadbalanceCreatePoolMember(d.Get("pool").(string), req)
	if err != nil {
		return err
	}

	id, err := client.LoadbalanceGetPoolMemberID(d.Get("pool").(string), d.Get("name").(string))
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceFortiadcLoadbalancePoolMemberRead(d, m)
}

func resourceFortiadcLoadbalancePoolMemberRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFortiadcLoadbalancePoolMemberUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceFortiadcLoadbalancePoolMemberRead(d, m)
}

func resourceFortiadcLoadbalancePoolMemberDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	return client.LoadbalanceDeletePoolMember(d.Get("pool").(string), d.Id())
}
