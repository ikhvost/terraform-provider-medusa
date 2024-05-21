package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ikhvost/medusajs-go-sdk/medusa"
	"github.com/ikhvost/terraform-provider-medusa/internal/utils"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"

	types "github.com/hashicorp/terraform-plugin-framework/types"
	basetypes "github.com/oapi-codegen/runtime/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &medusaProvider{}
)

type OptionFunc func(p *medusaProvider)

func WithRetryableClient(retries int) OptionFunc {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = retries

	return func(p *medusaProvider) {
		p.httpClient = retryClient.StandardClient()
	}
}

func WithDebugClient() OptionFunc {
	return func(p *medusaProvider) {
		p.httpClient.Transport = NewDebugTransport(p.httpClient.Transport)
	}
}

func WithRecorderClient(file string, mode recorder.Mode) (OptionFunc, func() error) {
	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       file,
		Mode:               mode,
		SkipRequestLatency: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	//Strip all fields we are not interested in
	hook := func(i *cassette.Interaction) error {
		i.Response.Headers = utils.CleanHeaders(i.Response.Headers, "Content-Type")
		i.Request.Headers = utils.CleanHeaders(i.Request.Headers)
		return nil
	}
	r.AddHook(hook, recorder.AfterCaptureHook)

	stop := func() error {
		return r.Stop()
	}

	return func(p *medusaProvider) {
		p.httpClient = r.GetDefaultClient()
	}, stop
}

// New is a helper function to simplify provider server and testing implementation.
func New(opts ...OptionFunc) provider.Provider {
	tp := http.DefaultTransport

	var p = &medusaProvider{
		httpClient: &http.Client{Transport: tp},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// medusaProvider is the provider implementation.
type medusaProvider struct {
	httpClient *http.Client
}

// medusaProviderModel maps provider schema data to a Go type.
type medusaProviderModel struct {
	URL      types.String `tfsdk:"url"`
	Email    types.String `tfsdk:"email"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *medusaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "medusa"
}

// Schema defines the provider-level schema for configuration data.
func (p *medusaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Medusa API.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "Admin API base URL",
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "Admin user email",
				Required:    true,
				Sensitive:   true,
			},
			"password": schema.StringAttribute{
				Description: "Admin user password",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a MedusaJS API client for data sources and resources.
func (p *medusaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Medusa client")

	// Retrieve provider data from configuration
	var config medusaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	url := os.Getenv("MEDUSA_URL")
	email := os.Getenv("MEDUSA_ADMIN_EMAIL")
	password := os.Getenv("MEDUSA_ADMIN_PASSWORD")

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	if !config.Email.IsNull() {
		email = config.Email.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "medusa_url", url)
	ctx = tflog.SetField(ctx, "medusa_email", email)
	ctx = tflog.SetField(ctx, "medusa_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "medusa_email")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "medusa_password")

	tflog.Debug(ctx, "Creating Medusa client")

	client, err := medusa.NewClientWithResponses(url)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Medusa API Client", err.Error())
	}

	token, err := login(client, medusa.PostTokenJSONRequestBody{
		Email:    basetypes.Email(email),
		Password: password,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Login to Medusa API", err.Error())
	}

	tokenProvider, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Storyblok API Client", err.Error())
	}

	client, err = medusa.NewClientWithResponses(url, medusa.WithRequestEditorFn(tokenProvider.Intercept))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Medusa API Client",
			"An unexpected error occurred when creating the Medusa API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Medusa Error: "+err.Error(),
		)
		return
	}

	// Make the Medusa client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Medusa client", map[string]any{"success": true})
}

func login(client *medusa.ClientWithResponses, credentials medusa.PostTokenJSONRequestBody) (string, error) {
	resp, err := client.PostTokenWithResponse(context.Background(), credentials)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("login failed: %s", resp.Status())
	}

	return *resp.JSON200.AccessToken, nil
}

// DataSources defines the data sources implemented in the provider.
func (p *medusaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *medusaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRegionResource,
	}
}
