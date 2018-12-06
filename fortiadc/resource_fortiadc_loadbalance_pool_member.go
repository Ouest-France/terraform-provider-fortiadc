package fortiadc

import (
	"fmt"
	"strconv"

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
				ForceNew: true,
			},
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	client := m.(*gofortiadc.Client)

	res, err := client.LoadbalanceGetPoolMember(d.Get("pool").(string), d.Id())
	if err != nil {
		return err
	}

	d.Set("name", res.RealServerID)
	d.Set("status", res.Status)

	port, err := strconv.ParseInt(res.Port, 10, 64)
	if err != nil {
		return err
	}
	d.Set("port", port)

	weight, err := strconv.ParseInt(res.Weight, 10, 64)
	if err != nil {
		return err
	}
	d.Set("weight", weight)

	connLimit, err := strconv.ParseInt(res.Connlimit, 10, 64)
	if err != nil {
		return err
	}
	d.Set("conn_limit", connLimit)

	recover, err := strconv.ParseInt(res.Recover, 10, 64)
	if err != nil {
		return err
	}
	d.Set("recover", recover)

	warmup, err := strconv.ParseInt(res.Warmup, 10, 64)
	if err != nil {
		return err
	}
	d.Set("warmup", warmup)

	warmrate, err := strconv.ParseInt(res.Warmrate, 10, 64)
	if err != nil {
		return err
	}
	d.Set("warmrate", warmrate)

	connectionRateLimit, err := strconv.ParseInt(res.ConnectionRateLimit, 10, 64)
	if err != nil {
		return err
	}
	d.Set("conn_rate_limit", connectionRateLimit)

	return nil
}

func resourceFortiadcLoadbalancePoolMemberUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalancePoolMember{
		Mkey:                     d.Id(),
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

	err := client.LoadbalanceUpdatePoolMember(d.Get("pool").(string), d.Id(), req)
	if err != nil {
		return err
	}

	return resourceFortiadcLoadbalancePoolMemberRead(d, m)
}

func resourceFortiadcLoadbalancePoolMemberDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	return client.LoadbalanceDeletePoolMember(d.Get("pool").(string), d.Id())
}
