package provider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func resourceRobotAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceRobotAccountCreate,
		Read:   resourceRobotAccountRead,
		Update: resourceRobotAccountUpdate,
		Delete: resourceRobotAccountDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^robot\$.+`), "validation error: name must begin with 'robot$'"),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"robot_account_access": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceRobotAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	robot, err := client.GetRobotAccount(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	err = d.Set("name", robot.Name)
	if err != nil {
		return err
	}
	err = d.Set("description", robot.Description)
	if err != nil {
		return err
	}
	err = d.Set("disabled", robot.Disabled)
	if err != nil {
		return err
	}
	return nil
}

func resourceRobotAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	access := &[]harbor.RobotAccountAccess{}

	if v, ok := d.GetOk("robot_account_access"); ok {
		access = getRobotAccountAccessFromData(d.Get("project_id").(string), v.(*schema.Set).List())
	}

	robot := &harbor.RobotAccountCreate{
		Name:        strings.Replace(d.Get("name").(string), "robot$", "", 1),
		Description: d.Get("description").(string),
		Access:      *access,
	}

	body, location, err := client.NewRobotAccount(d.Get("project_id").(string), robot)
	if err != nil {
		return err
	}

	d.SetId(location)
	err = d.Set("token", body.Token)
	if err != nil {
		return err
	}
	err = d.Set("name", body.Name)
	if err != nil {
		return err
	}

	return resourceRobotAccountRead(d, meta)
}

func getRobotAccountAccessFromData(projectID string, data []interface{}) *[]harbor.RobotAccountAccess {
	accessList := make([]harbor.RobotAccountAccess, 0, len(data))
	resourcePrefix := strings.Replace(projectID, "projects", "project", 1)
	for _, d := range data {
		accessData := d.(map[string]interface{})
		access := harbor.RobotAccountAccess{}
		if accessData["resource"].(string) == "image" {
			access.Action = accessData["action"].(string)
			access.Resource = fmt.Sprintf("%s/repository", resourcePrefix)
		} else if accessData["resource"].(string) == "helm-chart" {
			if accessData["action"].(string) == "pull" {
				access.Action = "read"
				access.Resource = fmt.Sprintf("%s/helm-chart", resourcePrefix)
			} else if accessData["action"].(string) == "push" {
				access.Action = "create"
				access.Resource = fmt.Sprintf("%s/helm-chart-version", resourcePrefix)
			}
		}
		accessList = append(accessList, access)
	}

	return &accessList
}

func resourceRobotAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	robot := &harbor.RobotAccountUpdate{
		Disabled: d.Get("disabled").(bool),
	}

	robotID := d.Id()

	err := client.UpdateRobotAccount(robotID, robot)
	if err != nil {
		return err
	}

	return resourceRobotAccountRead(d, meta)
}

func resourceRobotAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	robotID := d.Id()

	err := client.DeleteRobotAccount(robotID)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
