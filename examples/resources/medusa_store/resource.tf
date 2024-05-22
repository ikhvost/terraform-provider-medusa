resource "medusa_store" "my-store" {
  name                  = "my-store"
  default_currency_code = "usd"
  currencies            = ["usd"]
}
