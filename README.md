## Example Usage
```terraform
terraform {
  required_providers{
    zoom ={
      version ="0.2"
      source = "hashicorp.com/edu/zoom"
    }
  }
}

data "zoom_user" "user" {
  email = "user@domain.com"
}

output "user" {
  value = data.zoom_user.user
}


resource "zoom_user" "user1" {
  email = "user@domain.com"
  first_name = "User_First_Name"
  last_name = "User_Last_Name"
  type = 1
}
```
## Argument Reference

* `email`(required, string)         - Email of the user.
* `first_name`(optional, string) - First Name of the User.
* `Last_name`(optional, string) - Last Name of the User.
*  `type`(optional, string) - Type of the User.
