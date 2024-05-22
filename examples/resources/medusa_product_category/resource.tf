resource "medusa_product_category" "my-parent-product-category" {
  name        = "my-parent-product-category"
  description = "my parent product category"
  handle      = "my-parent-handle"
  is_internal = false
  is_active   = false
}

resource "medusa_product_category" "my-child-product-category" {
  name               = "my-child-product-category"
  description        = "my child product category"
  handle             = "my-child-handle"
  is_internal        = false
  is_active          = false
  parent_category_id = medusa_product_category.my-parent-product-category.id
}
