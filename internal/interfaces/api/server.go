package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/container"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/middleware"
	"gorm.io/gorm"
)

// Server represents the API server
type Server struct {
	config     *config.Config
	router     *mux.Router
	httpServer *http.Server
	logger     logger.Logger
	container  container.Container
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, db *gorm.DB, logger logger.Logger) *Server {
	// Initialize dependency container
	diContainer := container.NewContainer(cfg, db, logger)

	// Initialize default payment providers
	paymentProviderService := diContainer.Services().PaymentProviderService()
	if err := paymentProviderService.InitializeDefaultProviders(); err != nil {
		logger.Error("Failed to initialize default payment providers: %v", err)
	} else {
		logger.Info("Default payment providers initialized successfully")
	}

	router := mux.NewRouter()

	server := &Server{
		config:    cfg,
		router:    router,
		logger:    logger,
		container: diContainer,
	}

	// Apply CORS middleware to all routes
	// corsMiddleware := diContainer.Middlewares().CorsMiddleware()
	// router.Use(corsMiddleware.ApplyCors)

	server.setupRoutes()

	// Create HTTP server
	server.httpServer = &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return server
}

// setupRoutes configures all routes for the API
func (s *Server) setupRoutes() {
	// Extract handlers from container
	userHandler := s.container.Handlers().UserHandler()
	productHandler := s.container.Handlers().ProductHandler()
	categoryHandler := s.container.Handlers().CategoryHandler()
	checkoutHandler := s.container.Handlers().CheckoutHandler()
	orderHandler := s.container.Handlers().OrderHandler()
	paymentHandler := s.container.Handlers().PaymentHandler()
	paymentProviderHandler := s.container.Handlers().PaymentProviderHandler()
	webhookHandlers := s.container.Handlers().WebhookHandlerProvider()
	discountHandler := s.container.Handlers().DiscountHandler()
	shippingHandler := s.container.Handlers().ShippingHandler()
	currencyHandler := s.container.Handlers().CurrencyHandler()
	healthHandler := s.container.Handlers().HealthHandler()
	emailTestHandler := s.container.Handlers().EmailTestHandler()

	// Extract middleware from container
	authMiddleware := s.container.Middlewares().AuthMiddleware()
	// corsMiddleware := s.container.Middlewares().CorsMiddleware()

	// Health check routes (no prefix, for load balancers and monitoring)
	s.router.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	// Register routes
	api := s.router.PathPrefix("/api").Subrouter()
	// api.Use(corsMiddleware.ApplyCors)

	// Webhook routes (separate subrouter without CORS middleware for server-to-server communication)
	webhooks := s.router.PathPrefix("/api/webhooks").Subrouter()

	// Public routes
	api.HandleFunc("/auth/register", userHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/auth/signin", userHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/products/{productId:[0-9]+}", productHandler.GetProduct).Methods(http.MethodGet)

	api.HandleFunc("/products/search", productHandler.SearchProducts).Methods(http.MethodGet)
	api.HandleFunc("/categories", categoryHandler.ListCategories).Methods(http.MethodGet)
	api.HandleFunc("/categories/{id:[0-9]+}", categoryHandler.GetCategory).Methods(http.MethodGet)
	api.HandleFunc("/categories/{id:[0-9]+}/children", categoryHandler.GetChildCategories).Methods(http.MethodGet)
	api.HandleFunc("/payment/providers", paymentHandler.GetAvailablePaymentProviders).Methods(http.MethodGet)

	// Webhook routes (public, no authentication or CORS required for server-to-server communication)
	webhooks.HandleFunc("/stripe", webhookHandlers.StripeHandler().HandleWebhook).Methods(http.MethodPost)
	webhooks.HandleFunc("/mobilepay", webhookHandlers.MobilePayHandler().HandleWebhook).Methods(http.MethodPost)

	// Public discount routes
	api.HandleFunc("/discounts/validate", discountHandler.ValidateDiscountCode).Methods(http.MethodPost)

	// Public currency routes
	api.HandleFunc("/currencies", currencyHandler.ListEnabledCurrencies).Methods(http.MethodGet)
	api.HandleFunc("/currencies/default", currencyHandler.GetDefaultCurrency).Methods(http.MethodGet)
	api.HandleFunc("/currencies/convert", currencyHandler.ConvertAmount).Methods(http.MethodPost)

	// Public shipping routes
	api.HandleFunc("/shipping/options", shippingHandler.CalculateShippingOptions).Methods(http.MethodPost)
	// api.HandleFunc("/shipping/methods", shippingHandler.ListShippingMethods).Methods(http.MethodGet)
	// api.HandleFunc("/shipping/methods/{shippingMethodId:[0-9]+}", shippingHandler.GetShippingMethodByID).Methods(http.MethodGet)
	// api.HandleFunc("/shipping/rates/{shippingRateId:[0-9]+}/cost", shippingHandler.GetShippingCost).Methods(http.MethodPost)

	// Guest checkout routes (no authentication required)
	api.HandleFunc("/checkout", checkoutHandler.GetCheckout).Methods(http.MethodGet)
	api.HandleFunc("/checkout/items", checkoutHandler.AddToCheckout).Methods(http.MethodPost)
	api.HandleFunc("/checkout/items/{sku}", checkoutHandler.UpdateCheckoutItem).Methods(http.MethodPut)
	api.HandleFunc("/checkout/items/{sku}", checkoutHandler.RemoveFromCheckout).Methods(http.MethodDelete)
	api.HandleFunc("/checkout", checkoutHandler.ClearCheckout).Methods(http.MethodDelete)
	api.HandleFunc("/checkout/shipping-address", checkoutHandler.SetShippingAddress).Methods(http.MethodPut)
	api.HandleFunc("/checkout/billing-address", checkoutHandler.SetBillingAddress).Methods(http.MethodPut)
	api.HandleFunc("/checkout/customer-details", checkoutHandler.SetCustomerDetails).Methods(http.MethodPut)
	api.HandleFunc("/checkout/shipping-method", checkoutHandler.SetShippingMethod).Methods(http.MethodPut)
	api.HandleFunc("/checkout/currency", checkoutHandler.SetCurrency).Methods(http.MethodPut)
	api.HandleFunc("/checkout/discount", checkoutHandler.ApplyDiscount).Methods(http.MethodPost)
	api.HandleFunc("/checkout/discount", checkoutHandler.RemoveDiscount).Methods(http.MethodDelete)
	api.HandleFunc("/checkout/complete", checkoutHandler.CompleteOrder).Methods(http.MethodPost)
	// api.HandleFunc("/checkout/convert", checkoutHandler.ConvertGuestCheckoutToUserCheckout).Methods(http.MethodPost)

	// Routes with optional authentication (accessible via auth or checkout session)
	optionalAuth := api.PathPrefix("").Subrouter()
	optionalAuth.Use(authMiddleware.OptionalAuthenticate)
	optionalAuth.HandleFunc("/orders/{orderId:[0-9]+}", orderHandler.GetOrder).Methods(http.MethodGet)

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.Authenticate)

	// User routes
	protected.HandleFunc("/users/me", userHandler.GetProfile).Methods(http.MethodGet)
	protected.HandleFunc("/users/me", userHandler.UpdateProfile).Methods(http.MethodPut)
	protected.HandleFunc("/users/me/password", userHandler.ChangePassword).Methods(http.MethodPut)

	// Order routes (authenticated users only)
	protected.HandleFunc("/orders", orderHandler.ListOrders).Methods(http.MethodGet)

	// Admin routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminOnly)
	admin.HandleFunc("/users", userHandler.ListUsers).Methods(http.MethodGet)
	admin.HandleFunc("/orders", orderHandler.ListAllOrders).Methods(http.MethodGet)
	admin.HandleFunc("/orders/{orderId:[0-9]+}/status", orderHandler.UpdateOrderStatus).Methods(http.MethodPut)

	// Admin checkout routes
	admin.HandleFunc("/checkouts", checkoutHandler.ListAdminCheckouts).Methods(http.MethodGet)
	admin.HandleFunc("/checkouts/{checkoutId:[0-9]+}", checkoutHandler.GetAdminCheckout).Methods(http.MethodGet)
	admin.HandleFunc("/checkouts/{checkoutId:[0-9]+}", checkoutHandler.DeleteAdminCheckout).Methods(http.MethodDelete)

	// Admin currency routes
	admin.HandleFunc("/currencies/all", currencyHandler.ListCurrencies).Methods(http.MethodGet)
	admin.HandleFunc("/currencies", currencyHandler.CreateCurrency).Methods(http.MethodPost)
	admin.HandleFunc("/currencies", currencyHandler.UpdateCurrency).Methods(http.MethodPut)
	admin.HandleFunc("/currencies", currencyHandler.DeleteCurrency).Methods(http.MethodDelete)
	admin.HandleFunc("/currencies/default", currencyHandler.SetDefaultCurrency).Methods(http.MethodPut)

	// Admin email test route
	admin.HandleFunc("/test/email", emailTestHandler.TestEmail).Methods(http.MethodPost)

	// Admin category routes
	admin.HandleFunc("/categories", categoryHandler.CreateCategory).Methods(http.MethodPost)
	admin.HandleFunc("/categories/{id:[0-9]+}", categoryHandler.UpdateCategory).Methods(http.MethodPut)
	admin.HandleFunc("/categories/{id:[0-9]+}", categoryHandler.DeleteCategory).Methods(http.MethodDelete)

	// Shipping management routes (admin only)
	admin.HandleFunc("/shipping/methods", shippingHandler.CreateShippingMethod).Methods(http.MethodPost)
	// admin.HandleFunc("/shipping/methods/{shippingMethodId:[0-9]+}", shippingHandler.UpdateShippingMethod).Methods(http.MethodPut)
	admin.HandleFunc("/shipping/zones", shippingHandler.CreateShippingZone).Methods(http.MethodPost)
	// admin.HandleFunc("/shipping/zones", shippingHandler.ListShippingZones).Methods(http.MethodGet)
	// admin.HandleFunc("/shipping/zones/{shippingZoneId:[0-9]+}", shippingHandler.GetShippingZoneByID).Methods(http.MethodGet)
	// admin.HandleFunc("/shipping/zones/{shippingZoneId:[0-9]+}", shippingHandler.UpdateShippingZone).Methods(http.MethodPut)
	admin.HandleFunc("/shipping/rates", shippingHandler.CreateShippingRate).Methods(http.MethodPost)
	// admin.HandleFunc("/shipping/rates/{shippingRateId:[0-9]+}", shippingHandler.GetShippingRateByID).Methods(http.MethodGet)
	// admin.HandleFunc("/shipping/rates/{shippingRateId:[0-9]+}", shippingHandler.UpdateShippingRate).Methods(http.MethodPut)
	admin.HandleFunc("/shipping/rates/weight", shippingHandler.CreateWeightBasedRate).Methods(http.MethodPost)
	admin.HandleFunc("/shipping/rates/value", shippingHandler.CreateValueBasedRate).Methods(http.MethodPost)

	// Discount routes
	admin.HandleFunc("/discounts", discountHandler.CreateDiscount).Methods(http.MethodPost)
	admin.HandleFunc("/discounts/{discountId:[0-9]+}", discountHandler.UpdateDiscount).Methods(http.MethodPut)
	admin.HandleFunc("/discounts/{discountId:[0-9]+}", discountHandler.DeleteDiscount).Methods(http.MethodDelete)
	admin.HandleFunc("/discounts", discountHandler.ListDiscounts).Methods(http.MethodGet)
	admin.HandleFunc("/discounts/active", discountHandler.ListActiveDiscounts).Methods(http.MethodGet)
	admin.HandleFunc("/discounts/apply/{orderId:[0-9]+}", discountHandler.ApplyDiscountToOrder).Methods(http.MethodPost)
	admin.HandleFunc("/discounts/remove/{orderId:[0-9]+}", discountHandler.RemoveDiscountFromOrder).Methods(http.MethodDelete)
	admin.HandleFunc("/discounts/{discountId:[0-9]+}", discountHandler.GetDiscount).Methods(http.MethodGet)

	// Payment management routes (admin only)
	admin.HandleFunc("/payments/{paymentId}/capture", paymentHandler.CapturePayment).Methods(http.MethodPost)
	admin.HandleFunc("/payments/{paymentId}/cancel", paymentHandler.CancelPayment).Methods(http.MethodPost)
	admin.HandleFunc("/payments/{paymentId}/refund", paymentHandler.RefundPayment).Methods(http.MethodPost)
	admin.HandleFunc("/payments/{paymentId}/force-approve", paymentHandler.ForceApproveMobilePayPayment).Methods(http.MethodPost)

	// Payment provider management routes (admin only)
	admin.HandleFunc("/payment-providers", paymentProviderHandler.GetPaymentProviders).Methods(http.MethodGet)
	admin.HandleFunc("/payment-providers/enabled", paymentProviderHandler.GetEnabledPaymentProviders).Methods(http.MethodGet)
	admin.HandleFunc("/payment-providers/{providerType}/enable", paymentProviderHandler.EnablePaymentProvider).Methods(http.MethodPut)
	admin.HandleFunc("/payment-providers/{providerType}/configuration", paymentProviderHandler.UpdateProviderConfiguration).Methods(http.MethodPut)
	admin.HandleFunc("/payment-providers/{providerType}/webhook", paymentProviderHandler.RegisterWebhook).Methods(http.MethodPost)
	admin.HandleFunc("/payment-providers/{providerType}/webhook", paymentProviderHandler.DeleteWebhook).Methods(http.MethodDelete)
	admin.HandleFunc("/payment-providers/{providerType}/webhook", paymentProviderHandler.GetWebhookInfo).Methods(http.MethodGet)

	admin.HandleFunc("/products", productHandler.ListProducts).Methods(http.MethodGet)
	admin.HandleFunc("/products", productHandler.CreateProduct).Methods(http.MethodPost)
	admin.HandleFunc("/products/{productId:[0-9]+}", productHandler.UpdateProduct).Methods(http.MethodPut)
	admin.HandleFunc("/products/{productId:[0-9]+}", productHandler.DeleteProduct).Methods(http.MethodDelete)

	// Product variant routes
	admin.HandleFunc("/products/{productId:[0-9]+}/variants", productHandler.AddVariant).Methods(http.MethodPost)
	admin.HandleFunc("/products/{productId:[0-9]+}/variants/{variantId:[0-9]+}", productHandler.UpdateVariant).Methods(http.MethodPut)
	admin.HandleFunc("/products/{productId:[0-9]+}/variants/{variantId:[0-9]+}", productHandler.DeleteVariant).Methods(http.MethodDelete)
}

// GetContainer returns the dependency injection container
func (s *Server) GetContainer() container.Container {
	return s.container
}

// Start starts the server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
