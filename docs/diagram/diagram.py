from diagrams import Diagram, Cluster, Edge
from diagrams.programming.language import Go, Rust
from diagrams.onprem.database import PostgreSQL
from diagrams.onprem.client import Client
from diagrams.onprem.network import Internet
from diagrams.onprem.queue import RabbitMQ

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

            supabase = Client("Supabase\nStorage")

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

    # External APIs in separate clusters
    with Cluster("Stripe"):
        stripe_api = Internet("Stripe API")

    with Cluster("Auth0"):
        auth0_api = Internet("Auth0 API")

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
