resource "onespansign_account_signing_logos" "example" {
  logo {
    language = "en"
    image    = "data:image/png;base64,<BASE64_IMAGE_DATA>"
  }

  logo {
    language = "fr"
    image    = "data:image/png;base64,<BASE64_IMAGE_DATA>"
  }
}