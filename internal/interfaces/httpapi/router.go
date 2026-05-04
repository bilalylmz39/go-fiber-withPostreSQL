package httpapi

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	app "ybilaly/internal/application/product"
	domain "ybilaly/internal/domain/product"
)

type productService interface {
	Create(ctx context.Context, input domain.CreateInput) (domain.Product, error)
	Update(ctx context.Context, id int64, input domain.UpdateInput) (domain.Product, error)
	GetByID(ctx context.Context, id int64) (domain.Product, error)
	List(ctx context.Context, filter domain.ListFilter) ([]domain.Product, error)
	Delete(ctx context.Context, id int64) error
	ChangeStock(ctx context.Context, id int64, delta int) (domain.Product, error)
	SetActive(ctx context.Context, id int64, active bool) (domain.Product, error)
}

type Server struct {
	app            *fiber.App
	productService productService
	requestTimeout time.Duration
}

func NewServer(productService productService, requestTimeout time.Duration) *Server {
	server := &Server{
		productService: productService,
		requestTimeout: requestTimeout,
	}

	server.app = fiber.New(fiber.Config{
		AppName:      "Products API",
		ErrorHandler: server.errorHandler,
	})
	server.registerRoutes()

	return server
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) registerRoutes() {
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	s.app.Get("/products", s.listProducts)
	s.app.Post("/products", s.createProduct)
	s.app.Get("/products/:id", s.getProduct)
	s.app.Put("/products/:id", s.updateProduct)
	s.app.Delete("/products/:id", s.deleteProduct)
	s.app.Patch("/products/:id/stock", s.changeStock)
	s.app.Patch("/products/:id/activate", s.activateProduct)
	s.app.Patch("/products/:id/deactivate", s.deactivateProduct)
}

func (s *Server) listProducts(c *fiber.Ctx) error {
	ctx, cancel := s.context()
	defer cancel()

	filter, err := parseListFilter(c)
	if err != nil {
		return err
	}

	products, err := s.productService.List(ctx, filter)
	if err != nil {
		return err
	}

	return c.JSON(products)
}

func (s *Server) getProduct(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	ctx, cancel := s.context()
	defer cancel()

	product, err := s.productService.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(product)
}

func (s *Server) createProduct(c *fiber.Ctx) error {
	var input domain.CreateInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product payload")
	}

	ctx, cancel := s.context()
	defer cancel()

	product, err := s.productService.Create(ctx, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (s *Server) updateProduct(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	var input domain.UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product payload")
	}

	ctx, cancel := s.context()
	defer cancel()

	product, err := s.productService.Update(ctx, id, input)
	if err != nil {
		return err
	}

	return c.JSON(product)
}

func (s *Server) deleteProduct(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	ctx, cancel := s.context()
	defer cancel()

	if err := s.productService.Delete(ctx, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) changeStock(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	var input struct {
		Delta int `json:"delta"`
	}
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid stock payload")
	}

	ctx, cancel := s.context()
	defer cancel()

	product, err := s.productService.ChangeStock(ctx, id, input.Delta)
	if err != nil {
		return err
	}

	return c.JSON(product)
}

func (s *Server) activateProduct(c *fiber.Ctx) error {
	return s.setProductActive(c, true)
}

func (s *Server) deactivateProduct(c *fiber.Ctx) error {
	return s.setProductActive(c, false)
}

func (s *Server) setProductActive(c *fiber.Ctx, active bool) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	ctx, cancel := s.context()
	defer cancel()

	product, err := s.productService.SetActive(ctx, id, active)
	if err != nil {
		return err
	}

	return c.JSON(product)
}

func (s *Server) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.requestTimeout)
}

func (s *Server) errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"

	switch {
	case app.IsValidationError(err):
		code = fiber.StatusBadRequest
		message = err.Error()
	case app.IsNotFound(err):
		code = fiber.StatusNotFound
		message = "product not found"
	default:
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			message = fiberErr.Message
		}
	}

	return c.Status(code).JSON(fiber.Map{"error": message})
}

func parseID(c *fiber.Ctx) (int64, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}
	return id, nil
}

func parseListFilter(c *fiber.Ctx) (domain.ListFilter, error) {
	filter := domain.ListFilter{
		Query:  c.Query("query"),
		Limit:  c.QueryInt("limit", 20),
		Offset: c.QueryInt("offset", 0),
	}

	if value := c.Query("active"); value != "" {
		active, err := strconv.ParseBool(value)
		if err != nil {
			return domain.ListFilter{}, fiber.NewError(fiber.StatusBadRequest, "invalid active filter")
		}
		filter.Active = &active
	}

	if value := c.Query("min_price"); value != "" {
		minPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return domain.ListFilter{}, fiber.NewError(fiber.StatusBadRequest, "invalid min_price filter")
		}
		filter.MinPrice = &minPrice
	}

	if value := c.Query("max_price"); value != "" {
		maxPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return domain.ListFilter{}, fiber.NewError(fiber.StatusBadRequest, "invalid max_price filter")
		}
		filter.MaxPrice = &maxPrice
	}

	return filter, nil
}
