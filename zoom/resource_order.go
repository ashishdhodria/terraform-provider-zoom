package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Type int `json:"type"`
}

type Create struct {
	Action string    `json:"action"`
	User   User `json:"user_info"`
}

func resourceOrder() *schema.Resource {
	return &schema.Resource{
		UpdateContext: resourceOrderUpdate,
		CreateContext: resourceOrderCreate,
		ReadContext:   resourceOrderRead,
		DeleteContext: resourceOrderDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"first_name":&schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type: schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceOrderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	url := "https://api.zoom.us/v2/users"

	bearer := "Bearer " + "ZOOM_TOKEN"
	email := d.Get("email").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	typeT := d.Get("type").(int)

	values := Create{
		Action: "create",
		User: User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Type:      typeT,
		},
	}
	jsonStr, err := json.Marshal(values)

	if err != nil {
		return diag.FromErr(err)
	}


	// Create a new request using http
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonStr)))

	req.Header.Set("Content-Type", "application/json")
	// add auth header to the req
	req.Header.Set("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}

	defer resp.Body.Close()

	resourceOrderRead(ctx, d, m)
	return diags

}

func resourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	var diags diag.Diagnostics
	user := User{}
	Id := d.Get("email").(string)

	req, err := http.NewRequest("GET", fmt.Sprintf( "https://api.zoom.us/v2/users/%s", Id), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	bearer := "Bearer " + "ZOOM_TOKEN"
	req.Header.Add("Authorization", bearer)

	r, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}
	err = json.NewDecoder(r.Body).Decode(&user)

	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("email", user.Email)
	d.Set("type", user.Type)
	defer r.Body.Close()
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags

}

func resourceOrderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var _ diag.Diagnostics

	Id := d.Get("email").(string)

	url := "https://api.zoom.us/v2/users/" + Id

	bearer := "Bearer " + "ZOOM_TOKEN"
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)

	values := User{
		FirstName: firstName,
		LastName:  lastName,
	}

	jsonStr, err := json.Marshal(values)

	if err != nil {
		return diag.FromErr(err)
	}

	// Create a new request using http
	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(jsonStr)))

	req.Header.Set("Content-Type", "application/json")

	// add auth header to the req
	req.Header.Set("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error on response: \n", err)
	}

	defer resp.Body.Close()

	d.Set("last_updated", time.Now().Format(time.RFC850))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceOrderRead(ctx, d, m)
}

func resourceOrderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	Id := d.Get("email").(string)

	url := "https://api.zoom.us/v2/users/" + Id

	bearer := "Bearer " + 

	// Create a new request using http
	req, err := http.NewRequest("DELETE", url, nil)

	req.Header.Set("Content-Type", "application/json")

	// add auth header to the req
	req.Header.Set("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error on response: \n", err)
	}

	defer resp.Body.Close()

	d.Set("status ", resp.Status)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
