import * as amqp from 'amqplib';
import { config } from '../config/config';
import type { PaymentRequest } from './stripe.service';
import { StripeService } from './stripe.service';

export class RabbitMQService {
    private connection?: amqp.Connection;
    private channel?: amqp.Channel;
    private stripeService: StripeService;

    constructor() {
        this.stripeService = new StripeService();
    }

    async connect(): Promise<void> {
        try {
            this.connection = await amqp.connect(config.rabbitmq.url);
            this.channel = await this.connection.createChannel();

            // Ensure queues exist
            await this.channel.assertQueue(config.rabbitmq.paymentRequestQueue);
            await this.channel.assertQueue(config.rabbitmq.paymentResponseQueue);

            console.log('Connected to RabbitMQ');
        } catch (error) {
            console.error('Error connecting to RabbitMQ:', error);
            throw error;
        }
    }

    async startListening(): Promise<void> {
        if (!this.channel) throw new Error('RabbitMQ channel not initialized');

        this.channel.consume(config.rabbitmq.paymentRequestQueue, async (msg) => {
            if (!this.channel) throw new Error('RabbitMQ channel not initialized');
            if (!msg) return;
            try {
                const paymentRequest: PaymentRequest = JSON.parse(msg.content.toString());
                const checkoutUrl = await this.stripeService.createCheckoutSession(paymentRequest);

                // Send response back to the specified reply-to queue
                if (msg.properties.replyTo) {
                    this.channel.sendToQueue(
                        msg.properties.replyTo,
                        Buffer.from(JSON.stringify({ checkoutUrl })),
                        { correlationId: msg.properties.correlationId }
                    );
                } else {
                    console.error('No reply-to queue specified in the request');
                }

                // Acknowledge the message
                this.channel.ack(msg);
            } catch (error) {
                console.error('Error processing payment request:', error);
                // Negative acknowledge in case of error
                this.channel.nack(msg);
            }
        });
    }

    async close(): Promise<void> {
        await this.channel?.close();
        await this.connection?.close();
    }
} 