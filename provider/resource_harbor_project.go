package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Display name of the project.",
				Required:    true,
				ForceNew:    true,
				ValidateFunc: validation.All(
					validation.StringMatch(
						regexp.MustCompile(`^[a-z0-9]([a-z0-9_.-]*[a-z0-9])?$`),
						"validation error:  project name should use lower case characters, numbers and ._- and must start and end with characters or numbers. '",
					),
					validation.StringLenBetween(1, 255),
				),
			},
			"public": {
				Type:        schema.TypeBool,
				Description: "When true, anyone has read permissions to repositories under this project.",
				Optional:    true,
				Default:     false,
			},
			"auto_scan": {
				Type:        schema.TypeBool,
				Description: "When true, it will auto scan on push.",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func mapDataToProjectReq(d *schema.ResourceData, project *harbor.ProjectReq) error {
	project.ProjectName = d.Get("name").(string)

	project.Metadata = harbor.ProjectMetadata{
		Public:   d.Get("public").(bool),
		AutoScan: d.Get("auto_scan").(bool),
	}
	return nil
}

func mapProjectToData(d *schema.ResourceData, project *harbor.Project) error {
	err := d.Set("name", project.Name)
	if err != nil {
		return err
	}
	err = d.Set("public", project.Metadata.Public)
	if err != nil {
		return err
	}
	err = d.Set("auto_scan", project.Metadata.AutoScan)
	if err != nil {
		return err
	}
	return nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	projectID := d.Id()
	project, err := client.GetProject(projectID)
	if err != nil {
		return handleNotFoundError(err, d)
	}

	return mapProjectToData(d, project)
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	project := &harbor.ProjectReq{}
	err := mapDataToProjectReq(d, project)
	if err != nil {
		return err
	}

	location, err := client.NewProject(project)
	if err != nil {
		return err
	}

	d.SetId(location)
	return resourceProjectRead(d, meta)
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	project := &harbor.ProjectReq{}
	err := mapDataToProjectReq(d, project)
	if err != nil {
		return err
	}

	err = client.UpdateProject(d.Id(), project)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	projectName := d.Get("name").(string)

	repos, err := client.GetRepositories(projectName)
	if err != nil {
		return handleNotFoundError(err, d)
	}

	if len(repos) > 0 {
		err = client.DeleteRepositories(projectName, repos)
		if err != nil {
			return err
		}
	}

	charts, err := client.GetCharts(projectName)
	// this can return a 404 if chartmuseum is disabled
	if err != nil && !harbor.ErrorIs404(err) {
		return err
	}

	if len(charts) > 0 {
		err = client.DeleteCharts(projectName, charts)
		if err != nil {
			return err
		}
	}

	err = client.DeleteProject(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.SetId("")
	return nil
}
