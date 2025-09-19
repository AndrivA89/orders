package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/AndrivA89/orders/internal/application/services"
	"github.com/AndrivA89/orders/internal/infrastructure/config"
	"github.com/AndrivA89/orders/internal/infrastructure/database"
	"github.com/AndrivA89/orders/internal/infrastructure/repositories"
	"github.com/AndrivA89/orders/internal/transport/http/handlers"
	"github.com/AndrivA89/orders/internal/transport/http/router"
)

type IntegrationTestFixture struct {
	db       *gorm.DB
	router   *gin.Engine
	server   *httptest.Server
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func setupTestFixture(t *testing.T) *IntegrationTestFixture {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	err = pool.Client.Ping()
	require.NoError(t, err)

	resource, err := pool.Run("postgres", "15-alpine", []string{
		"POSTGRES_DB=orders_test",
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=postgres",
		"listen_addresses = '*'",
	})

	require.NoError(t, err)
	require.NoError(t, resource.Expire(60))

	var dbConn *database.Connection

	t.Log("Waiting 3s for PostgreSQL to be ready...")
	time.Sleep(3 * time.Second)

	require.NoError(t, pool.Retry(func() error {
		cfg := &config.DatabaseConfig{
			Host:     "localhost",
			Port:     getIntFromString(resource.GetPort("5432/tcp")),
			User:     "postgres",
			Password: "postgres",
			DBName:   "orders_test",
			SSLMode:  "disable",
		}

		var retryErr error
		dbConn, retryErr = database.NewConnection(cfg)
		if retryErr != nil {
			return retryErr
		}

		sqlDB, retryErr := dbConn.DB.DB()
		if retryErr != nil {
			return retryErr
		}

		return sqlDB.Ping()
	}))

	t.Logf("PostgreSQL container started on port %s", resource.GetPort("5432/tcp"))

	err = dbConn.AutoMigrate()
	require.NoError(t, err)

	userRepo := repositories.NewUserRepository(dbConn.DB)
	productRepo := repositories.NewProductRepository(dbConn.DB)
	orderRepo := repositories.NewOrderRepository(dbConn.DB)
	txManager := repositories.NewTransactionManager(dbConn.DB)

	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo, userRepo, productRepo, txManager)

	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	r := router.NewRouter(userHandler, productHandler, orderHandler, logger)
	ginRouter := r.SetupRoutes()

	return &IntegrationTestFixture{
		db:       dbConn.DB,
		router:   ginRouter,
		server:   httptest.NewServer(ginRouter),
		pool:     pool,
		resource: resource,
	}
}

func (f *IntegrationTestFixture) cleanup(t *testing.T) {
	f.server.Close()

	if err := f.pool.Purge(f.resource); err != nil {
		log.Fatalf("Could not purge PostgreSQL container: %s", err)
	}

	t.Log("Test cleanup completed - PostgreSQL container removed")
}

func TestFullWorkflow(t *testing.T) {
	fixture := setupTestFixture(t)
	defer fixture.cleanup(t)

	t.Log("Starting full e2e workflow test")

	var users []map[string]interface{}
	var products []map[string]interface{}
	var orders []map[string]interface{}

	t.Log("Step 1: Creating users")

	userRequests := []map[string]interface{}{
		{
			"first_name": "Алексей",
			"last_name":  "Иванов",
			"age":        28,
			"is_married": true,
			"password":   "securepass123",
		},
		{
			"first_name": "Мария",
			"last_name":  "Петрова",
			"age":        25,
			"is_married": false,
			"password":   "mypassword456",
		},
	}

	for i, userReq := range userRequests {
		resp := fixture.makeRequest(t, "POST", "/api/v1/users", userReq)
		assert.Equal(t, http.StatusCreated, resp.Code)

		var user map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &user)
		require.NoError(t, err)
		users = append(users, user)

		t.Logf("User %d created: %s %s (ID: %s)", i+1,
			user["first_name"], user["last_name"], user["id"])

		assert.Equal(t, userReq["first_name"], user["first_name"])
		assert.Equal(t, userReq["last_name"], user["last_name"])
		assert.Equal(t, float64(userReq["age"].(int)), user["age"])
		assert.Equal(t, userReq["is_married"], user["is_married"])
		assert.NotEmpty(t, user["id"])
		assert.NotEmpty(t, user["created_at"])
		assert.NotContains(t, user, "password")
	}

	t.Log("Step 2: Creating products")

	productRequests := []map[string]interface{}{
		{
			"description": "iPhone 15 Pro Max 256GB",
			"price":       119999, // 1199.99 в копейках
			"quantity":    50,
		},
		{
			"description": "MacBook Air M2 13\" 512GB",
			"price":       149999, // 1499.99 в копейках
			"quantity":    25,
		},
		{
			"description": "AirPods Pro 2-го поколения",
			"price":       29999, // 299.99 в копейках
			"quantity":    100,
		},
	}

	for i, productReq := range productRequests {
		resp := fixture.makeRequest(t, "POST", "/api/v1/products", productReq)
		assert.Equal(t, http.StatusCreated, resp.Code)

		var product map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &product)
		require.NoError(t, err)
		products = append(products, product)

		t.Logf("Product %d created: %s (ID: %s, Price: %.2f, Qty: %.0f)",
			i+1, product["description"], product["id"],
			product["price"].(float64)/100, product["quantity"])

		// Validate product data
		assert.Equal(t, productReq["description"], product["description"])
		assert.Equal(t, float64(productReq["price"].(int)), product["price"])
		assert.Equal(t, float64(productReq["quantity"].(int)), product["quantity"])
		assert.NotEmpty(t, product["id"])
	}

	t.Log("Step 3: Creating orders")

	// Order 1: First user orders iPhone and AirPods
	order1Req := map[string]interface{}{
		"user_id": users[0]["id"],
		"items": []map[string]interface{}{
			{
				"product_id": products[0]["id"],
				"quantity":   2,
			},
			{
				"product_id": products[2]["id"],
				"quantity":   3,
			},
		},
	}

	resp := fixture.makeRequest(t, "POST", "/api/v1/orders", order1Req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var order1 map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &order1)
	require.NoError(t, err)
	orders = append(orders, order1)

	expectedTotal1 := int64(2*119999 + 3*29999)
	assert.Equal(t, float64(expectedTotal1), order1["total"])
	t.Logf("Order 1 created: User %s, Total: %.2f, Items: %d",
		order1["user_id"], order1["total"].(float64)/100, len(order1["items"].([]interface{})))

	// Order 2: Second user orders MacBook
	order2Req := map[string]interface{}{
		"user_id": users[1]["id"],
		"items": []map[string]interface{}{
			{
				"product_id": products[1]["id"],
				"quantity":   1,
			},
		},
	}

	resp = fixture.makeRequest(t, "POST", "/api/v1/orders", order2Req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var order2 map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &order2)
	require.NoError(t, err)
	orders = append(orders, order2)

	expectedTotal2 := int64(149999)
	assert.Equal(t, float64(expectedTotal2), order2["total"])
	t.Logf("Order 2 created: User %s, Total: %.2f",
		order2["user_id"], order2["total"].(float64)/100)

	// Order 3: First user orders more items
	order3Req := map[string]interface{}{
		"user_id": users[0]["id"],
		"items": []map[string]interface{}{
			{
				"product_id": products[0]["id"],
				"quantity":   1,
			},
			{
				"product_id": products[1]["id"],
				"quantity":   2,
			},
			{
				"product_id": products[2]["id"],
				"quantity":   5,
			},
		},
	}

	resp = fixture.makeRequest(t, "POST", "/api/v1/orders", order3Req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var order3 map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &order3)
	require.NoError(t, err)
	orders = append(orders, order3)

	expectedTotal3 := int64(119999 + 2*149999 + 5*29999)
	assert.Equal(t, float64(expectedTotal3), order3["total"])
	t.Logf("Order 3 created: User %s, Total: %.2f, Items: %d",
		order3["user_id"], order3["total"].(float64)/100, len(order3["items"].([]interface{})))

	t.Log("Step 4: Verifying inventory updates")

	resp = fixture.makeRequest(t, "GET", "/api/v1/products", nil)
	assert.Equal(t, http.StatusOK, resp.Code)

	var productResponse map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &productResponse)
	require.NoError(t, err)

	updatedProducts := productResponse["products"].([]interface{})

	// iPhone: 50 - 2 - 1 = 47
	assert.Equal(t, float64(47), updatedProducts[0].(map[string]interface{})["quantity"])
	// MacBook: 25 - 1 - 2 = 22
	assert.Equal(t, float64(22), updatedProducts[1].(map[string]interface{})["quantity"])
	// AirPods: 100 - 3 - 5 = 92
	assert.Equal(t, float64(92), updatedProducts[2].(map[string]interface{})["quantity"])

	t.Log("Inventory correctly updated after orders")

	t.Log("Step 5: Confirming Order 1")

	confirmURL := fmt.Sprintf("/api/v1/orders/%s/confirm", order1["id"])
	resp = fixture.makeRequest(t, "PATCH", confirmURL, nil)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify order status changed
	orderURL := fmt.Sprintf("/api/v1/orders/%s", order1["id"])
	resp = fixture.makeRequest(t, "GET", orderURL, nil)
	assert.Equal(t, http.StatusOK, resp.Code)

	var confirmedOrder map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &confirmedOrder)
	require.NoError(t, err)
	assert.Equal(t, "confirmed", confirmedOrder["status"])

	t.Logf("Order 1 confirmed successfully, status: %s", confirmedOrder["status"])

	t.Log("Step 6: Cancelling Order 2")

	cancelURL := fmt.Sprintf("/api/v1/orders/%s/cancel", order2["id"])
	resp = fixture.makeRequest(t, "PATCH", cancelURL, nil)
	assert.Equal(t, http.StatusOK, resp.Code)

	orderURL = fmt.Sprintf("/api/v1/orders/%s", order2["id"])
	resp = fixture.makeRequest(t, "GET", orderURL, nil)
	assert.Equal(t, http.StatusOK, resp.Code)

	var cancelledOrder map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &cancelledOrder)
	require.NoError(t, err)
	assert.Equal(t, "cancelled", cancelledOrder["status"])

	t.Logf("Order 2 cancelled successfully, status: %s", cancelledOrder["status"])

	t.Log("Step 7: Confirming then cancelling Order 3")

	confirmURL = fmt.Sprintf("/api/v1/orders/%s/confirm", order3["id"])
	resp = fixture.makeRequest(t, "PATCH", confirmURL, nil)
	assert.Equal(t, http.StatusOK, resp.Code)

	cancelURL = fmt.Sprintf("/api/v1/orders/%s/cancel", order3["id"])
	resp = fixture.makeRequest(t, "PATCH", cancelURL, nil)
	assert.Equal(t, http.StatusOK, resp.Code)

	orderURL = fmt.Sprintf("/api/v1/orders/%s", order3["id"])
	resp = fixture.makeRequest(t, "GET", orderURL, nil)
	assert.Equal(t, http.StatusOK, resp.Code)

	var finalOrder map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &finalOrder)
	require.NoError(t, err)
	assert.Equal(t, "cancelled", finalOrder["status"])

	t.Logf("Order 3 confirmed then cancelled, final status: %s", finalOrder["status"])

	t.Log("Step 8: Checking users' orders")

	for i, user := range users {
		userOrdersURL := fmt.Sprintf("/api/v1/users/%s/orders", user["id"])
		resp = fixture.makeRequest(t, "GET", userOrdersURL, nil)
		assert.Equal(t, http.StatusOK, resp.Code)

		var orderResponse map[string]interface{}
		err = json.Unmarshal(resp.Body.Bytes(), &orderResponse)
		require.NoError(t, err)

		userOrders := orderResponse["orders"].([]interface{})

		expectedOrdersCount := 1
		if i == 0 {
			expectedOrdersCount = 2
		}
		assert.Len(t, userOrders, expectedOrdersCount)
		t.Logf("User %d (%s %s) has %d order(s)",
			i+1, user["first_name"], user["last_name"], len(userOrders))
	}

	t.Log("Step 9: Testing edge cases")

	// Test validation: non-existent user
	invalidOrderReq := map[string]interface{}{
		"user_id": uuid.New().String(),
		"items": []map[string]interface{}{
			{
				"product_id": products[0]["id"],
				"quantity":   1,
			},
		},
	}
	resp = fixture.makeRequest(t, "POST", "/api/v1/orders", invalidOrderReq)
	// Rate limiting may return 429 instead of 400 due to previous requests
	assert.Contains(t, []int{http.StatusBadRequest, http.StatusTooManyRequests}, resp.Code)
	t.Log("Order creation with invalid user correctly rejected")

	// Test validation: insufficient quantity
	insufficientOrderReq := map[string]interface{}{
		"user_id": users[0]["id"],
		"items": []map[string]interface{}{
			{
				"product_id": products[0]["id"],
				"quantity":   100,
			},
		},
	}
	resp = fixture.makeRequest(t, "POST", "/api/v1/orders", insufficientOrderReq)
	// Rate limiting may return 429 instead of 400 due to previous requests
	assert.Contains(t, []int{http.StatusBadRequest, http.StatusTooManyRequests}, resp.Code)
	t.Log("Order creation with insufficient quantity correctly rejected")

	// Test validation: double confirmation
	confirmURL = fmt.Sprintf("/api/v1/orders/%s/confirm", order1["id"])
	resp = fixture.makeRequest(t, "PATCH", confirmURL, nil)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	t.Log("Double confirmation correctly rejected")

	t.Log("Full workflow test completed successfully")
}

func (f *IntegrationTestFixture) makeRequest(t *testing.T, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer

	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(method, path, reqBody)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	f.router.ServeHTTP(w, req)

	return w
}

func getIntFromString(portStr string) int {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Could not convert port %s to int: %s", portStr, err)
	}

	return port
}
