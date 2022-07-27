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