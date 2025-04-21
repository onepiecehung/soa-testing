package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"product-management/internal/dto"
	"product-management/internal/models"
	"product-management/internal/repositories"
	"product-management/internal/services"
	"product-management/internal/types"
	"product-management/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productRepo    *repositories.ProductRepository
	productService *services.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productRepo *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{
		productRepo:    productRepo,
		productService: services.NewProductService(),
	}
}

// ListProducts godoc
// @Summary      List products
// @Description  Get a paginated list of products with optional filters
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page       query     int     false  "Page number"
// @Param        page_size      query     int     false  "Items per page"
// @Param        categoryId query     int     false  "Filter by category ID"
// @Param        search     query     string  false  "Search term"
// @Param        sort       query     string  false  "Sort field (name, price, created_at)"
// @Param        statuses   query     []string false "Filter by statuses"
// @Success      200        {object}  types.ProductListResponse
// @Failure      400        {object}  types.ErrorResponse
// @Failure      500        {object}  types.ErrorResponse
// @Router       /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var req dto.ProductSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	products, total, err := h.productService.ListProducts(
		req.Page,
		req.PageSize,
		req.CategoryID,
		req.Search,
		req.Sort,
		req.Statuses,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.NewProductListResponse(products, total, req.Page, req.PageSize))
}

// GetProduct godoc
// @Summary      Get a product
// @Description  Get a product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  types.APIResponse
// @Failure      400  {object}  types.ErrorResponse
// @Failure      404  {object}  types.ErrorResponse
// @Failure      500  {object}  types.ErrorResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    product,
	})
}

// CreateProduct godoc
// @Summary      Create a product
// @Description  Create a new product with categories
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        product  body      dto.CreateProductRequest  true  "Product details"
// @Success      201      {object}  types.APIResponse
// @Failure      400      {object}  types.ErrorResponse
// @Failure      500      {object}  types.ErrorResponse
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	// Validate and get categories
	categories, err := h.validateCategories(req.Categories)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	// Create product
	product := &models.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.Quantity,
		Status:        models.StatusActive,
	}

	if err := h.productService.CreateProduct(product, categories); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, types.APIResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update an existing product
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id       path      int                     true  "Product ID"
// @Param        product  body      dto.UpdateProductRequest true  "Product details to update"
// @Success      200      {object}  types.APIResponse
// @Failure      400      {object}  types.ErrorResponse
// @Failure      500      {object}  types.ErrorResponse
// @Router       /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	// Update product
	product := &models.Product{
		BaseModel:     models.BaseModel{ID: uint(id)},
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.Quantity,
		Status:        models.ProductStatus(req.Status),
	}

	if err := h.productService.UpdateProduct(product, req.Categories); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  types.SuccessResponse
// @Failure      400  {object}  types.ErrorResponse
// @Failure      500  {object}  types.ErrorResponse
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{Message: "Product deleted successfully"})
}

// GetWishlist godoc
// @Summary      Get wishlist
// @Description  Get the user's wishlist
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page  query     int  false  "Page number"
// @Param        limit query     int  false  "Items per page"
// @Success      200   {object}  types.WishlistResponse
// @Failure      500   {object}  types.ErrorResponse
// @Router       /products/wishlist [get]
func (h *ProductHandler) GetWishlist(c *gin.Context) {
	pagination := utils.ParsePaginationParams(
		c.DefaultQuery("page", "1"),
		c.DefaultQuery("limit", "10"),
	)

	currentUserID := c.GetUint("userID")
	wishlist, total, err := h.productService.GetWishlist(currentUserID, pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.NewWishlistResponse(wishlist, total, pagination.Page, pagination.Limit))
}

// AddToWishlist godoc
// @Summary      Add to wishlist
// @Description  Add a product to the user's wishlist if it's not already added
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        product_id path      int  true  "Product ID"
// @Success      200        {object}  types.APIResponse
// @Failure      400        {object}  types.ErrorResponse
// @Failure      500        {object}  types.ErrorResponse
// @Router       /products/wishlist/{product_id} [post]
func (h *ProductHandler) AddToWishlist(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	currentUserID := c.GetUint("userID")

	// Check if product is already in wishlist
	isInWishlist, err := h.productService.IsProductInWishlist(currentUserID, uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	if isInWishlist {
		c.JSON(http.StatusOK, types.APIResponse{
			Success: true,
			Message: "Product is already in wishlist",
		})
		return
	}

	// Add to wishlist if not already added
	if err := h.productService.AddToWishlist(currentUserID, uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Product added to wishlist successfully",
	})
}

// RemoveFromWishlist godoc
// @Summary      Remove from wishlist
// @Description  Remove a product from the user's wishlist
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        product_id path      int  true  "Product ID"
// @Success      200        {object}  types.SuccessResponse
// @Failure      400        {object}  types.ErrorResponse
// @Failure      500        {object}  types.ErrorResponse
// @Router       /products/wishlist/{product_id} [delete]
func (h *ProductHandler) RemoveFromWishlist(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	currentUserID := c.GetUint("userID")
	if err := h.productService.RemoveFromWishlist(currentUserID, uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{Message: "Product removed from wishlist"})
}

// GetTotalWishlistCount godoc
// @Summary      Get total wishlist count
// @Description  Get the total number of wishlist items
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.APIResponse
// @Failure      500  {object}  types.ErrorResponse
// @Security     Bearer
// @Router       /products/wishlist/count [get]
func (h *ProductHandler) GetTotalWishlistCount(c *gin.Context) {
	count, err := h.productRepo.CountTotalWishlistItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to count wishlist items",
		})
		return
	}
	userID, exists := c.Get("userID")
	myWishlistCount, err := h.productRepo.CountUserWishlistItems(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to count user wishlist items",
		})
		return
	}

	if !exists {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Total wishlist count retrieved successfully",
		Data: gin.H{
			"total_wishlist_count": count,
			"my_wishlist_count":    myWishlistCount,
		},
	})
}

// validateCategories checks for duplicate category IDs and validates their existence
func (h *ProductHandler) validateCategories(categoryIDs []uint) ([]models.Category, error) {
	categoryMap := make(map[uint]bool)
	for _, id := range categoryIDs {
		if categoryMap[id] {
			return nil, fmt.Errorf("duplicate category ID found: %d", id)
		}
		categoryMap[id] = true
	}

	var categories []models.Category
	for _, categoryID := range categoryIDs {
		var category models.Category
		if err := h.productRepo.DB().First(&category, categoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("category not found with ID: %d", categoryID)
			}
			return nil, fmt.Errorf("failed to fetch category: %v", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}
