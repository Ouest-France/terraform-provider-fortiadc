package fortiadc

import (
	"errors"

	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFortiadcLoadbalanceContentRouting() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalanceContentRoutingCreate,
		Read:   resourceFortiadcLoadbalanceContentRoutingRead,
		Update: resourceFortiadcLoadbalanceContentRoutingUpdate,
		Delete: resourceFortiadcLoadbalanceContentRoutingDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "l7-content-routing",
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "comments",
			},
			"ipv4": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.0.0/0",
			},
			"ipv6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "::/0",
			},
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFortiadcLoadbalanceContentRoutingCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalanceContentRouting{
		Mkey:                  d.Get("name").(string),
		Type:                  d.Get("type").(string),
		Comments:              d.Get("comment").(string),
		IP:                    d.Get("ipv4").(string),
		IP6:                   d.Get("ipv6").(string),
		Pool:                  d.Get("pool").(string),
		PacketFwdMethod:       "inherit",
		SourcePoolList:        "",
		Persistence:           "",
		PersistenceInherit:    "enable",
		Method:                "",
		MethodInherit:         "enable",
		ConnectionPool:        "",
		ConnectionPoolInherit: "enable",
		ScheduleList:          "disable",
		SchedulePoolList:      "",
	}

	err := client.LoadbalanceCreateContentRouting(req)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return resourceFortiadcLoadbalanceContentRoutingRead(d, m)
}

func resourceFortiadcLoadbalanceContentRoutingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	res, err := client.LoadbalanceGetContentRouting(d.Id())
	if errors.Is(err, gofortiadc.ErrNotFound) {
		// If not found, remove from state
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	arguments := map[string]interface{}{
		"type":    res.Type,
		"comment": res.Comments,
		"ipv4":    res.IP,
		"ipv6":    res.IP6,
		"pool":    res.Pool,
	}

	for arg, value := range arguments {
		err = d.Set(arg, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceFortiadcLoadbalanceContentRoutingUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalanceContentRouting{
		Mkey:                  d.Get("name").(string),
		Type:                  d.Get("type").(string),
		Comments:              d.Get("comment").(string),
		IP:                    d.Get("ipv4").(string),
		IP6:                   d.Get("ipv6").(string),
		Pool:                  d.Get("pool").(string),
		PacketFwdMethod:       "inherit",
		SourcePoolList:        "",
		Persistence:           "",
		PersistenceInherit:    "enable",
		Method:                "",
		MethodInherit:         "enable",
		ConnectionPool:        "",
		ConnectionPoolInherit: "enable",
		ScheduleList:          "disable",
		SchedulePoolList:      "",
	}

	err := client.LoadbalanceUpdateContentRouting(req)
	if err != nil {
		return err
	}

	return resourceFortiadcLoadbalanceContentRoutingRead(d, m)
}

func resourceFortiadcLoadbalanceContentRoutingDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	return client.LoadbalanceDeleteContentRouting(d.Id())
}
