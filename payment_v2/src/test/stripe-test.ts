import { StripeService, type PaymentRequest } from '../services/stripe.service';

async function testStripeCheckout() {
    console.log('Začenjam test Stripe checkout funkcionalnosti...');

    const stripeService = new StripeService();

    // Testni podatki
    const testPayment: PaymentRequest = {
        amount: 1500, // 15.00 EUR
        currency: 'eur',
        productName: 'Test Produkt',
        quantity: 1
    };

    try {
        console.log('Pošiljam zahtevo za checkout:', testPayment);
        const checkoutUrl = await stripeService.createCheckoutSession(testPayment);

        if (checkoutUrl && checkoutUrl.startsWith('https://checkout.stripe.com')) {
            console.log('✅ Test uspešen! Checkout URL ustvarjen:', checkoutUrl);
        } else {
            console.error('❌ Test neuspešen! Neveljaven checkout URL:', checkoutUrl);
        }
    } catch (error) {
        console.error('❌ Test neuspešen! Napaka:', error);
    }
}

// Poženi test
console.log('=== Stripe Integration Test ===');
testStripeCheckout()
    .then(() => console.log('Test zaključen.'))
    .catch(error => console.error('Nepričakovana napaka:', error));
