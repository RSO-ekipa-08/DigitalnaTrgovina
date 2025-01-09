import { RabbitMQService } from './services/rabbitmq.service';

async function main() {
    const rabbitMQService = new RabbitMQService();

    try {
        await rabbitMQService.connect();
        await rabbitMQService.startListening();
        
        console.log('Payment service is running and listening for messages...');

        // Graceful shutdown
        process.on('SIGINT', async () => {
            console.log('Shutting down...');
            await rabbitMQService.close();
            process.exit(0);
        });
    } catch (error) {
        console.error('Failed to start the service:', error);
        process.exit(1);
    }
}

main();
