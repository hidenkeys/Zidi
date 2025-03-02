package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s Server) GetAllCustomers(c *fiber.Ctx) error {
	response, err := s.customerService.GetAllCustomers(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Success",
		"Data":    response,
	})
}

func (s Server) CreateCustomer(c *fiber.Ctx) error {
	newCustomer := new(api.Customer)
	if err := c.BodyParser(newCustomer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   err.Error(),
		})
	}

	createdCustomer, err := s.customerService.CreateCustomer(context.Background(), *newCustomer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdCustomer)
}

func (s Server) GetCustomersByOrganization(c *fiber.Ctx, params api.GetCustomersByOrganizationParams) error {
	orgID := params.OrganizationId
	if orgID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid organization ID",
		})
	}

	customers, err := s.customerService.GetCustomersByOrganization(context.Background(), orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(customers)
}

func (s Server) DeleteCustomer(c *fiber.Ctx, id openapi_types.UUID) error {
	if id == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid customer ID",
		})
	}

	err := s.customerService.DeleteCustomer(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			ErrorCode: "404",
			Message:   "Customer not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (s Server) GetCustomerById(c *fiber.Ctx, id openapi_types.UUID) error {
	if id == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid customer ID",
		})
	}

	customer, err := s.customerService.GetCustomerByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			ErrorCode: "404",
			Message:   "Customer not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(customer)
}

func (s Server) UpdateCustomer(c *fiber.Ctx, id openapi_types.UUID) error {
	if id == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid customer ID",
		})
	}

	updateRequest := new(api.Customer)
	if err := c.BodyParser(updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   err.Error(),
		})
	}

	updatedCustomer, err := s.customerService.UpdateCustomer(context.Background(), id, updateRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedCustomer)
}
