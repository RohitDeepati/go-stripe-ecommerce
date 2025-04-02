package stripe

import (
	"encoding/json"
	"fmt"
	"go-ecommerce/beans"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	// "github.com/stripe/stripe-go/checkout/session"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/webhook"
)

type StripeDB struct {
	db *sqlx.DB
}

func StripeHandler(db *sqlx.DB) StripeDB {
	return StripeDB{
		db: db,
	}
}


func (h *StripeDB) RouteGroup(r *gin.Engine) {
	routeGroup := r.Group("/payments")
	// routeGroup.POST("/create-payment-intent", h.createPaymentIntent)
	routeGroup.POST("/webhook", h.handleStripeWebhook)
	routeGroup.POST("/checkout-session", h.handleCheckoutSession)
	routeGroup.GET("/order/success", h.getCheckOutSessionDetails)
}

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == ""{
		log.Fatal("Stripe secret key is not set")
	}
}


// func (h *StripeDB) createPaymentIntent(ctx *gin.Context) {
// 	var requestData beans.StripePayment


// 	if err := ctx.ShouldBindJSON(&requestData); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}


// 	params := &stripe.PaymentIntentParams{
// 		Amount:   stripe.Int64(requestData.Amount),
// 		Currency: stripe.String(string(stripe.CurrencyUSD)),
// 		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
// 			Enabled: stripe.Bool(true),
// 		},
// 	}

// 	if len(requestData.PaymentMethodTypes) > 0{
// 		params.PaymentMethodTypes = stripe.StringSlice(requestData.PaymentMethodTypes)
// 	}

	
// 	if requestData.CustomerID != "" {
// 		params.Customer = stripe.String(requestData.CustomerID)
// 	}

// 	result, err := paymentintent.New(params)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}


// 	ctx.JSON(http.StatusOK, gin.H{"client_secret": result})
// }

func (h *StripeDB) handleStripeWebhook(ctx *gin.Context) {
	const maxBodyBytes = int64(65536) 
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxBodyBytes)

	// Read the request body
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Unable to read request body"})
		return
	}

	// Deserialize the event from the payload
	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("Webhook error while parsing JSON: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse webhook JSON"})
		return
	}

	// Log the signature for debugging
	signatureHeader := ctx.GetHeader("Stripe-Signature")
	log.Printf("Stripe-Signature: %s\n", signatureHeader)

	// Verify the webhook signature
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	event, err = webhook.ConstructEventWithOptions(
    payload,
    signatureHeader,
    webhookSecret,
    webhook.ConstructEventOptions{
        IgnoreAPIVersionMismatch: true,
    },
)

if err != nil {
    log.Printf("⚠️ Webhook signature verification failed: %v\n", err)
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook signature"})
    return
}

	fmt.Println("webhook-secret", webhookSecret)

	// Process the event based on its type
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("Error parsing payment intent: %v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse payment intent"})
			return
		}
		log.Printf("------checkoutsession ID------: %s\n", session.ID)
		// Further handling for payment intent success (e.g., updating the database)
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "paymentStatus": "paid"})
		
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			log.Printf("Error parsing payment intent: %v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse payment intent"})
			return
		}
		log.Printf("------Payment Intent: %v\n", paymentIntent)
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "paymentStatus": "paid"})

	case "payment_intenet.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil{
			fmt.Printf("error parsing the payment intent: %v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("---payment failed intent: %v", err)

	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"Unandled event type"})
	}

	// Respond to Stripe that the event was successfully processed
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h * StripeDB) handleCheckoutSession(ctx *gin.Context){
	var cartItems []beans.Cart

	if err := ctx.ShouldBindJSON(&cartItems); err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	
	var lineItems []*stripe.CheckoutSessionLineItemParams
	var totalAmount int64 = 0

	for _, item := range cartItems{
		lineItem := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				UnitAmount: stripe.Int64(item.Price),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.Name),
					Description: stripe.String(item.Title),
					Images: stripe.StringSlice([]string{item.ImageURL}),
				},
			},
			Quantity: stripe.Int64(item.Quantity),
		}
		lineItems = append(lineItems, lineItem)
		totalAmount += item.Price * int64(item.Quantity)
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("http://localhost:5173/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL: stripe.String("http://localhost:5173/cart"),
		LineItems: lineItems,
	}

	session, err := session.New(params)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	fmt.Print("session_id", session.URL)
	ctx.JSON(http.StatusOK, gin.H{ "session":session})
}

func (h *StripeDB) getCheckOutSessionDetails(ctx *gin.Context){
	sessionID, _ := ctx.GetQuery(("session_id"))

	if sessionID == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "session_id parameter is required"})
		return
	}

	session, err := session.Get(sessionID, nil)
	if err != nil{
		log.Printf("Error fetching checkout session: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch sessionid"})
		return
	}
	log.Printf("payment status for session %s: %s", session.ID, session.PaymentStatus)

	ctx.JSON(http.StatusOK, gin.H{
		"paymentStatus": session.PaymentStatus,
	})

}
