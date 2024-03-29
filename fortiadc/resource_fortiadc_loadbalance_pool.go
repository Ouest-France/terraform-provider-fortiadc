package fortiadc

import (
	"errors"
	"strings"

	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFortiadcLoadbalancePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalancePoolCreate,
		Read:   resourceFortiadcLoadbalancePoolRead,
		Update: resourceFortiadcLoadbalancePoolUpdate,
		Delete: resourceFortiadcLoadbalancePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pool_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ipv4",
			},
			"healtcheck_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"healtcheck_relationship": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AND",
			},
			"healtcheck_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"real_server_ssl_profile": {
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

	req := gofortiadc.LoadbalancePool{
		Mkey:                    d.Get("name").(string),
		PoolType:                d.Get("pool_type").(string),
		HealthCheck:             boolToEnable(d.Get("healtcheck_enable").(bool)),
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
	if errors.Is(err, gofortiadc.ErrNotFound) {
		// If not found, remove from state
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	healtcheckEnable := enableToBool(res.HealthCheck)

	hcList := strings.Split(res.HealthCheckList, " ")
	if len(hcList) > 1 {
		hcList = hcList[:len(hcList)-1]
	}
	if !healtcheckEnable {
		hcList = []string{}
	}

	arguments := map[string]interface{}{
		"name":                    res.Mkey,
		"pool_type":               res.PoolType,
		"healtcheck_enable":       healtcheckEnable,
		"healtcheck_relationship": res.HealthCheckRelationship,
		"healtcheck_list":         hcList,
		"real_server_ssl_profile": res.RsProfile,
	}

	for arg, value := range arguments {
		err = d.Set(arg, value)
		if err != nil {
			return err
		}
	}

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

	req := gofortiadc.LoadbalancePool{
		Mkey:                    d.Get("name").(string),
		PoolType:                d.Get("pool_type").(string),
		HealthCheck:             boolToEnable(d.Get("healtcheck_enable").(bool)),
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
