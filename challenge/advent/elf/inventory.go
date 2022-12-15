package elf

import (
    "github.com/GPKyte/banter/challenge/advent/food"
    "github.com/GPKyte/banter/challenge/advent/sack"
)

// Inventory is used to store, retrieve, and track or summarize information about collections of objects represented as data, e.g. food.
type Inventory struct {
    Foods *food.FoodCollection
    Sack *sack.Sack
}

// AddFood to Inventory's food category
// Use reference to food object and the inventory in case food is eaten, altered, or if inventory is transferred, rather than copy value of food.
func (i *Inventory) AddFood(f *food.Food) {
    i.Foods.Include(f)
}

func (i *Inventory) TotalCalories() int {
    var runningCalorieCount int = 0
    // Alt implementation: put this method over *food.FoodCollection

    for _, snack := range *i.Foods {
        runningCalorieCount += snack.Calories
    }

    return runningCalorieCount
}
