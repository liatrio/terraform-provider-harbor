package provider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "ID of the project the robot account corresponds to in the form '/projects/${ID_NUMBER}'",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^/projects/[0-9]+$`), "validation error: project_id should be of the form '/projects/${ID_NUMBER}'"),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the robot account, beginning with 'robot$'",
				ValidateFunc: validation.All(
					validation.StringMatch(
						regexp.MustCompile(`^robot\$[^~#$%]+`),
						"validation error: name must begin with 'robot$' and must otherwise not include the special characters(~#$%)",
					),
					validation.StringLenBetween(1, 255),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "A description of this robot account",
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this robot account is disabled.",
			},
			"access": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"push", "pull"}, false),
						},
						"resource": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"image", "helm-chart"}, false),
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

func mapRobotAccountToData(d *schema.ResourceData, robot *harbor.RobotAccount) error {
	err := d.Set("name", robot.Name)
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

func mapRobotAccountPostRepToData(d *schema.ResourceData, robot *harbor.RobotAccountPostRep) error {
	err := d.Set("token", robot.Token)
	if err != nil {
		return err
	}
	err = d.Set("name", robot.Name)
	if err != nil {
		return err
	}
	return nil
}

func mapDataToRobotAccountCreate(d *schema.ResourceData, robot *harbor.RobotAccountCreate) {
	access := &[]harbor.RobotAccountAccess{}
	mapDataToRobotAccountAccess(d, access)

	robot.Name = strings.Replace(d.Get("name").(string), "robot$", "", 1)
	robot.Description = d.Get("description").(string)
	robot.Access = *access
}

func mapDataToRobotAccountAccess(d *schema.ResourceData, accessList *[]harbor.RobotAccountAccess) {
	v, ok := d.GetOk("access")
	if !ok {
		return
	}

	projectID := d.Get("project_id").(string)
	projectID = strings.Replace(projectID, "projects", "project", 1)

	for _, dataAccess := range v.(*schema.Set).List() {
		dataAccess := dataAccess.(map[string]interface{})
		access := harbor.RobotAccountAccess{}
		if dataAccess["resource"].(string) == "image" {
			access.Action = dataAccess["action"].(string)
			access.Resource = fmt.Sprintf("%s/repository", projectID)
		} else if dataAccess["resource"].(string) == "helm-chart" {
			if dataAccess["action"].(string) == "pull" {
				access.Action = "read"
				access.Resource = fmt.Sprintf("%s/helm-chart", projectID)
			} else if dataAccess["action"].(string) == "push" {
				access.Action = "create"
				access.Resource = fmt.Sprintf("%s/helm-chart-version", projectID)
			}
		}
		*accessList = append(*accessList, access)
	}
}

func mapDataToRobotAccountUpdate(d *schema.ResourceData, robot *harbor.RobotAccountUpdate) {
	robot.Disabled = d.Get("disabled").(bool)
}

func resourceRobotAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	robot, err := client.GetRobotAccount(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	return mapRobotAccountToData(d, robot)
}

func resourceRobotAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	robot := &harbor.RobotAccountCreate{}

	mapDataToRobotAccountCreate(d, robot)

	body, location, err := client.NewRobotAccount(d.Get("project_id").(string), robot)
	if err != nil {
		return err
	}

	d.SetId(location)
	err = mapRobotAccountPostRepToData(d, body)
	if err != nil {
		return err
	}

	return resourceRobotAccountRead(d, meta)
}

func resourceRobotAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	robot := &harbor.RobotAccountUpdate{}

	mapDataToRobotAccountUpdate(d, robot)

	err := client.UpdateRobotAccount(d.Id(), robot)
	if err != nil {
		return err
	}

	return resourceRobotAccountRead(d, meta)
}

func resourceRobotAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	err := client.DeleteRobotAccount(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.SetId("")
	return nil
}
