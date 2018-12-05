package fortiadc

import (
	"errors"
	"strings"

	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFortiadcLoadbalancePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalancePoolCreate,
		Read:   resourceFortiadcLoadbalancePoolRead,
		Update: resourceFortiadcLoadbalancePoolUpdate,
		Delete: resourceFortiadcLoadbalancePoolDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pool_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ipv4",
			},
			"healtcheck_enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"healtcheck_relationship": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AND",
			},
			"healtcheck_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"real_server_ssl_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NONE",
			},
		},
	}
}

func resourceFortiadcLoadbalancePoolCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	hcList := []string{}
	if raw, ok := d.GetOk("healtcheck_list"); ok {
		for _, v := range raw.([]interface{}) {
			hcList = append(hcList, v.(string))
		}
	}

	if d.Get("healtcheck_enable").(bool) && len(hcList) == 0 {
		return errors.New("healtcheck_list cannot be empty when healtcheck_enable is set to true")
	}

	if !d.Get("healtcheck_enable").(bool) && len(hcList) > 0 {
		return errors.New("healtcheck_list must be empty when healtcheck_enable is set to false")
	}

	healtcheckEnable := "disable"
	if d.Get("healtcheck_enable").(bool) {
		healtcheckEnable = "enable"
	}

	req := gofortiadc.LoadbalancePoolReq{
		Mkey:                    d.Get("name").(string),
		PoolType:                d.Get("pool_type").(string),
		HealthCheck:             healtcheckEnable,
		HealthCheckRelationship: d.Get("healtcheck_relationship").(string),
		HealthCheckList:         strings.Join(hcList, " "),
		RsProfile:               d.Get("real_server_ssl_profile").(string),
	}

	err := client.LoadbalanceCreatePool(req)
	if err != nil {
		return err
	}

	d.SetId(req.Mkey)

	return resourceFortiadcLoadbalancePoolRead(d, m)
}

func resourceFortiadcLoadbalancePoolRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	res, err := client.LoadbalanceGetPool(d.Id())
	if err != nil {
		return err
	}

	healtcheckEnable := false
	if res.HealthCheck == "enable" {
		healtcheckEnable = true
	}

	hcList := strings.Split(res.HealthCheckList, " ")
	if len(hcList) > 1 {
		hcList = hcList[:len(hcList)-1]
	}
	if healtcheckEnable == false {
		hcList = []string{}
	}

	d.Set("name", res.Mkey)
	d.Set("pool_type", res.PoolType)
	d.Set("healtcheck_enable", healtcheckEnable)
	d.Set("healtcheck_relationship", res.HealthCheckRelationship)
	d.Set("healtcheck_list", hcList)
	d.Set("real_server_ssl_profile", res.RsProfile)

	return nil
}

func resourceFortiadcLoadbalancePoolUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	hcList := []string{}
	if raw, ok := d.GetOk("healtcheck_list"); ok {
		for _, v := range raw.([]interface{}) {
			hcList = append(hcList, v.(string))
		}
	}

	if d.Get("healtcheck_enable").(bool) && len(hcList) == 0 {
		return errors.New("healtcheck_list cannot be empty when healtcheck_enable is set to true")
	}

	if !d.Get("healtcheck_enable").(bool) && len(hcList) > 0 {
		return errors.New("healtcheck_list must be empty when healtcheck_enable is set to false")
	}

	healtcheckEnable := "disable"
	if d.Get("healtcheck_enable").(bool) {
		healtcheckEnable = "enable"
	}

	req := gofortiadc.LoadbalancePoolReq{
		Mkey:                    d.Get("name").(string),
		PoolType:                d.Get("pool_type").(string),
		HealthCheck:             healtcheckEnable,
		HealthCheckRelationship: d.Get("healtcheck_relationship").(string),
		HealthCheckList:         strings.Join(hcList, " "),
		RsProfile:               d.Get("real_server_ssl_profile").(string),
	}

	err := client.LoadbalanceUpdatePool(req)
	if err != nil {
		return err
	}

	return resourceFortiadcLoadbalancePoolRead(d, m)
}

func resourceFortiadcLoadbalancePoolDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	return client.LoadbalanceDeletePool(d.Get("name").(string))
}
