package fortiadc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFortiadcLoadbalancePoolMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalancePoolMemberCreate,
		Read:   resourceFortiadcLoadbalancePoolMemberRead,
		Update: resourceFortiadcLoadbalancePoolMemberUpdate,
		Delete: resourceFortiadcLoadbalancePoolMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				client := m.(*gofortiadc.Client)

				// Split id to pool and member names
				idParts := strings.Split(d.Id(), ".")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected POOL.MEMBER", d.Id())
				}
				pool := idParts[0]
				member := idParts[1]

				// Read member ID
				mkey, err := client.LoadbalanceGetPoolMemberID(pool, member)
				if err != nil {
					return nil, err
				}

				// Set pool and member
				err = d.Set("pool", pool)
				if err != nil {
					return nil, err
				}
				d.SetId(mkey)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pool": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enable",
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"conn_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"recover": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"warmup": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"warmrate": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"conn_rate_limit": {
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
	if errors.Is(err, gofortiadc.ErrNotFound) {
		// If not found, remove from state
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	arguments := map[string]interface{}{
		"name":   res.RealServerID,
		"status": res.Status,
	}

	port, err := strconv.ParseInt(res.Port, 10, 64)
	if err != nil {
		return err
	}
	arguments["port"] = port

	weight, err := strconv.ParseInt(res.Weight, 10, 64)
	if err != nil {
		return err
	}
	arguments["weight"] = weight

	connLimit, err := strconv.ParseInt(res.Connlimit, 10, 64)
	if err != nil {
		return err
	}
	arguments["conn_limit"] = connLimit

	recover, err := strconv.ParseInt(res.Recover, 10, 64)
	if err != nil {
		return err
	}
	arguments["recover"] = recover

	warmup, err := strconv.ParseInt(res.Warmup, 10, 64)
	if err != nil {
		return err
	}
	arguments["warmup"] = warmup

	warmrate, err := strconv.ParseInt(res.Warmrate, 10, 64)
	if err != nil {
		return err
	}
	arguments["warmrate"] = warmrate

	connectionRateLimit, err := strconv.ParseInt(res.ConnectionRateLimit, 10, 64)
	if err != nil {
		return err
	}
	arguments["conn_rate_limit"] = connectionRateLimit

	for arg, value := range arguments {
		err = d.Set(arg, value)
		if err != nil {
			return err
		}
	}

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
