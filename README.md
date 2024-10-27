# Auctions-go-routines

The Auctions project is a Go-based application that allows users to create and manage auctions, place bids, and find winning bids. The project uses MongoDB as the database and Gin as the web framework.

### Running the Application
To run the application using Docker Compose, follow these steps:

1. Clone the repository:

```sh
git clone https://github.com/carlosmeds/auctions-go-routines.git
cd auctions-go-routines
```

2. Set up environment variables: Ensure that the `.env` file in the auction directory is correctly configured

3. Build and run the application: Use Docker Compose to build and run the application:

```sh
docker-compose up --build
```

This command will build the Docker images and start the containers for the application and MongoDB.

4. Access the application: Once the containers are up and running, you can access the application at http://localhost:8080.

### Services
The application exposes the following services:

* Auction Service:

   - `GET /auction`: List all auctions.
   - `GET /auction/:auctionId`: Get details of a specific auction.
   - `POST /auction`: Create a new auction.
   - `GET /auction/winner/:auctionId`: Get the winning bid for a specific auction.

* Bid Service:

    - `POST /bid`: Place a new bid.
    - `GET /bid/:auctionId`: List all bids for a specific auction.

* User Service:

    - `GET /user/:userId`: Get details of a specific user.
