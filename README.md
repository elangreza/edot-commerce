# edot-commerce

built with _\*Orchestrated Saga Pattern_\* consist of multiple microservices to handle e-commerce operations such as

- API service. It runs as a gateway to route requests to other services. This is the public access only Includes:
  - User authentication,
  - Product catalog,
  - Order processing,
  - Warehouse management.
- Order service,
- Product service,
- Shop service,
- Warehouse service

here's the list of technologies used in this project:

- API gateway: go-chi router
- Authentication: JWT tokens
- Communication between services: gRPC with go
- Database: can be run with either PostgreSQL or Sqlite3

here's the list of API endpoints exposed by the API service:

- `POST /auth/register` - Register a new user  
  | Field | Value |
  |-------------------|--------------------------------------------|
  | **Endpoint** | `POST /auth/register` |
  | **URL** | `http://localhost:8080/auth/register` |
  | **Content-Type** | `application/json` |
  | **Success Code** | `201 Created` |
  | **Description** | Registers a new user account with email, password, and name. |

  **Example:**

  ```bash
  curl --location 'http://localhost:8080/auth/register' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "email":"test@test.com",
      "password":"test",
      "name":"test"
  }'
  ```

---

- `POST /auth/login` - Login and obtain a JWT token  
  | Field | Value |
  |-------------------|--------------------------------------------|
  | **Endpoint** | `POST /auth/login` |
  | **URL** | `http://localhost:8080/auth/login` |
  | **Content-Type** | `application/json` |
  | **Success Code** | `200 OK` |
  | **Description** | Authenticates a user and returns a JWT for protected endpoints. |

  **Example:**

  ```bash
  curl --location 'http://localhost:8080/auth/login' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "email":"test@test.com",
      "password":"test"
  }'
  ```

---

- `GET /products` - Get a list of products  
  | Field | Value |
  |-------------------|--------------------------------------------|
  | **Endpoint** | `GET /products` |
  | **URL** | `http://localhost:8080/products` |
  | **Content-Type** | — |
  | **Success Code** | `200 OK` |
  | **Description** | Retrieves a paginated, optionally filtered list of products. Supports `page`, `limit`, and `search` query parameters. |

  **Example:**

  ```bash
  curl --location 'http://localhost:8080/products?page=2&limit=10&search=men'
  ```

---

- `POST /cart` - Add a product to the cart  
  | Field | Value |
  |-------------------|--------------------------------------------|
  | **Endpoint** | `POST /cart` |
  | **URL** | `http://localhost:8080/cart` |
  | **Content-Type** | `application/json` |
  | **Authorization** | `Bearer <JWT>` |
  | **Success Code** | `201 Created` |
  | **Description** | Adds a specified quantity of a product to the authenticated user’s cart. |

  **Example:**

  ```bash
  curl --location 'http://localhost:8080/cart' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer {{token from login API}}' \
  --data '{
      "product_id":"019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca2",
      "quantity":79
  }'
  ```

---

- `GET /cart` - Get the current cart contents  
  | Field | Value |
  |-------------------|--------------------------------------------|
  | **Endpoint** | `GET /cart` |
  | **URL** | `http://localhost:8080/cart` |
  | **Content-Type** | — |
  | **Authorization** | `Bearer <JWT>` |
  | **Success Code** | `200 OK` |
  | **Description** | Returns the full contents of the authenticated user’s shopping cart. |

  **Example:**

  ```bash
  curl --location 'http://localhost:8080/cart' \
  --header 'Authorization: Bearer {{token from login API}}'
  ```

---

- `POST /order` - Create a new order based on the cart  
  | Field | Value |
  |-------------------|--------------------------------------------|
  | **Endpoint** | `POST /order` |
  | **URL** | `http://localhost:8080/order` |
  | **Content-Type** | `application/json` |
  | **Authorization** | `Bearer <JWT>` |
  | **Success Code** | `201 Created` |
  | **Description** | Converts the user’s current cart into a confirmed order. Uses an `idempotency_key` to prevent duplicate submissions. |

  **Example:**

  ```bash
  curl --location 'http://localhost:8080/order' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer {{token from login API}}' \
  --data '{
      "idempotency_key":"75b12b36-8547-4c02-9783-d42007f6a92a"
  }'
  ```

---

- `POST /warehouse/status` - Set warehouse status (active/inactive)  
  | Field | Value |
  |-------------------|--------------------------------------------|
  | **Endpoint** | `POST /warehouse/status` |
  | **URL** | `http://localhost:8080/warehouse/status` |
  | **Content-Type** | `application/json` |
  | **Authorization** | `Bearer <JWT>` |
  | **Success Code** | `200 OK` |
  | **Description** | Updates the operational status (`is_active`) of a warehouse. Typically restricted to admin users. |

  **Example:**

  ```bash
  curl --location 'http://localhost:8080/warehouse/status' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer {{token from login API}}' \
  --data '{
      "warehouse_id":1,
      "is_active": true
  }'
  ```

---

- `POST /warehouse/transfer` - Transfer stock between warehouses  
  | Field | Value |
  |-------------------|--------------------------------------------|
  | **Endpoint** | `POST /warehouse/transfer` |
  | **URL** | `http://localhost:8080/warehouse/transfer` |
  | **Content-Type** | `application/json` |
  | **Authorization** | `Bearer <JWT>` |
  | **Success Code** | `200 OK` |
  | **Description** | Moves a specified quantity of a product from one warehouse to another. Requires admin privileges. |

  **Example:**

  ```bash
  curl --location 'http://localhost:8080/warehouse/transfer' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer {{token from login API}}' \
  --data '{
      "from_warehouse_id": 1,
      "to_warehouse_id": 3,
      "product_id": "019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca2",
      "quantity": 10
  }'
  ```
