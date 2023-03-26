package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lawrencejews/e-commerce/database"
	"github.com/lawrencejews/e-commerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongodb data collection types
type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

// Handles creation of collections
func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

// Adding an item to cart
func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return

			userQueryID := c.Query("userID")
			if userQueryID == "" {
				log.Println("user is empty")
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("user is empty"))
				return
			}

			productID, err := primitive.ObjectIDFromHex(productQueryID)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
			}
			c.IndentedJSON(200, "Successfully added to cart")
		}
	}
}

// Remove an item from the cart
func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return

			userQueryID := c.Query("userID")
			if userQueryID == "" {
				log.Println("user is empty")
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("user is empty"))
				return
			}

			productID, err := primitive.ObjectIDFromHex(productQueryID)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
			}

			c.IndentedJSON(200, "Successfully removed item from cart")
		}
	}
}

// Getting an item from a cart using MongoDB Aggregation approach
func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid"})
			c.Abort()
			return
		}

		userOne_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Finding the User and create an aggregation function
		var filledcart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userOne_id}}).Decode(&filledcart)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}

		// Filter match section
		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userOne_id}}}}

		// Unwind section to locates cart-items for the user
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}

		// Grouping
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})

		if err != nil {
			log.Println(err)
		}

		var listing []bson.M
		if err = pointCursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledcart.UserCart)
		}
		ctx.Done()
	}
}

// Buying an item from cart
func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")

		if userQueryID == "" {
			log.Panic("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Successfully placed order")

	}
}

// Instantly buying an item from cart
func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return

			userQueryID := c.Query("userID")
			if userQueryID == "" {
				log.Println("user is empty")
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("user is empty"))
				return
			}

			productID, err := primitive.ObjectIDFromHex(productQueryID)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = database.InstantBuy(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
			}

			c.IndentedJSON(200, "Successfully placed the order")
		}
	}
}
