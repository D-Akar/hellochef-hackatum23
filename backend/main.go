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
	tags         []string
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
	utils.AddTagsToRecipes(recipes)
	recipesShort = utils.ConvertRecipesToShortRecipes(recipes)
	tags = utils.GetTags(recipes)

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/recipes", getRecipes)
	router.GET("/tags", getTags)
	router.GET("/recipe/:id", getRecipe)
	router.GET("/image/:path", getImage)
	log.Fatal(router.Run(":8080"))

	return
}

func getRecipes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, recipesShort)
}

func getImage(ctx *gin.Context) {
	name := ctx.Param("path")
	if name[len(name)-4:] == ".png" {
		ctx.Header("Content-Type", "image/png")
	} else if name[len(name)-4:] == ".jpg" || name[len(name)-5:] == ".jpeg" {
		ctx.Header("Content-Type", "image/jpeg")
	}
	ctx.File(fmt.Sprintf("database/images/%s", name))
}

func getTags(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, tags)
}

func getRecipe(ctx *gin.Context) {
	id := ctx.GetInt("id")
	ctx.JSON(http.StatusOK, recipes[id])
}

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}
