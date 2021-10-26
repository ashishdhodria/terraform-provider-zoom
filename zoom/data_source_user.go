package zoom

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
	"time"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"first_name":&schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
			},
		},
	}
}


func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	email := d.Get("email").(string)

	url := "https://api.zoom.us/v2/users/"+email

	bearer := "Bearer " + "ZOOM_TOKEN"
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add auth header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}

	defer resp.Body.Close()
	user := User{}
	err = json.NewDecoder(resp.Body).Decode(&user)

	oi := make(map[string]interface{})

	oi["first_name"] = user.FirstName
	oi["last_name"] = user.LastName
	oi["email"] = user.Email
	oi["type"] = user.Type

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("email", user.Email)
	d.Set("type", user.Type)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags

}

