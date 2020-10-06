package fortiadc

import (
	"errors"

	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceFortiadcLoadbalanceContentRewritingCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalanceContentRewritingConditionCreate,
		Read:   resourceFortiadcLoadbalanceContentRewritingConditionRead,
		Update: resourceFortiadcLoadbalanceContentRewritingConditionUpdate,
		Delete: resourceFortiadcLoadbalanceContentRewritingConditionDelete,

		Schema: map[string]*schema.Schema{
			"content_rewriting": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"http-host-header",
						"http-request-url",
						"http-referer-header",
						"ip-source-address",
						"http-location-header",
					}, false),
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "string",
				ValidateFunc: validation.StringInSlice(
					[]string{
						"string",
						"regular-expression",
					}, false),
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
			"ignore_case": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceFortiadcLoadbalanceContentRewritingConditionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalanceContentRewritingCondition{
		Mkey:       "",
		Object:     d.Get("object").(string),
		Type:       d.Get("type").(string),
		Content:    d.Get("content").(string),
		Reverse:    boolToEnable(d.Get("reverse").(bool)),
		Ignorecase: boolToEnable(d.Get("ignore_case").(bool)),
	}

	err := client.LoadbalanceCreateContentRewritingCondition(d.Get("content_rewriting").(string), req)
	if err != nil {
		return err
	}

	id, err := client.LoadbalanceGetContentRewritingConditionID(d.Get("content_rewriting").(string), req)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceFortiadcLoadbalanceContentRewritingConditionRead(d, m)
}

func resourceFortiadcLoadbalanceContentRewritingConditionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	res, err := client.LoadbalanceGetContentRewritingCondition(d.Get("content_rewriting").(string), d.Id())
	if errors.Is(err, gofortiadc.ErrNotFound) {
		// If not found, remove from state
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	arguments := map[string]interface{}{
		"object":      res.Object,
		"type":        res.Type,
		"content":     res.Content,
		"reverse":     enableToBool(res.Reverse),
		"ignore_case": d.Get("ignore_case").(bool),
	}

	for arg, value := range arguments {
		err = d.Set(arg, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceFortiadcLoadbalanceContentRewritingConditionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalanceContentRewritingCondition{
		Mkey:       d.Id(),
		Object:     d.Get("object").(string),
		Type:       d.Get("type").(string),
		Content:    d.Get("content").(string),
		Reverse:    boolToEnable(d.Get("reverse").(bool)),
		Ignorecase: boolToEnable(d.Get("ignore_case").(bool)),
	}

	err := client.LoadbalanceUpdateContentRewritingCondition(d.Get("content_rewriting").(string), req)
	if err != nil {
		return err
	}

	return resourceFortiadcLoadbalanceContentRewritingConditionRead(d, m)
}

func resourceFortiadcLoadbalanceContentRewritingConditionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	return client.LoadbalanceDeleteContentRewritingCondition(d.Get("content_rewriting").(string), d.Id())
}
