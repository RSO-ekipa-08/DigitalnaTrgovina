import Stripe from "stripe";
import { config } from "../config/config";

const stripe = new Stripe(config.stripe.secretKey);

export interface PaymentRequest {
  amount: number;
  currency: string;
  productName: string;
  quantity: number;
}

export class StripeService {
  async createCheckoutSession(paymentRequest: PaymentRequest): Promise<string> {
    try {
      const session = await stripe.checkout.sessions.create({
        line_items: [
          {
            price_data: {
              currency: paymentRequest.currency,
              product_data: {
                name: paymentRequest.productName,
              },
              unit_amount: paymentRequest.amount,
            },
            quantity: paymentRequest.quantity,
          },
        ],
        mode: "payment",
        success_url: config.stripe.successUrl,
        cancel_url: config.stripe.cancelUrl,
      });

      return session.url ?? "";
    } catch (error) {
      console.error("Error creating checkout session:", error);
      throw error;
    }
  }
}
