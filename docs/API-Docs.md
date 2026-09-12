# API Docs Bukuin

Dokumentasi API ini berfungsi sebagai panduan pengembangan dan testing untuk aplikasi Bukuin.

## Base URL

```text
http://localhost:8080
```

## Authentication

Beberapa endpoint membutuhkan session login. Untuk request yang memerlukan login, lakukan login terlebih dahulu dan simpan cookie/session dari server.

## Endpoints

### Auth

#### 1. Register

- Method: `POST`
- Endpoint: `/api/auth/register`
- Body:

```json
{
  "name": "User Test",
  "email": "test@example.com",
  "password": "password123"
}
```

#### 2. Login

- Method: `POST`
- Endpoint: `/api/auth/login`
- Body:

```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

#### 3. Logout

- Method: `POST`
- Endpoint: `/api/auth/logout`

#### 4. Get Current User

- Method: `GET`
- Endpoint: `/api/auth/me`

---

### Books

#### 1. Get All Books

- Method: `GET`
- Endpoint: `/api/books`

#### 2. Get Book By ID

- Method: `GET`
- Endpoint: `/api/books/{id}`

---

### Cart

#### 1. Get Cart

- Method: `GET`
- Endpoint: `/api/cart`

#### 2. Add Book to Cart

- Method: `POST`
- Endpoint: `/api/cart/items`
- Body:

```json
{
  "book_id": 1,
  "quantity": 1
}
```

#### 3. Update Cart Item

- Method: `PUT`
- Endpoint: `/api/cart/items/{id}`
- Body:

```json
{
  "quantity": 2
}
```

#### 4. Delete Cart Item

- Method: `DELETE`
- Endpoint: `/api/cart/items/{id}`

---

### Orders

#### 1. Checkout

- Method: `POST`
- Endpoint: `/api/orders/checkout`
- Auth required: `Yes`

Response contoh:

```json
{
  "id": 1,
  "user_id": 1,
  "total_amount": 150000,
  "status": "pending",
  "payment_url": "https://app.sandbox.midtrans.com/snap/v2/...",
  "items": [
    {
      "id": 1,
      "book_id": 1,
      "title": "Belajar Go",
      "author": "Author",
      "price": 150000,
      "quantity": 1,
      "subtotal": 150000
    }
  ]
}
```

#### 2. Get User Orders

- Method: `GET`
- Endpoint: `/api/orders`
- Auth required: `Yes`

#### 3. Get Order By ID

- Method: `GET`
- Endpoint: `/api/orders/{id}`
- Auth required: `Yes`

#### 4. Get User Invoice / Struk

- Method: `GET`
- Endpoint: `/api/orders/{id}/invoice`
- Auth required: `Yes`

#### 5. Midtrans Notification

- Method: `POST`
- Endpoint: `/api/orders/notification`

Catatan:
- Endpoint ini dipakai oleh Midtrans untuk callback status pembayaran.
- Ketika status berubah menjadi `paid`, order akan otomatis unlock untuk My Library.

---

### Library

#### 1. Get My Library

- Method: `GET`
- Endpoint: `/api/library`
- Auth required: `Yes`

Response contoh:

```json
[
  {
    "id": 1,
    "order_id": 1,
    "book_id": 1,
    "title": "Belajar Go",
    "author": "Author",
    "category": "Programming",
    "cover_url": "/static/uploads/cover.jpg",
    "file_path": "/static/uploads/books/sample.pdf",
    "paid_at": "2026-09-12T12:00:00Z"
  }
]
```

---

### Admin

#### 1. Get Admin Stats

- Method: `GET`
- Endpoint: `/api/admin/stats`
- Auth required: `Admin`

#### 2. Get All Orders for Admin

- Method: `GET`
- Endpoint: `/api/admin/orders`
- Auth required: `Admin`

#### 3. Get Invoice for Admin

- Method: `GET`
- Endpoint: `/api/admin/orders/{id}/invoice`
- Auth required: `Admin`

#### 4. Update Order Status

- Method: `PUT`
- Endpoint: `/api/admin/orders/{id}/status`
- Auth required: `Admin`
- Body:

```json
{
  "status": "paid"
}
```

#### 5. Generate Report

- Method: `GET`
- Endpoint: `/api/admin/report`
- Auth required: `Admin`

#### 6. Create Book

- Method: `POST`
- Endpoint: `/api/admin/books`
- Auth required: `Admin`
- Form-data:
  - `title`
  - `author`
  - `category`
  - `description`
  - `price`
  - `cover` (optional)
  - `ebook` (optional)

#### 7. Update Book

- Method: `PUT`
- Endpoint: `/api/admin/books/{id}`
- Auth required: `Admin`
- Form-data sama seperti create book

#### 8. Delete Book

- Method: `DELETE`
- Endpoint: `/api/admin/books/{id}`
- Auth required: `Admin`

---

## Response Codes

- `200 OK`
- `201 Created`
- `400 Bad Request`
- `401 Unauthorized`
- `403 Forbidden`
- `404 Not Found`
- `500 Internal Server Error`

## Tips Menggunakan Bruno

1. Buat collection `Bukuin API`.
2. Tambahkan environment variable `baseUrl`.
3. Simpan request untuk setiap endpoint.
4. Gunakan request login terlebih dahulu jika endpoint membutuhkan autentikasi.
5. Simpan response untuk dokumentasi contoh request/response.

## Catatan

- Untuk local development, gunakan `http://localhost:8080`.
- Untuk production, ganti `baseUrl` dengan domain yang aktif.
- `APP_BASE_URL` juga dapat digunakan agar redirect Midtrans mengarah ke URL yang benar.
