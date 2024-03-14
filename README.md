# URL Shortening Service

- This URL Shortening Service is a modern web application designed to create short URLs from longer ones, making it easier to share links online. 
- Built with Go and the Gin web framework, it leverages the Domain-Driven Design (DDD) approach for clarity, maintainability, and scalability. 
- Redis is used for efficient data storage and retrieval.

## Domain-Driven Design (DDD)

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

## Avoiding Collisions

Collision avoidance in URL shortening is crucial for ensuring each long URL is associated with a unique short code. A two-fold strategy was implemented to achieve this

- **Hashing and Encoding:** Initially, the system generate a candidate short code by hashing the original URL and encoding the hash into a Base62 string. This process ensures that most URLs will naturally map to unique short codes.
- **Uniqueness Check:** Before finalizing a short code, the system also checked its uniqueness within the Redis data store. If a collision is detected (the generated short code already exists), the system applied a sequence number to the original URL and regenerate the hash. This process repeats until a unique short code is found.

This method balances efficiency with the guarantee of uniqueness, allowing the service to scale while maintaining integrity.

## Getting Started

To run the URL Shortening Service locally:

1. Ensure Go and Redis are installed on your system.
2. Clone the repository and navigate to the project directory.
3. Install dependencies:
   ```
   go mod tidy
   ```
4. Start the server:
   ```
   go run main.go
   ```

The service will be available at `http://localhost:9000`. Use the following API endpoints to interact with the system:

1. **Add a URL:**
   ```
   curl --location 'http://localhost:9000/api/v1/url/add' \
   --header 'Content-Type: application/json' \
   --data '{
   "original_url": "https://research.tsmc.com/chinese/collaborations/academic/university-centers.html",
   "expiry": "2024-04-02T00:00:00Z"
   }'
   ```
    The response will include the shortened URL.
   ```
   {
   "shortened_url": "http://localhost:9000/api/v1/redirect/3EMjtvea"
   }
   ```

2. **Redirect to Original URL:**
   ```
   curl --location 'localhost:9000/api/v1/redirect/3EMjtvea'
   ```
3. **Display all the mapped URLs:** 
   ```
   curl --location 'http://localhost:9000/api/v1/url/display'
   ```
   The response will include the shortened URL.
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
4. **Swagger API Documentation:**
   - The Swagger API documentation is available at `http://localhost:9000/swagger/index.html`
   - ![screen shot of redis client](https://github.com/terenzio/URL-Shortening-Service/blob/main/screenshots/Redis_ScreenShot.png?raw=true)
   - For convenience a PDF version can also be seen here, without having to run the application: [PDF LINK](https://github.com/terenzio/URL-Shortening-Service/blob/main/screenshots/Swagger_ScreenShot.pdf)
5. **Redis Data Store:**
   - The Redis data store can be accessed using the Redis CLI or a GUI tool like RedisInsight. 
   - The data store will contain the original URLs, their shortened codes, and expiry dates.
   - For convenience a screenshot of the Redis data store can be seen below, without having to run the application:
   - ![screen shot of redis client](https://github.com/terenzio/URL-Shortening-Service/blob/main/screenshots/Redis_ScreenShot.png?raw=true)

     