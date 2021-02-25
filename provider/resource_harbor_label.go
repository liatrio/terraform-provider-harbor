package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func resourceLabel() *schema.Resource {
	return &schema.Resource{
		Create: resourceLabelCreate,
		Read:   resourceLabelRead,
		Update: resourceLabelUpdate,
		Delete: resourceLabelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "Display name of the label.",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 125),
			},
			"color": {
				Type:         schema.TypeString,
				Description:  "Display color of the label.",
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^#[0-9A-F]{6}$`), "validation error: color should an RGB hex code of the form '#12AB34'"),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the label.",
				Optional:    true,
			},
			"project_id": {
				Type:         schema.TypeString,
				Description:  "If set, the project the label will be created under. If not set, label will be global.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^/projects/[0-9]+$`), "validation error: project_id should be of the form '/projects/${ID_NUMBER}'"),
			},
		},
	}
}

func mapDataToLabel(d *schema.ResourceData, label *harbor.Label) error {
	label.Name = d.Get("name").(string)
	label.Color = d.Get("color").(string)
	label.Description = d.Get("description").(string)
	projectPath := d.Get("project_id").(string)
	if projectPath == "" {
		label.ProjectID = 0
		label.Scope = "g"
	} else {
		ID, err := strconv.ParseInt(strings.Split(projectPath, "/")[2], 10, 64)
		if err != nil {
			return err
		}
		label.ProjectID = ID
		label.Scope = "p"
	}
	return nil
}

func mapLabelToData(d *schema.ResourceData, label *harbor.Label) error {
	err := d.Set("name", label.Name)
	if err != nil {
		return err
	}
	err = d.Set("color", label.Color)
	if err != nil {
		return err
	}
	err = d.Set("description", label.Description)
	if err != nil {
		return err
	}
	if label.Scope == "p" {
		err = d.Set("project_id", fmt.Sprintf("/projects/%d", label.ProjectID))
	} else {
		err = d.Set("project_id", "")
	}
	if err != nil {
		return err
	}
	return nil
}

func resourceLabelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	labelID := d.Id()
	label, err := client.GetLabel(labelID)
	if err != nil {
		return handleNotFoundError(err, d)
	}

	return mapLabelToData(d, label)
}

func resourceLabelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	label := &harbor.Label{}
	err := mapDataToLabel(d, label)
	if err != nil {
		return err
	}

	location, err := client.NewLabel(label)
	if err != nil {
		return err
	}

	d.SetId(location)
	return resourceLabelRead(d, meta)
}

func resourceLabelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	label := &harbor.Label{}
	err := mapDataToLabel(d, label)
	if err != nil {
		return err
	}

	err = client.UpdateLabel(d.Id(), label)
	if err != nil {
		return err
	}

	return resourceLabelRead(d, meta)
}

func resourceLabelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	err := client.DeleteLabel(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.SetId("")
	return nil
}
