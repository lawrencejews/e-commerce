package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lawrencejews/e-commerce/models"
	"go.mongodb.org/mongo-driver/bson"
)

// HashPassword
func HashPassword(password string) string{

}

// VerifyPassword
func VerifyPassword(userPassword string, givenPassword string)(bool, string){

}

// SignUp
func SignUp() gin.HandlerFunc {

	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// User
		var user models.User 
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"err.Error()"})
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone number is already in use"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
	}
}

// SignIn
func Login() gin.HandlerFunc{

}

// ProductViewerAdmin
func ProductViewerAdmin() gin.HandlerFunc{

}

// SearchProduct
func SearchProduct() gin.HandlerFunc{

}

// SearchProductByQuery
func SearchProductByQuery() gin.HandlerFunc{
	
}