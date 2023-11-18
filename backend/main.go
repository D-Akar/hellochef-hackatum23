package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"backend/database"
	"backend/utils"
)

var (
	users        []*database.User
	ingredients  []*database.Ingredient
	recipes      []*database.RecipeWithProperties
	recipesShort []*database.RecipeShort
)

func main() {

	utils.OpenAndUnmarshal("database/users.json", any(&users))
	utils.AssignIdUsers(users)
	utils.OpenAndUnmarshal("database/ingredients.json", any(&ingredients))
	utils.AssignIdIngredients(ingredients)
	var recipesTemp []*database.Recipe
	utils.OpenAndUnmarshal("database/recipes.json", any(&recipesTemp))
	utils.AssignIdRecipes(recipesTemp)
	recipes = utils.RecipesToRecipesWithProperties(recipesTemp, ingredients)
	recipesShort = utils.ConvertRecipesToShortRecipes(recipes)
	fmt.Println(recipesShort[0].Name)

	router := gin.Default()
	router.GET("/recipes", getRecipes)
	log.Fatal(router.Run(":8080"))

	return
}

func getRecipes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, recipesShort)
}
