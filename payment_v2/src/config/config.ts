import dotenv from 'dotenv';

dotenv.config();

export const config = {
    stripe: {
        secretKey: process.env.STRIPE_SECRET_KEY ?? '',
        successUrl: process.env.SUCCESS_URL ?? 'http://localhost:3000/success',
        cancelUrl: process.env.CANCEL_URL ?? 'http://localhost:3000/cancel',
    },
    rabbitmq: {
        url: process.env.RABBITMQ_URL ?? 'amqp://localhost:5672',
        paymentRequestQueue: 'payment_requests',
        paymentResponseQueue: 'payment_responses',
    },
}; 