package controller

import (
	"net/http"
	pb "ungraded_5/internal/product"

	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductController struct {
	Client pb.ProductServiceClient
}

func NewProductController(client pb.ProductServiceClient) ProductController {
	return ProductController{
		Client: client,
	}
}

func (p ProductController) Create(c echo.Context) error {
	var requestData pb.AddProductRequest
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	res, err := p.Client.AddProduct(c.Request().Context(), &requestData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Product Created",
		"data": res,
	})
}

func (p ProductController) GetAll(c echo.Context) error {
	res, err := p.Client.GetProducts(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Products",
		"data": res.Products,
	})
}

func (p ProductController) Update(c echo.Context) error {
	var requestData pb.UpdateProductRequest
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	requestData.Id = c.Param("id")

	_, err := p.Client.UpdateProduct(c.Request().Context(), &requestData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Product updated",
		"data": &requestData,
	})
}

func (p ProductController) Delete(c echo.Context) error {
	var requestData pb.DeleteProductRequest
	requestData.Id = c.Param("id")

	res, err := p.Client.DeleteProduct(c.Request().Context(), &requestData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Product deleted",
		"data": res,
	})
}