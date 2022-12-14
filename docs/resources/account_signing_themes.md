---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "onespansign_account_signing_themes Resource - terraform-provider-onespansign"
subcategory: ""
description: |-
  OneSpan Sign account's customized signing themes.
  Please note that this resource is a singleton, which means that only one instance of this resource should
  exist. Having multiple instances may produce unexpected result.
---

# onespansign_account_signing_themes (Resource)

OneSpan Sign account's customized signing themes.
		
Please note that this resource is a singleton, which means that only one instance of this resource should
exist. Having multiple instances may produce unexpected result.

## Example Usage

```terraform
resource "onespansign_account_signing_themes" "example" {
  theme {
    name                      = "default"
    primary                   = "#129EFB"
    success                   = "#474DBE"
    warning                   = "#708AB6"
    error                     = "#F80B00"
    info                      = "#968949"
    signature_button          = "#F7003D"
    optional_signature_button = "#F8DCD3"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `theme` (Block Set, Min: 1, Max: 1) Customized signing theme for the account. (see [below for nested schema](#nestedblock--theme))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--theme"></a>
### Nested Schema for `theme`

Required:

- `error` (String) Error notification color hex code.
- `info` (String) Info notification color hex code.
- `name` (String) Name of the theme.
- `optional_signature_button` (String) Color hex code for the optional signature buttons.
- `primary` (String) Primary color hex code.
- `signature_button` (String) Color hex code for the required signature buttons.
- `success` (String) Success notification color hex code.
- `warning` (String) Warning notification color hex code.


