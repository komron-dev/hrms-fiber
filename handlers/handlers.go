package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/komron-dev/hrms-fiber/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func handleError(c *fiber.Ctx, statusCode int, err error) error {
	log.Error(err)
	return c.Status(statusCode).SendString(err.Error())
}

func CreateEmployee(c *fiber.Ctx) error {
	collection := models.MG.Db.Collection("employees")
	employee := new(models.Employee)

	if err := c.BodyParser(employee); err != nil {
		return handleError(c, 400, err)
	}
	employee.ID = ""
	insertionRes, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return handleError(c, 500, err)
	}

	filter := bson.D{{Key: "_id", Value: insertionRes.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdEmployee := new(models.Employee)
	err = createdRecord.Decode(&createdEmployee)
	if err != nil {
		return handleError(c, 400, err)
	}

	return c.Status(201).JSON(createdEmployee)
}

func ListEmployees(c *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := models.MG.Db.Collection("employees").Find(c.Context(), &query)
	if err != nil {
		return handleError(c, 500, err)
	}

	var employees = make([]models.Employee, 0)

	if err := cursor.All(c.Context(), &employees); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(employees)
}

func GetEmployeeById(c *fiber.Ctx) error {
	id := c.Params("id")

	employeeID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.SendStatus(400)
	}
	query := bson.D{
		{
			Key:   "_id",
			Value: employeeID,
		},
	}
	employee := new(models.Employee)
	if err := models.MG.Db.Collection("employees").FindOne(c.Context(), query).Decode(&employee); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(404).SendString("Employee not found")
		}
		return handleError(c, 500, err)
	}

	return c.JSON(employee)
}

func UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return handleError(c, 400, err)
	}

	employee := new(models.Employee)

	if err := c.BodyParser(employee); err != nil {
		return handleError(c, 400, err)
	}

	query := bson.D{{Key: "_id", Value: employeeID}}
	updatedEmployee := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "age", Value: employee.Age},
				{Key: "salary", Value: employee.Salary},
			},
		},
	}
	err = models.MG.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, updatedEmployee).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return handleError(c, 404, err)
		}
		return handleError(c, 500, err)
	}

	employee.ID = id

	return c.Status(200).JSON(employee)
}

func DeleteEmployee(c *fiber.Ctx) error {
	employeeID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return handleError(c, 400, err)
	}

	query := bson.D{{Key: "_id", Value: employeeID}}
	result, err := models.MG.Db.Collection("employees").DeleteOne(c.Context(), &query)
	if err != nil {
		return handleError(c, 500, err)
	}
	if result.DeletedCount < 1 {
		return handleError(c, 404, err)
	}

	return c.Status(200).JSON("record deleted")
}
