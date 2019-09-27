package fortiadc

import (
	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFortiadcLoadbalanceContentRoutingCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalanceContentRoutingConditionCreate,
		Read:   resourceFortiadcLoadbalanceContentRoutingConditionRead,
		Update: resourceFortiadcLoadbalanceContentRoutingConditionUpdate,
		Delete: resourceFortiadcLoadbalanceContentRoutingConditionDelete,

		Schema: map[string]*schema.Schema{
			"content_routing": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"reverse": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceFortiadcLoadbalanceContentRoutingConditionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	reverse := "disable"
	if d.Get("reverse").(bool) {
		reverse = "enable"
	}

	req := gofortiadc.LoadbalanceContentRoutingCondition{
		Mkey:    "",
		Object:  d.Get("object").(string),
		Type:    d.Get("type").(string),
		Content: d.Get("content").(string),
		Reverse: reverse,
	}

	err := client.LoadbalanceCreateContentRoutingCondition(d.Get("content_routing").(string), req)
	if err != nil {
		return err
	}

	id, err := client.LoadbalanceGetContentRoutingConditionID(d.Get("content_routing").(string), req)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceFortiadcLoadbalanceContentRoutingConditionRead(d, m)
}

func resourceFortiadcLoadbalanceContentRoutingConditionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	res, err := client.LoadbalanceGetContentRoutingCondition(d.Get("content_routing").(string), d.Id())
	if err != nil {
		return err
	}

	reverse := false
	if res.Reverse == "enable" {
		reverse = true
	}

	arguments := map[string]interface{}{
		"object":  res.Object,
		"type":    res.Type,
		"content": res.Content,
		"reverse": reverse,
	}

	for arg, value := range arguments {
		err = d.Set(arg, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceFortiadcLoadbalanceContentRoutingConditionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	reverse := "disable"
	if d.Get("reverse").(bool) {
		reverse = "enable"
	}

	req := gofortiadc.LoadbalanceContentRoutingCondition{
		Mkey:    d.Id(),
		Object:  d.Get("object").(string),
		Type:    d.Get("type").(string),
		Content: d.Get("content").(string),
		Reverse: reverse,
	}

	err := client.LoadbalanceUpdateContentRoutingCondition(d.Get("content_routing").(string), req)
	if err != nil {
		return err
	}

	return resourceFortiadcLoadbalanceContentRoutingConditionRead(d, m)
}

func resourceFortiadcLoadbalanceContentRoutingConditionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	return client.LoadbalanceDeleteContentRoutingCondition(d.Get("content_routing").(string), d.Id())
}
