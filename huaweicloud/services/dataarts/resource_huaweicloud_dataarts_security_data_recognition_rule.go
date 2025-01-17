package dataarts

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// API: DataArtsStudio POST /v1/{project_id}/security/data-classification/rule
// API: DataArtsStudio DELETE /v1/{project_id}/security/data-classification/rule/{id}
// API: DataArtsStudio GET /v1/{project_id}/security/data-classification/rule/{id}
// API: DataArtsStudio PUT /v1/{project_id}/security/data-classification/rule/{id}
func ResourceSecurityRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityRuleCreate,
		ReadContext:   resourceSecurityRuleRead,
		UpdateContext: resourceSecurityRuleUpdate,
		DeleteContext: resourceSecurityRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDataArtsStudioImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secrecy_level_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"builtin_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"content_expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"column_expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment_expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secrecy_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secrecy_level_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateOrUpdateSecurityRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"rule_type":          d.Get("rule_type").(string),
		"secrecy_level_id":   d.Get("secrecy_level_id").(string),
		"name":               d.Get("name").(string),
		"method":             utils.ValueIngoreEmpty(d.Get("method").(string)),
		"content_expression": utils.ValueIngoreEmpty(d.Get("content_expression").(string)),
		"column_expression":  utils.ValueIngoreEmpty(d.Get("column_expression").(string)),
		"commit_expression":  utils.ValueIngoreEmpty(d.Get("comment_expression").(string)),
		"builtin_rule_id":    utils.ValueIngoreEmpty(d.Get("builtin_rule_id").(string)),
		"description":        utils.ValueIngoreEmpty(d.Get("description").(string)),
		"category_id":        utils.ValueIngoreEmpty(d.Get("category_id").(string)),
	}
	return bodyParams
}

func resourceSecurityRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	createRuleHttpUrl := "v1/{project_id}/security/data-classification/rule"
	createRuleProduct := "dataarts"

	createRuleClient, err := conf.NewServiceClient(createRuleProduct, region)

	if err != nil {
		return diag.Errorf("error creating DataArts Studio V1 client: %s", err)
	}

	createRulePath := createRuleClient.Endpoint + createRuleHttpUrl
	createRulePath = strings.ReplaceAll(createRulePath, "{project_id}", createRuleClient.ProjectID)

	createRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	createRuleOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateSecurityRuleBodyParams(d))
	createRuleResp, err := createRuleClient.Request("POST", createRulePath, &createRuleOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Security data recognition rule: %s", err)
	}

	createRuleRespBody, err := utils.FlattenResponse(createRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("uuid", createRuleRespBody)
	if err != nil || id == nil {
		return diag.Errorf("error creating DataArts Security data recognition rule: ID is not found in API response")
	}

	d.SetId(id.(string))

	return resourceSecurityRuleRead(ctx, d, meta)
}

func resourceSecurityRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	getRuleHttpUrl := "v1/{project_id}/security/data-classification/rule/{id}"
	getRuleProduct := "dataarts"

	getRuleClient, err := cfg.NewServiceClient(getRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio V1 Client: %s", err)
	}

	getRulePath := getRuleClient.Endpoint + getRuleHttpUrl
	getRulePath = strings.ReplaceAll(getRulePath, "{project_id}", getRuleClient.ProjectID)
	getRulePath = strings.ReplaceAll(getRulePath, "{id}", d.Id())

	getRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	getRuleResp, err := getRuleClient.Request("GET", getRulePath, &getRuleOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, parseSecurityRuleError(err), "error retrieving DataArts Security data recognition rule")
	}

	getRuleRespBody, err := utils.FlattenResponse(getRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rule_type", utils.PathSearch("rule_type", getRuleRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRuleRespBody, nil)),
		d.Set("method", utils.PathSearch("method", getRuleRespBody, nil)),
		d.Set("content_expression", utils.PathSearch("content_expression", getRuleRespBody, nil)),
		d.Set("column_expression", utils.PathSearch("column_expression", getRuleRespBody, nil)),
		d.Set("comment_expression", utils.PathSearch("commit_expression", getRuleRespBody, nil)),
		d.Set("builtin_rule_id", utils.PathSearch("builtin_rule_id", getRuleRespBody, nil)),
		d.Set("category_id", utils.PathSearch("category_id", getRuleRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRuleRespBody, nil)),
		d.Set("secrecy_level", utils.PathSearch("secrecy_level", getRuleRespBody, nil)),
		d.Set("secrecy_level_num", utils.PathSearch("secrecy_level_num", getRuleRespBody, nil)),
		d.Set("enable", utils.PathSearch("enable", getRuleRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("created_at", getRuleRespBody, nil).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("updated_at", getRuleRespBody, nil).(float64))/1000, false)),
		d.Set("created_by", utils.PathSearch("created_by", getRuleRespBody, nil)),
		d.Set("updated_by", utils.PathSearch("updated_by", getRuleRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// The example of error message is: {"error_code": "DLS.4106","error_msg": "Rule is not exist."}
func parseSecurityRuleError(err error) error {
	var errCode golangsdk.ErrDefault400
	if errors.As(err, &errCode) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return err
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return err
		}

		if errorCode == "DLS.4106" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return err
}

func resourceSecurityRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateRuleHttpUrl := "v1/{project_id}/security/data-classification/rule/{id}"
	updateRuleProduct := "dataarts"

	updateRuleClient, err := cfg.NewServiceClient(updateRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio V1 Client: %s", err)
	}
	updateRulePath := updateRuleClient.Endpoint + updateRuleHttpUrl
	updateRulePath = strings.ReplaceAll(updateRulePath, "{project_id}", updateRuleClient.ProjectID)
	updateRulePath = strings.ReplaceAll(updateRulePath, "{id}", d.Id())

	updateRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	updateRuleOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateSecurityRuleBodyParams(d))
	_, err = updateRuleClient.Request("PUT", updateRulePath, &updateRuleOpt)
	if err != nil {
		return diag.Errorf("error updating DataArts Security data recognition rule: %s", err)
	}

	return resourceSecurityRuleRead(ctx, d, meta)
}

func resourceSecurityRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	deleteRuleHttpUrl := "v1/{project_id}/security/data-classification/rule/{id}"
	deleteRuleProduct := "dataarts"

	deleteRuleClient, err := cfg.NewServiceClient(deleteRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio V1 Client: %s", err)
	}
	deleteRulePath := deleteRuleClient.Endpoint + deleteRuleHttpUrl
	deleteRulePath = strings.ReplaceAll(deleteRulePath, "{project_id}", deleteRuleClient.ProjectID)
	deleteRulePath = strings.ReplaceAll(deleteRulePath, "{id}", d.Id())

	deleteRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	_, err = deleteRuleClient.Request("DELETE", deleteRulePath, &deleteRuleOpt)
	if err != nil {
		return diag.Errorf("error deleting DataArts Security data recognition rule: %s", err)
	}

	return nil
}
