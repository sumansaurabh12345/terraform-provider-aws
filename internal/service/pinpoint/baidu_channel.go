package pinpoint

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pinpoint"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKResource("aws_pinpoint_baidu_channel")
func ResourceBaiduChannel() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceBaiduChannelUpsert,
		ReadWithoutTimeout:   resourceBaiduChannelRead,
		UpdateWithoutTimeout: resourceBaiduChannelUpsert,
		DeleteWithoutTimeout: resourceBaiduChannelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceBaiduChannelUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).PinpointConn()

	applicationId := d.Get("application_id").(string)

	params := &pinpoint.BaiduChannelRequest{}

	params.Enabled = aws.Bool(d.Get("enabled").(bool))
	params.ApiKey = aws.String(d.Get("api_key").(string))
	params.SecretKey = aws.String(d.Get("secret_key").(string))

	req := pinpoint.UpdateBaiduChannelInput{
		ApplicationId:       aws.String(applicationId),
		BaiduChannelRequest: params,
	}

	_, err := conn.UpdateBaiduChannelWithContext(ctx, &req)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Pinpoint Baidu Channel for application %s: %s", applicationId, err)
	}

	d.SetId(applicationId)

	return append(diags, resourceBaiduChannelRead(ctx, d, meta)...)
}

func resourceBaiduChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).PinpointConn()

	log.Printf("[INFO] Reading Pinpoint Baidu Channel for application %s", d.Id())

	output, err := conn.GetBaiduChannelWithContext(ctx, &pinpoint.GetBaiduChannelInput{
		ApplicationId: aws.String(d.Id()),
	})
	if err != nil {
		if tfawserr.ErrCodeEquals(err, pinpoint.ErrCodeNotFoundException) {
			log.Printf("[WARN] Pinpoint Baidu Channel for application %s not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}

		return sdkdiag.AppendErrorf(diags, "getting Pinpoint Baidu Channel for application %s: %s", d.Id(), err)
	}

	d.Set("application_id", output.BaiduChannelResponse.ApplicationId)
	d.Set("enabled", output.BaiduChannelResponse.Enabled)
	// ApiKey and SecretKey are never returned

	return diags
}

func resourceBaiduChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).PinpointConn()

	log.Printf("[DEBUG] Deleting Pinpoint Baidu Channel for application %s", d.Id())
	_, err := conn.DeleteBaiduChannelWithContext(ctx, &pinpoint.DeleteBaiduChannelInput{
		ApplicationId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, pinpoint.ErrCodeNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Pinpoint Baidu Channel for application %s: %s", d.Id(), err)
	}
	return diags
}
