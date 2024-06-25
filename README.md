# Project: Receipt Processor Service

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

# Final Note
- **Overview:** To meet the requirement that data need not persist across server restarts, I utilized an in-memory map in Go to store receipts. This approach avoids the overhead of setting up a database like PostgreSQL. However, since Go maps are not concurrency-safe, I implemented a mutex lock `rs.mu.Lock()` to ensure safe concurrent access to the `inMemoryReceiptMap`. The lock is acquired before performing get/put operations and released afterward using defer `rs.mu.Unlock()`.
- **Scalability Considerations:** While this solution works for the given requirements, it does have limitations regarding horizontal scalability. Running multiple instances of the application would lead to each instance having its own in-memory store, resulting in data integrity issues due to the lack of a single source of truth.
- **Reflection:** Working on this assignment was an enjoyable experience. It provided a great opportunity to refresh my knowledge of best production code practices and reinforced the importance of a layered architecture in application development.
