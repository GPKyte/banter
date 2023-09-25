// This package helps create a collection with convenience methods to describe food and its calories.
package food

import (
    "fmt"
    "text/scanner"
    "strconv"
    "io"
)

type FoodCollection []*Food

const (
    Debug bool = false
)

func New(source io.Reader) *FoodCollection {
    foods := make(FoodCollection, 0)
    var s scanner.Scanner
    s.Init(source)
    for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
        cal, err := strconv.Atoi(s.TokenText())
        if err != nil {
            examineInput(s.TokenText())
        } else {
            foods = append(foods, &Food{Calories: cal})
        }
    }

    return &foods
}

func (fc *FoodCollection) TotalCalories() int {
    var total int = 0
    for _, food := range *fc {
        total += food.Calories
    }

    return total
}

func (fc *FoodCollection) Include(f *Food) {
    *fc = append(*fc, f)
}

func examineInput(text string) {
    if Debug {
        fmt.Println(text)
    }
}

type Food struct {
    Calories int
}

