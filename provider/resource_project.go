package provider

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
			},

			"project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func setProjectData(d *schema.ResourceData, project *harbor.Project) error {
	d.SetId(project.Name)

	err := d.Set("project_id", project.ProjectId)
	if err != nil {
		return err
	}
	err = d.Set("name", project.Name)
	if err != nil {
		return err
	}
	err = d.Set("public", project.Metadata.Public)
	if err != nil {
		return err
	}
	return nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	projectName := d.Get("name").(string)
	project, err := client.GetProject(projectName)
	if err != nil {
		return err
	}

	return setProjectData(d, project)
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	project := &harbor.ProjectReq{
		ProjectName: d.Get("name").(string),
		Metadata: harbor.ProjectMetadata{
			Public: d.Get("public").(string),
		},
	}

	err := client.NewProject(project)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	project := &harbor.ProjectReq{
		ProjectName: d.Get("name").(string),
		Metadata: harbor.ProjectMetadata{
			Public: d.Get("public").(string),
		},
	}

	projectId := strconv.Itoa(d.Get("project_id").(int))

	err := client.UpdateProject(projectId, project)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	projectId := strconv.Itoa(d.Get("project_id").(int))

	err := client.DeleteProject(projectId)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
