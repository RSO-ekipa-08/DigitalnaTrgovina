# Auth0 Authentication Microservice

## Running the App

To run the app, make sure you have **go** installed.

Add `.env` file and provide your Auth0 credentials.

```bash
# .env

AUTH0_CLIENT_ID=7QJuJ3TENmqgqqJPa2ayKVpA5pchLdDd
AUTH0_DOMAIN=samolego.eu.auth0.com
AUTH0_CLIENT_SECRET=1gkzkJHP2GT3lO9V8j6tsouUtpocSwnz2UaUWBFImwyUxxEs1yGKtCOiQptXwdcW
AUTH0_CALLBACK_URL=http://localhost:3000/callback
```

Once you've set your Auth0 credentials in the `.env` file, run `go mod vendor` to download the Go dependencies.

Run `go run main.go` to start the app and navigate to [http://localhost:3000/](http://localhost:3000/).
