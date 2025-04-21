package handlers

import (
	"net/http"
	"strconv"

	"product-management/internal/dto"
	"product-management/internal/services"
	"product-management/internal/types"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a new category with name and optional description
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request  body      dto.CreateCategoryRequest  true  "Category details"
// @Success      201     {object}   types.APIResponse
// @Failure      400     {object}   types.ErrorResponse
// @Failure      500     {object}   types.ErrorResponse
// @Router       /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	c.JSON(http.StatusCreated, types.APIResponse{
		Success: true,
		Message: "Category created successfully",
		Data:    response,
	})
}

// GetCategoryByID godoc
// @Summary      Get a category
// @Description  Get a category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  types.APIResponse
// @Failure      400  {object}  types.ErrorResponse
// @Failure      404  {object}  types.ErrorResponse
// @Security     Bearer
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    category,
	})
}

// GetAllCategories godoc
// @Summary      List categories
// @Description  Get all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.APIResponse
// @Failure      500  {object}  types.ErrorResponse
// @Security     Bearer
// @Router       /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    categories,
	})
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Update an existing category with name and optional description
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id       path      int                        true  "Category ID"
// @Param        request  body      dto.UpdateCategoryRequest  true  "Category details"
// @Success      200     {object}   types.APIResponse
// @Failure      400     {object}   types.ErrorResponse
// @Failure      404     {object}   types.ErrorResponse
// @Failure      500     {object}   types.ErrorResponse
// @Router       /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.categoryService.UpdateCategory(uint(id), req)
	if err != nil {
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Category updated successfully",
		Data:    response,
	})
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Delete a category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      204  {object}  types.SuccessResponse
// @Failure      400  {object}  types.ErrorResponse
// @Failure      500  {object}  types.ErrorResponse
// @Security     Bearer
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		if err.Error() == "cannot delete category with associated products" {
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{Message: "Category deleted successfully"})
}

// GetProductsByCategoryID godoc
// @Summary      Get category products
// @Description  Get all products in a specific category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  types.APIResponse
// @Failure      400  {object}  types.ErrorResponse
// @Failure      500  {object}  types.ErrorResponse
// @Security     Bearer
// @Router       /categories/{id}/products [get]
func (h *CategoryHandler) GetProductsByCategoryID(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	products, err := h.categoryService.GetProductsByCategoryID(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    products,
	})
}

// AddProductToCategory godoc
// @Summary      Add product to category
// @Description  Add a product to a specific category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id         path      int  true  "Category ID"
// @Param        productId  path      int  true  "Product ID"
// @Success      204        {object}  types.SuccessResponse
// @Failure      400        {object}  types.ErrorResponse
// @Failure      500        {object}  types.ErrorResponse
// @Security     Bearer
// @Router       /categories/{id}/products/{productId} [post]
func (h *CategoryHandler) AddProductToCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	if err := h.categoryService.AddProductToCategory(uint(categoryID), uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{Message: "Product added to category successfully"})
}

// RemoveProductFromCategory godoc
// @Summary      Remove product from category
// @Description  Remove a product from a specific category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id         path      int  true  "Category ID"
// @Param        productId  path      int  true  "Product ID"
// @Success      204        {object}  types.SuccessResponse
// @Failure      400        {object}  types.ErrorResponse
// @Failure      500        {object}  types.ErrorResponse
// @Security     Bearer
// @Router       /categories/{id}/products/{productId} [delete]
func (h *CategoryHandler) RemoveProductFromCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	if err := h.categoryService.RemoveProductFromCategory(uint(categoryID), uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{Message: "Product removed from category successfully"})
}

// GetCategoryDistribution godoc
// @Summary      Get category distribution
// @Description  Get the distribution of products across categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.APIResponse
// @Failure      500  {object}  types.ErrorResponse
// @Security     Bearer
// @Router       /categories/distribution [get]
func (h *CategoryHandler) GetCategoryDistribution(c *gin.Context) {
	distributions, err := h.categoryService.GetCategoryDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get category distribution: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    distributions,
	})
}
