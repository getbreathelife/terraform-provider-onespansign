package provider

import (
	"errors"
	"fmt"
	"testing"

	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testImg = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAoAAAANCAYAAACQN/8FAAAABGdBTUEAALGPC" +
	"/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAhGVYSWZNTQAqAAAACAAFARIAAwAAAAEA" +
	"AQAAARoABQAAAAEAAABKARsABQAAAAEAAABSASgAAwAAAAEAAgAAh2kABAAAAAEAAABaAAAAAAAAASAAAAABAAABIAAAAAEAA6A" +
	"BAAMAAAABAAEAAKACAAQAAAABAAAACqADAAQAAAABAAAADQAAAADcXFczAAAACXBIWXMAACxLAAAsSwGlPZapAAABWWlUWHRYTU" +
	"w6Y29tLmFkb2JlLnhtcAAAAAAAPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgN" +
	"i4wLjAiPgogICA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMi" +
	"PgogICAgICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgICAgICAgICB4bWxuczp0aWZmPSJodHRwOi8vbnMuYWR" +
	"vYmUuY29tL3RpZmYvMS4wLyI+CiAgICAgICAgIDx0aWZmOk9yaWVudGF0aW9uPjE8L3RpZmY6T3JpZW50YXRpb24+CiAgICAgID" +
	"wvcmRmOkRlc2NyaXB0aW9uPgogICA8L3JkZjpSREY+CjwveDp4bXBtZXRhPgoZXuEHAAAByElEQVQoFS1Rz2sTQRR+b2Z2N7Ek2" +
	"mpyyEHi1RQ8VA+lBXuyzcGbK3gTLxZJpf+A0IOIV7HtoVcRweCpWBoUBCWe7HEREVpFLNg1lCbZbfbHzPNNcGB4v775+N43uNjq" +
	"joR0PakKYIwGncWvyhfO3m2vNVIgQkAk4KNQOJ7J019G5+8BqFYs1+70//4+4tnqzP09tQeQWaBQThFI4E5nffYe3+Zp/7DP/aY" +
	"d1gAc33/t2lxlySBHgHlbACAZ070OSka22t66Gtvo+yRxaaX7UDoTz/Ik3gWEb6wp4ofGEEyxlBPH4Iu3m7Nfcan1aUG6pQ9Cen" +
	"YRYM2WBIxOwS2eg9HwiKXBnCCQz/M01kncu0Eaq9pElygZ1MHE1dHwzy3HKwMaeqqcQnk6Pz152dmYfzemAgj/RxveLLY+f2Tt0" +
	"4rYO97a2K6/FrgQBDpsVLAShNRu39ZsIpsJmUIUDKHEAkeH+2rQqJgKVMX+ZMkarS2IzZAqS4dc4xUL3N66ObbD5vY0V757hsLL" +
	"QBAr/qDNM5P1B6zlC88CbhaYgdkw0Sa8Vjpfr0W9g0fqYuqs/jz+IVjCsnInZojGchmHoJPIDHsHTzobc4//ASnLxOSzgBDYAAA" +
	"AAElFTkSuQmCC"

func TestAccResourceSigningLogos(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "onespansign_account_signing_logos" "foo" {
					logo {
						language = "en"
						image = "%s"
					}

					logo {
						language = "fr"
						image = "%s"
					}
				}
				`, testImg, testImg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onespansign_account_signing_logos.foo", "logo.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"onespansign_account_signing_logos.foo", "logo.*", map[string]string{
							"language": "en",
							"image":    testImg,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"onespansign_account_signing_logos.foo", "logo.*", map[string]string{
							"language": "fr",
							"image":    testImg,
						}),
					testAccCheckSigningLogosResourceMatches([]ossign.SigningLogo{
						{
							Language: "en",
							Image:    testImg,
						},
						{
							Language: "fr",
							Image:    testImg,
						},
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "onespansign_account_signing_logos" "foo" {
					logo {
						language = "en"
						image = "%s"
					}
				}
				`, testImg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onespansign_account_signing_logos.foo", "logo.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"onespansign_account_signing_logos.foo", "logo.*", map[string]string{
							"language": "en",
							"image":    testImg,
						}),
					testAccCheckSigningLogosResourceMatches([]ossign.SigningLogo{
						{
							Language: "en",
							Image:    testImg,
						},
					}),
				),
			},
			{
				ResourceName:      "onespansign_account_signing_logos.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: `resource "onespansign_account_signing_logos" "foo" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onespansign_account_signing_logos.foo", "logo.#", "0"),
					testAccCheckSigningLogosResourceMatches([]ossign.SigningLogo{}),
				),
			},
		},
	})
}

func testAccCheckSigningLogosResourceMatches(m []ossign.SigningLogo) resource.TestCheckFunc {
	return func(*terraform.State) error {
		c := getTestApiClient()

		l, err := c.GetAccountSigningLogos()
		if err != nil {
			return err.GetError()
		}

		var match bool

		for _, v1 := range l {
			match = false

			for _, v2 := range m {
				if cmp.Equal(v1, v2) {
					match = true
					break
				}
			}

			if !match {
				return errors.New("Signing logos resource does not match expectation")
			}
		}

		return nil
	}
}
