from diagrams import Diagram, Cluster
from diagrams.programming.language import Go, Rust
from diagrams.onprem.database import PostgreSQL
from diagrams.onprem.client import Client
from diagrams.onprem.network import Internet
from diagrams.azure.storage import BlobStorage

with Diagram("Arhitektura DigitalnaTrgovina", show=False, direction="TB", outformat="pdf", filename="arhitektura", graph_attr={
    "fontsize": "20",
    "bgcolor": "transparent"
}, node_attr={"fontsize": "20"}) as d:

    # Frontend
    client = Client("Spletna Aplikacija")

    # Microservices
    with Cluster("Mikrostoritve"):
        # App Service
        with Cluster("Aplikacijska Storitev"):
            app = Go("Aplikacijska\nStoritev")
            app_db = PostgreSQL("Aplikacijska\nPB")

            with Cluster("Shramba Datotek"):
                minio = Client("MinIO")
                azure = BlobStorage("Azure Blob\nStorage")

            app - app_db
            app - minio
            app - azure

        # Reviews Service
        with Cluster("Storitev za Ocene"):
            reviews = Rust("Storitev za\nOcene")
            reviews_db = PostgreSQL("Ocene PB")

            reviews - reviews_db

        # Payment Service
        with Cluster("Plačilna Storitev"):
            payment = Go("Plačilna\nStoritev")
            payment_db = PostgreSQL("Plačilna PB")

            payment - payment_db

        # Auth Service
        with Cluster("Avtentikacijska Storitev"):
            auth = Go("Avtentikacijska\nStoritev")
            auth_db = PostgreSQL("Avtentikacijska\nPB")

            auth - auth_db

    # External APIs
    with Cluster("Zunanje Storitve"):
        stripe_api = Internet("Stripe API")
        auth0_api = Internet("Auth0 API")

    # Service Connections
    client >> app
    client >> reviews
    client >> payment
    client >> auth

    # Cross-service communication
    app << reviews
    app << payment
    auth << app
    auth << payment
    auth << reviews

    # External API connections
    payment - stripe_api
    auth - auth0_api
