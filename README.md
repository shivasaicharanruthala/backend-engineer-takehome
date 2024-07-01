# Project: Receipt Processor Service
- Build a webservice that fulfils the documented API. The API is described below. A formal definition is provided in the [api.yml](https://github.com/fetch-rewards/receipt-processor-challenge/blob/main/api.yml) file
- Data does not need to persist when your application stops. It is sufficient to store information in memory.

These rules collectively define how many points should be awarded to a receipt.

- One point for every alphanumeric character in the retailer name.
- 50 points if the total is a round dollar amount with no cents.
- 25 points if the total is a multiple of `0.25`.
- 5 points for every two items on the receipt.
- If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
- 6 points if the day in the purchase date is odd.
- 10 points if the time of purchase is after 2:00pm and before 4:00pm.

## Examples
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```
```text
Total Points: 28
Breakdown:
     6 points - retailer name has 6 characters
    10 points - 4 items (2 pairs @ 5 points each)
     3 Points - "Emils Cheese Pizza" is 18 characters (a multiple of 3)
                item price of 12.25 * 0.2 = 2.45, rounded up is 3 points
     3 Points - "Klarbrunn 12-PK 12 FL OZ" is 24 characters (a multiple of 3)
                item price of 12.00 * 0.2 = 2.4, rounded up is 3 points
     6 points - purchase day is odd
  + ---------
  = 28 points
```

```json
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
```

```text
Total Points: 109
Breakdown:
    50 points - total is a round dollar amount
    25 points - total is a multiple of 0.25
    14 points - retailer name (M&M Corner Market) has 14 alphanumeric characters
                note: '&' is not alphanumeric
    10 points - 2:33pm is between 2:00pm and 4:00pm
    10 points - 4 items (2 pairs @ 5 points each)
  + ---------
  = 109 points
```


#  Notes
- **Design Approach:** To meet the requirement that data need not persist across server restarts, I utilized an in-memory map in Go to store receipts. This approach avoids the overhead of setting up a database like PostgreSQL. However, **since Go maps are not concurrency-safe, I implemented a mutex lock `rs.mu.Lock()` to ensure safe concurrent access to the `inMemoryReceiptMap`. The lock is acquired before performing get/put operations and released afterward using `defer rs.mu.Unlock()`.**
- **Code Structure:** The codebase is well-structured and clean, organized into different packages. It uses a layered architecture with the following components:
  - **Handler Layer:** Handles requests and validates request parameters.
  - **Service Layer:** Contains the business logic.
  - **Data Layer**: Manages data operations.
  - Each layer interacts with the interfaces of the next layer, and dependency injection is used to inject objects from one layer into another, achieving inversion of control.
- **Scalability Considerations:** While this solution works for the given requirements, it does have limitations regarding horizontal scalability. Running multiple instances of the application would lead to each instance having its own in-memory store, resulting in data integrity issues due to the lack of a single source of truth.
- **Reflection:** While this layered architecture may seem like over-engineering for a simple task, but it is designed with future code scalability in mind. Additionally, a logger and custom error messages are implemented to enhance robustness and maintainability. Other than that Working on this assignment was an enjoyable experience. It provided a great opportunity to refresh my knowledge of best production code practices and reinforced the importance of a layered architecture in application development.


## How to run
### Installation
```
1. Install docker on the desktop - https://docs.docker.com/compose/install/
2. Check for docker & docker-compose version -  `docker --version`, `docker-compose --version`
```

### Running the project
1. Clone the project
> git clone git@github.com:shivasaicharanruthala/backend-engineer-takehome.git
2. Go project directory
> cd backend-engineer-takehome
3. Run the docker-compose file
> docker-compose up -d
4. Check the logs of application container 
```text
  - check for active logs: `docker ps`
  - check for container id of the image: `docker ps --filter "ancestor=shiva5128/backend-engineer-takehome:latest"`
  - check for logs: `docker logs <container-id>`
  
  Logs are of these format 
  2024/06/25 03:42:16 INFO: {"level":"INFO","msg":"Logger initialized successfully"}
  2024/06/25 03:42:16 INFO: {"level":"INFO","msg":"Receipts Server starting to listen on port 8080"}
```


## How to test
### Using CuRL
1. Endpoint: Process Receipts
   - Path: `/v1/receipts/process`
   - Method: `POST` 
```bash
curl --location --request POST 'http://localhost:8080/v1/receipts/process' \
--header 'Content-Type: application/json' \
--data '{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },
    {
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },
    {
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },
    {
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },
    {
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}' -i -v
```

2. verify the response sent back after storing the receipt
```bash
{"id": "370fc237-9d4c-4d2f-a056-023080c755e2"}
```

3. Endpoint: Get Points
    - Path: `/v1/receipts//{id}/points`
    - Method: `GET` 
- Replace `{id}` in path `http://localhost:8080/v1/receipts/{id}/points` with the uuid seen in step 2.
- execute below GET CuRL 
```bash
curl -v -X GET 'http://localhost:8080/v1/receipts/{id}/points' -i -v
```

4. Verify the response got back after executing GET CuRL
```bash
{ "points" : "28" }
```

### Using Postman
![postman_testing.gif](tests%2Fpostman_testing.gif)

### Using Unit & Integration tests
- You can go through individual test cases of a function in each package or to run all test cases in one go execute below command
> go test -v ./...
