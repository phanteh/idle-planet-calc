package main

func calculateIngredients(order []Ingredient) (bill map[string]Ingredient) {
	bill = make(map[string]Ingredient)
	for _, o := range order {
		ingredients := getIngredients(o.Item)
		for _, i := range ingredients {
			value := i.Item.Value * i.Amount * o.Amount
			amount := i.Amount * o.Amount
			if item, found := bill[i.Item.Name]; !found {
				bill[i.Item.Name] = Ingredient{i.Item, value, amount}
			} else {
				item.Value += value
				item.Amount += amount
			}
		}
	}
	return
}

func getIngredients(item GameItem) (ingredients map[string]Ingredient) {
	ingredients = make(map[string]Ingredient, 0)
	for _, i := range item.Ingredients {
		if newItem, found := ingredients[i.Item.Name]; !found {
			ingredients[i.Item.Name] = Ingredient{i.Item, i.Value, i.Amount}
		} else {
			newItem.Amount += i.Amount
		}
		if len(i.Item.Ingredients) > 0 {
			subIngredients := getIngredients(i.Item)
			for k, v := range subIngredients {
				if newItem, found := ingredients[k]; !found {
					ingredients[k] = Ingredient{v.Item, v.Value * i.Amount, v.Amount * i.Amount}
				} else {
					newItem.Amount += v.Amount * i.Amount
					ingredients[k] = newItem
				}
			}
		}
	}
	return
}
