package fortiadc

import (
	"github.com/Ouest-France/gofortiadc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceFortiadcLoadbalanceContentRewriting() *schema.Resource {
	return &schema.Resource{
		Create: resourceFortiadcLoadbalanceContentRewritingCreate,
		Read:   resourceFortiadcLoadbalanceContentRewritingRead,
		Update: resourceFortiadcLoadbalanceContentRewritingUpdate,
		Delete: resourceFortiadcLoadbalanceContentRewritingDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"request",
						"response",
					}, false),
			},
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"rewrite_http_header",
						"redirect",
						"send-403-forbidden",
						"add_http_header",
						"delete_http_header",
					}, false),
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"host_match": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"url_match": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"referer_match": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"redirect": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"header_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceFortiadcLoadbalanceContentRewritingCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalanceContentRewriting{
		Mkey:           d.Get("name").(string),
		ActionType:     d.Get("action_type").(string),
		Action:         d.Get("action").(string),
		URLStatus:      "disable",
		URLContent:     d.Get("url_match").(string),
		RefererStatus:  "disable",
		RefererContent: d.Get("referer_match").(string),
		Redirect:       d.Get("redirect").(string),
		Location:       d.Get("location").(string),
		HeaderName:     "header-name",
		Comments:       d.Get("comment").(string),
		HostStatus:     "disable",
		HostContent:    d.Get("host_match").(string),
	}

	if d.Get("url_match").(string) != "" {
		req.URLStatus = "enable"
	}
	if d.Get("referer_match").(string) != "" {
		req.RefererStatus = "enable"
	}
	if d.Get("host_match").(string) != "" {
		req.HostStatus = "enable"
	}

	err := client.LoadbalanceCreateContentRewriting(req)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return resourceFortiadcLoadbalanceContentRewritingRead(d, m)
}

func resourceFortiadcLoadbalanceContentRewritingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	res, err := client.LoadbalanceGetContentRewriting(d.Id())
	if err != nil {
		return err
	}

	arguments := map[string]interface{}{
		"action_type": res.ActionType,
		"action":      res.Action,
		"host_match":  res.HostContent,
	}

	if res.HeaderName == "header-name" {
		arguments["header_name"] = ""
	} else {
		arguments["header_name"] = res.HeaderName
	}

	if res.Location == "http://" {
		arguments["location"] = ""
	} else {
		arguments["location"] = res.Location
	}

	if res.Redirect == "redirect" {
		arguments["redirect"] = ""
	} else {
		arguments["redirect"] = res.Redirect
	}

	if res.RefererContent == "http://" {
		arguments["referer_match"] = ""
	} else {
		arguments["referer_match"] = res.RefererContent
	}

	if res.URLContent == "/url" {
		arguments["url_match"] = ""
	} else {
		arguments["url_match"] = res.URLContent
	}

	if res.Comments == "comments" {
		arguments["comment"] = ""
	} else {
		arguments["comment"] = res.Comments
	}

	for arg, value := range arguments {
		err = d.Set(arg, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceFortiadcLoadbalanceContentRewritingUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	req := gofortiadc.LoadbalanceContentRewriting{
		Mkey:           d.Get("name").(string),
		ActionType:     d.Get("action_type").(string),
		Action:         d.Get("action").(string),
		URLStatus:      "disable",
		URLContent:     d.Get("url_match").(string),
		RefererStatus:  "disable",
		RefererContent: d.Get("referer_match").(string),
		Redirect:       d.Get("redirect").(string),
		Location:       d.Get("location").(string),
		HeaderName:     "header-name",
		Comments:       d.Get("comment").(string),
		HostStatus:     "disable",
		HostContent:    d.Get("host_match").(string),
	}

	if d.Get("url_match").(string) != "" {
		req.URLStatus = "enable"
	}
	if d.Get("referer_match").(string) != "" {
		req.RefererStatus = "enable"
	}
	if d.Get("host_match").(string) != "" {
		req.HostStatus = "enable"
	}

	err := client.LoadbalanceUpdateContentRewriting(req)
	if err != nil {
		return err
	}

	return resourceFortiadcLoadbalanceContentRewritingRead(d, m)
}

func resourceFortiadcLoadbalanceContentRewritingDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gofortiadc.Client)

	return client.LoadbalanceDeleteContentRewriting(d.Id())
}
