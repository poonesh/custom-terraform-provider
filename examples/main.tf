terraform {
  required_providers {
    food = {
      version = "0.1"
      source  = "test/default/food"
    }
  }
}

locals {
  foods = [
    {
      name   = "Bread"
      origin = "France"
    },
    {
      name   = "Poutine"
      origin = "Canada"
    },
    {
      name   = "beer"
      origin = "Germany" 
    }
  ]
}

// Food Resource
resource "food" "sample" {
  count  = length(local.foods)
  name   = local.foods[count.index].name
  origin = local.foods[count.index].origin
}

