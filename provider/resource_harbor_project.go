package provider

import (
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"projectid": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: false,
			},
		},
	}
}

func mapDataToProjectReq(d *schema.ResourceData, project *harbor.ProjectReq) error {
	project.ProjectName = d.Get("name").(string)

	project.Metadata = harbor.ProjectMetadata{
		Public: strconv.FormatBool(d.Get("public").(bool)),
	}

	return nil
}

func mapProjectToData(d *schema.ResourceData, project *harbor.Project) error {
	err := d.Set("name", project.Name)
	if err != nil {
		return err
	}
	public, err := strconv.ParseBool(project.Metadata.Public)
	if err != nil {
		return err
	}
	err = d.Set("public", public)
	if err != nil {
		return err
	}
	err = d.Set("projectid", strconv.Itoa(int(project.ProjectID)))
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
	log.Printf("[DEBUG] RESOURCE PROJECT DELETE")
	client := meta.(*harbor.Client)

	repos, err := client.GetRepositories(d.Get("projectid").(string))
	if err != nil {
		return err
	}

	err = client.DeleteRepositories(repos)
	if err != nil {
		return err
	}

	err = client.DeleteProject(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
