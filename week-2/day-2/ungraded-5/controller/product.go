package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	pb "ungraded_5/internal/product"
	"ungraded_5/models"
	"ungraded_5/service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductController struct {
	Client pb.ProductServiceClient
	Redis *redis.Client
	Mb service.MessageBroker
}

func NewProductController(client pb.ProductServiceClient, redis *redis.Client, mb service.MessageBroker) ProductController {
	return ProductController{
		Client: client,
		Redis: redis,
		Mb: mb,
	}
}

func (p ProductController) Create(c echo.Context) error {
	var requestData pb.AddProductRequest
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	ctx := grpcMetadata.AppendToOutgoingContext(c.Request().Context(), "authorization", "Bearer "+ os.Getenv("GRPC_AUTH"))
	res, err := p.Client.AddProduct(ctx, &requestData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if err := p.Mb.PublishMessage("User created product"); err != nil {
		return err
	}
	
	if err := p.TriggerCache(); err != nil {
		return err
	}
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Product Created",
		"data": res,
	})
}

func (p ProductController) GetAll(c echo.Context) error {
	productsStr, err := p.Redis.Get(context.TODO(), "products").Result()
	if err == redis.Nil {

		ctx := grpcMetadata.AppendToOutgoingContext(c.Request().Context(), "authorization", "Bearer "+ os.Getenv("GRPC_AUTH"))
		res, err := p.Client.GetProducts(ctx, &emptypb.Empty{})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		productDatas, err := json.Marshal(res.Products)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		p.Redis.Set(context.TODO(), "products", productDatas, 0)
		productsStr = string(productDatas)

	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if err := p.Mb.PublishMessage("User getting all products"); err != nil {
		return err
	}
	
	var productsDatas []models.Product
	if err := json.Unmarshal([]byte(productsStr), &productsDatas); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Products",
		"data": productsDatas,
	})
}

func (p ProductController) Update(c echo.Context) error {
	var requestData pb.UpdateProductRequest
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	requestData.Id = c.Param("id")

	ctx := grpcMetadata.AppendToOutgoingContext(c.Request().Context(), "authorization", "Bearer "+ os.Getenv("GRPC_AUTH"))
	_, err := p.Client.UpdateProduct(ctx, &requestData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if err := p.Mb.PublishMessage("User updated a product"); err != nil {
		return err
	}
	
	if err := p.TriggerCache(); err != nil {
		return err
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Product updated",
		"data": &requestData,
	})
}

func (p ProductController) Delete(c echo.Context) error {
	var requestData pb.DeleteProductRequest
	requestData.Id = c.Param("id")

	ctx := grpcMetadata.AppendToOutgoingContext(c.Request().Context(), "authorization", "Bearer "+ os.Getenv("GRPC_AUTH"))
	res, err := p.Client.DeleteProduct(ctx, &requestData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if err := p.Mb.PublishMessage("User deleted a product"); err != nil {
		return err
	}

	if err := p.TriggerCache(); err != nil {
		return err
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Product deleted",
		"data": res,
	})
}

func (p ProductController) TriggerCache() error {
	ctx := grpcMetadata.AppendToOutgoingContext(context.TODO(), "authorization", "Bearer "+ os.Getenv("GRPC_AUTH"))
	res, err := p.Client.GetProducts(ctx, &emptypb.Empty{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	productDatas, err := json.Marshal(res.Products)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	p.Redis.Set(context.TODO(), "products", productDatas, 0)
	return nil
}