package food

import (
    "testing"
    "strings"
    "fmt"
    "io"
)

func TestFoodType(t *testing.T) {
    o := &Food{Calories: 4000}
    if o.Calories != 4000 {
        t.FailNow()
    }
}

func TestCreateFoodCollection(t *testing.T) {
    type exampleWithExpectations struct {
        calorieListReader io.Reader
        numberOfFoodItems int
        totalCalorieCount int
    }

    examplesAndExpectations := []exampleWithExpectations{
        // Coded string of some calories
        {
            calorieListReader: strings.NewReader(
`4400
1100
2200
2200`),
            numberOfFoodItems: 4,
            totalCalorieCount: 9900,
        },
        {
            calorieListReader: strings.NewReader(fmt.Sprint("8000\n1000\n500\n")),
            numberOfFoodItems: 3,
            totalCalorieCount: 9500,
        },
    }

    for _, example := range examplesAndExpectations {
        food := New(example.calorieListReader)
        if len(*food) != example.numberOfFoodItems {
            t.Fail()
        }
        if food.TotalCalories() != example.totalCalorieCount {
            t.Fail()
        }
    }
}
