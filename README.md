# URL Shortening Service by Terence Liu for TSMC

- This URL Shortening Service is a modern web application designed to create short URLs from longer ones, making it easier to share links online. 
- Built with Go and the Gin web framework, it leverages the Domain-Driven Design (DDD) approach for clarity, maintainability, and scalability. 
- Redis is used for efficient data storage and retrieval.

## Features

1. A client (user) enters a long URL into the system and the system returns a shortened
   URL [See API Endpoints](#api-endpoints) (Point 1) 
2. The short URL should be readable [See API Endpoints](#api-endpoints) (Point 4)
3. The short URL should be collision-free [See API Endpoints](#api-endpoints) (Point 1)
4. The short URL should be non-predictable [See API Endpoints](#api-endpoints) (Point 1)
5. The client should be able to choose a custom short URL [See API Endpoints](#api-endpoints) (Point 2)
6. The client visiting the short URL must be redirected to the original long URL [See API Endpoints](#api-endpoints) (Point 3)
7. The client optionally defines the expiry time of the short URL [See API Endpoints](#api-endpoints) (Point 2)

##  Technical Design: Avoiding Collisions

Collision avoidance in URL shortening is crucial for ensuring each long URL is associated with a unique short code. A two-fold strategy was implemented to achieve this:

- **Hashing and Encoding:** Initially, the system generate a candidate short code by hashing the original URL and encoding the hash into a Base62 string. This process ensures that most URLs will naturally map to unique short codes.
- **Uniqueness Check:** Before finalizing a short code, the system also checked its uniqueness within the Redis data store. If a collision is detected (the generated short code already exists), the system applied a sequence number to the original URL and regenerate the hash. This process repeats until a unique short code is found.

This method balances efficiency with the guarantee of uniqueness, allowing the service to scale while maintaining integrity.

## System Analysis:

- **Traffic Estimates: 500 Million new URLs per month with 100:1 read/write ratio**
  - New URLs shortening per second = 500 million / (30 days * 24 hours * 60 minutes * 60 seconds) =~ 200 URLs per second
  - URLs redirection per second = 100 read ratio * 200 URLs per second = 20K redirections per second
- **Storage Estimates: for 5 years** 
  - 500 Million new URLs per month which is kept for 5 years
  - Total number of objects = 500 million * 5 years * 12 months = 30 billion
  - If we assume 0.5 KB per object, the total storage required = 30 billion * 0.5 KB = 15TB
- **Bandwidth Estimates: 100:1 read to write ratio**
  - Incoming read data = 200 new URLs per second * 0.5 KB = 100 KB per second
  - Outgoing write data = 20 000 redirections per second * 0.5 KB = 10 MB per second
- **Memory Estimates: Based on the Pareto Principle 20% of the total URLs are accessed per day**
  - 20k redirections per second * 3600 seconds * 24 hours =~ 1.7 billion redirections per day
  - Following the Pareto Principle we need to cache 20% of the daily redirections, this will lead to 0.2 * 1.7 billion * 0.5 KB =~ 170 GB of memory usage


## Architecture Philosophy: Domain-Driven Design (DDD)

- DDD was adopted to align the development practices with the requirements closely.
- This design philosophy emphasizes placing the project's primary focus on the core domain and domain logic, then builds outwards to application logic and infrastructure.

### Structure

The project is organized into several layers according to DDD principles:

- **Domain Layer:** Contains the core business logic, entities (such as URLs), and repository interfaces. It represents the heart of the business model.
- **Application Layer:** Encapsulates the application's use cases. It uses domain entities to fulfill the functionalities required by the application, such as URL shortening and redirection.
- **Infrastructure Layer:** Implements technical details that support the application, including data persistence (Redis repository) and the HTTP server setup (using Gin for routing and middleware).

### Advantages

- **Enhanced Modularity:** Separating concerns makes the system easier to understand, develop, and test.
- **Improved Scalability:** By decoupling core logic from infrastructure, the system can easily adapt to new requirements or technologies.
- **Focused Business Logic:** DDD helps us stay aligned with business objectives, making our application more effective and adaptable.

## Getting Started

To run the URL Shortening Service locally:

1. Ensure Go 1.20+ and Redis are installed on your system. 
   ```
   // For Mac users:
    > brew install go@1.20
    > brew install redis
    > redis-server
   ```
2. Clone the repository and navigate to the project directory.
3. Install dependencies:
   ```
    > go mod tidy
   ```
4. Start the server:
   ```
    > go run main.go
   ```
   Here is a screenshot of the GIN server running on Port 9000:
   ![screen shot of server running](https://github.com/terenzio/URL-Shortening-Service/blob/main/screenshots/GinServer_ScreenShot.png?raw=true)

## API Endpoints

The service will be available at `http://localhost:9000`. Use the following API endpoints and tools to interact with the system:

1. **Add a URL:**
   - Functional Requirement 1: A client (user) enters a long URL into the system and the system returns a shortened
      URL
   - Functional Requirement 3: The short URL should be collision-free
   - Functional Requirement 4: The short URL should be non-predictable
      ```
      curl --location 'http://localhost:9000/api/v1/url/add' \
      --header 'Content-Type: application/json' \
      --data '{
          "original_url": "https://research.tsmc.com/chinese/collaborations/academic/university-centers.html"
      }'
      ```
       The response will include the shortened URL.
      ```
      {
          "shortened_url": "http://localhost:9000/api/v1/redirect/3EMjtvea"
      }
      ```
2. **Add a URL with custom short code and expiry time:**
   - Functional Requirement 5: The client should be able to choose a custom short URL
   - Functional Requirement 7: The client optionally defines the expiry time of the short URL
      ```
      curl --location 'http://localhost:9000/api/v1/url/add' \
      --header 'Content-Type: application/json' \
      --data '{
          "original_url": "https://research.tsmc.com/chinese/collaborations/academic/university-centers1.html",
          "expiry": "2024-04-02T00:00:00Z",
          "custom_short_code": "abcde1"
      }'
      ```
      The response will include the shortened URL and the customized values.
       ```
       {
           "original_url": "https://research.tsmc.com/chinese/collaborations/academic/university-centers1.html",
           "expiry": "2024-04-02T00:00:00Z",
           "shortened_url": "http://localhost:9000/api/v1/redirect/abcde1"
       }
       ```

3. **Redirect to Original URL:**
    - Functional Requirement 6: The client visiting the short URL must be redirected to the original long URL
      ```
      curl --location 'localhost:9000/api/v1/redirect/3EMjtvea'
      ```
4. **Display all the mapped URLs:** 
   - Functional Requirement 2: The short URL should be readable
      ```
      curl --location 'http://localhost:9000/api/v1/url/display'
      ```
      The response will include the shortened code for easy readibility.
      ```
      [
        {
          "short_code": "2LzboGMR",
          "original_url": "https://research.tsmc.com/chinese/collaborations/academic/university-centers.html",
          "expiry": "2024-06-02T07:59:59.860239+08:00"
        },
        {
           "short_code": "4uODYpIv",
           "original_url": "https://www.tsmc.com/chinese/aboutTSMC/company_profile",
           "expiry": "2024-05-02T07:59:59.861373+08:00"
        }
      ]
      ```
5. **Swagger API Documentation:**
   - The Swagger API documentation is available at `http://localhost:9000/swagger/index.html`
   - ![screen shot of swagger](https://github.com/terenzio/URL-Shortening-Service/blob/main/screenshots/Swagger_ShortScreenShot.png?raw=true)
   - For convenience a PDF version can also be seen here, without having to run the application: [PDF LINK](https://github.com/terenzio/URL-Shortening-Service/blob/main/screenshots/Swagger_FullScreenShot.pdf)
6. **Redis Data Store:**
   - The Redis data store can be accessed using the Redis CLI or a GUI tool like RedisInsight. 
   - The data store will contain the original URLs, their shortened codes, and expiry dates.
   - For convenience a screenshot of the Redis data store can be seen below, without having to run the application:
   - ![screen shot of redis client](https://github.com/terenzio/URL-Shortening-Service/blob/main/screenshots/Redis_ScreenShot.png?raw=true)


## Testing

- Tests are included to ensure the application's correctness and robustness.
- A third-party library called `testify` was used to simplify the testing process.
- Along with another one called `Miniredis` to mock the Redis server. [Miniredis](https://github.com/alicebob/miniredis)
- Here are the results after running the tests: 
  ```
  > cd infrastructure/redis
  > go test -v
  
      === RUN   TestURLRepository_Store
      === RUN   TestURLRepository_Store/Valid_URL_with_24-hour_Expiry
      === RUN   TestURLRepository_Store/Expired_URL
      --- PASS: TestURLRepository_Store (0.00s)
      --- PASS: TestURLRepository_Store/Valid_URL_with_24-hour_Expiry (0.00s)
      --- PASS: TestURLRepository_Store/Expired_URL (0.00s)
      PASS
      ok      github.com/terenzio/URL-Shortening-Service/infrastructure/redis 0.222s
      ❯ go test -v
      === RUN   TestURLRepository_Store
      === RUN   TestURLRepository_Store/Valid_URL_with_24-hour_Expiry
      === RUN   TestURLRepository_Store/Expired_URL
      --- PASS: TestURLRepository_Store (0.00s)
      --- PASS: TestURLRepository_Store/Valid_URL_with_24-hour_Expiry (0.00s)
      --- PASS: TestURLRepository_Store/Expired_URL (0.00s)
      === RUN   TestURLRepository_FindByShortCode
      === RUN   TestURLRepository_FindByShortCode/URL_Found
      === RUN   TestURLRepository_FindByShortCode/URL_Not_Found
      --- PASS: TestURLRepository_FindByShortCode (0.00s)
      --- PASS: TestURLRepository_FindByShortCode/URL_Found (0.00s)
      --- PASS: TestURLRepository_FindByShortCode/URL_Not_Found (0.00s)
      PASS
      ok      github.com/terenzio/URL-Shortening-Service/infrastructure/redis 0.262s
      ❯ go test -v
      === RUN   TestURLRepository_Store
      === RUN   TestURLRepository_Store/Valid_URL_with_24-hour_Expiry
      === RUN   TestURLRepository_Store/Expired_URL
      --- PASS: TestURLRepository_Store (0.00s)
      --- PASS: TestURLRepository_Store/Valid_URL_with_24-hour_Expiry (0.00s)
      --- PASS: TestURLRepository_Store/Expired_URL (0.00s)
      === RUN   TestURLRepository_FindByShortCode
      === RUN   TestURLRepository_FindByShortCode/URL_Found
      === RUN   TestURLRepository_FindByShortCode/URL_Not_Found
      --- PASS: TestURLRepository_FindByShortCode (0.00s)
      --- PASS: TestURLRepository_FindByShortCode/URL_Found (0.00s)
      --- PASS: TestURLRepository_FindByShortCode/URL_Not_Found (0.00s)
      === RUN   TestURLRepository_IsUnique
      === RUN   TestURLRepository_IsUnique/ShortCode_is_Unique
      === RUN   TestURLRepository_IsUnique/ShortCode_is_Not_Unique
      --- PASS: TestURLRepository_IsUnique (0.00s)
      --- PASS: TestURLRepository_IsUnique/ShortCode_is_Unique (0.00s)
      --- PASS: TestURLRepository_IsUnique/ShortCode_is_Not_Unique (0.00s)
      === RUN   TestURLRepository_FetchAll
      === RUN   TestURLRepository_FetchAll/Successfully_Fetch_All_URLs
      ttl:  -1ns
      ttl:  -1ns
      --- PASS: TestURLRepository_FetchAll (0.00s)
      --- PASS: TestURLRepository_FetchAll/Successfully_Fetch_All_URLs (0.00s)
      PASS
      ok      github.com/terenzio/URL-Shortening-Service/infrastructure/redis 0.229s
  ```
  