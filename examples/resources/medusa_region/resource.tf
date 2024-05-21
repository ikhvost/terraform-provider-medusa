resource "medusa_region" "my-region" {
  name                  = "my-region"
  currency_code         = "usd"
  payment_providers     = ["manual"]
  fulfillment_providers = ["manual"]
  countries             = ["gb", "de", "dk", "se", "fr", "es", "it"]
  tax_rate              = 0
}
