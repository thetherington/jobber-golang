package stripe

import (
	"fmt"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/customer"
	"github.com/stripe/stripe-go/v79/paymentintent"
	"github.com/stripe/stripe-go/v79/refund"
	"github.com/thetherington/jobber-order/internal/adapters/config"
)

type StripeService struct {
}

func NewStripePayment(config *config.Stripe) *StripeService {
	stripe.Key = config.Key

	return &StripeService{}
}

func (s *StripeService) CreateCustomer(email string, buyerId string) (string, error) {
	params := &stripe.CustomerParams{
		Description: stripe.String("Jobber customer"),
		Email:       stripe.String(email),
		Metadata: map[string]string{
			"buyerId": buyerId,
		},
		PreferredLocales: stripe.StringSlice([]string{"en"}),
	}

	c, err := customer.New(params)
	if err != nil {
		return "", err
	}

	return c.ID, nil
}

func (s *StripeService) SearchCustomers(email string) (string, error) {
	params := &stripe.CustomerSearchParams{
		SearchParams: stripe.SearchParams{
			Query: fmt.Sprintf("email:'%s'", email),
		},
	}
	params.Single = true

	result := customer.Search(params)

	for result.Next() {
		c := result.Customer()
		if c.Email == email {
			return c.ID, nil
		}
	}

	if err := result.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func (s *StripeService) CreatePaymentIntent(customerId string, price float32) (string, string, error) {
	amount := int64(price * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Customer: stripe.String(customerId),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	result, err := paymentintent.New(params)
	if err != nil {
		return "", "", err
	}

	return result.ID, result.ClientSecret, nil
}

func (s *StripeService) RefundOrder(pi string) error {
	params := &stripe.RefundParams{
		PaymentIntent: &pi,
	}

	_, err := refund.New(params)
	if err != nil {
		return err
	}

	return nil
}
