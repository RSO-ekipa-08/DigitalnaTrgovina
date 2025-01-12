from diagrams import Diagram, Cluster, Edge
from diagrams.programming.language import Go, Rust
from diagrams.onprem.database import PostgreSQL
from diagrams.onprem.client import Client
from diagrams.onprem.network import Internet
from diagrams.onprem.queue import RabbitMQ
from diagrams.custom import Custom
from urllib.request import urlretrieve

# Download icons for external services
stripe_url = "https://cdn.prod.website-files.com/635637c5d13ee9bd2b5feb65/635a6521e448c975009bae34_6357d2186e13e5f26c16d622_stripe_logo_icon_167962.png"
stripe_icon = "stripe.png"
urlretrieve(stripe_url, stripe_icon)

auth0_url = "https://play-lh.googleusercontent.com/So22eXt1Cc7-9ishK7DAoBaCUqnfuehrxyA_kezuhspg5gg526eMIQeppffxFgZsjAXn"
auth0_icon = "auth0.png"
urlretrieve(auth0_url, auth0_icon)

supabase_url = "https://pipedream.com/s.v0/app_1dBhP3/logo/96"
supabase_icon = "supabase.png"
urlretrieve(supabase_url, supabase_icon)

# Create diagram
with Diagram(show=False, direction="LR", outformat="pdf", filename="arhitektura", graph_attr={
    "fontsize": "18",
    "bgcolor": "transparent"
}, node_attr={"fontsize": "18"}) as d:

    # Frontend
    client = Client("Uporabnik")

    # Microservices
    with Cluster("Mikrostoritve"):
        # API Gateway
        with Cluster("API Gateway"):
            gateway = Go("API Gateway")

        # Queue inside microservices cluster
        queue = RabbitMQ("RabbitMQ")

        # App Service
        with Cluster("Aplikacijska storitev"):
            app = Go("Aplikacijska\nstoritev")
            app_db = PostgreSQL("Aplikacijska\nPB")

            supabase = Custom("Supabase\nStorage", supabase_icon)

            app - app_db
            app - supabase

        # Reviews Service
        with Cluster("Storitev za ocene"):
            reviews = Rust("Storitev za\nocene")
            reviews_db = PostgreSQL("Ocene PB")

            reviews - reviews_db

        # Payment Service
        with Cluster("Plačilna storitev"):
            payment = Go("Plačilna\nstoritev")

        # Auth Service
        with Cluster("Avtentikacijska storitev"):
            auth = Go("Avtentikacijska\nstoritev")

    # Replace external service nodes with custom icons
    with Cluster("Stripe"):
        stripe_api = Custom("Stripe API", stripe_icon)

    with Cluster("Auth0"):
        auth0_api = Custom("Auth0 API", auth0_icon)

    # Client to Gateway
    client >> Edge(label="HTTP") >> gateway

    # Gateway to Services
    gateway >> Edge(label="REST") >> app
    gateway >> Edge(label="GraphQL") >> reviews
    gateway >> Edge(label="gRPC") >> auth

    # Payment Service Communication
    app - Edge() - queue
    queue - Edge() - payment

    # Service Dependencies
    auth << Edge(label="gRPC") << app
    auth << Edge(label="gRPC") << payment
    auth << Edge(label="gRPC") << reviews

    # External API connections
    payment - Edge(label="HTTP") - stripe_api
    auth - Edge(label="HTTP") - auth0_api
