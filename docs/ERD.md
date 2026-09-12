# ERD Bukuin

Berikut adalah diagram ERD untuk project Bukuin menggunakan Mermaid.

```mermaid
erDiagram
    USERS ||--o{ ORDERS : makes
    USERS ||--o{ CARTS : owns
    USERS ||--o{ LIBRARY_ITEMS : has

    BOOKS ||--o{ CART_ITEMS : contains
    BOOKS ||--o{ ORDER_ITEMS : includes

    CARTS ||--o{ CART_ITEMS : contains
    ORDERS ||--o{ ORDER_ITEMS : consists_of

    USERS {
        int id PK
        string name
        string email
        string password_hash
        string role
    }

    BOOKS {
        int id PK
        string title
        string author
        string category
        string description
        decimal price
        string cover_url
        string file_path
    }

    CARTS {
        int id PK
        int user_id FK
    }

    CART_ITEMS {
        int id PK
        int cart_id FK
        int book_id FK
        int quantity
    }

    ORDERS {
        int id PK
        int user_id FK
        decimal total_amount
        string status
        datetime created_at
        string payment_method
        string payment_url
        string transaction_id
        datetime paid_at
    }

    ORDER_ITEMS {
        int id PK
        int order_id FK
        int book_id FK
        string title
        string author
        decimal price
        int quantity
        decimal subtotal
    }

    LIBRARY_ITEMS {
        int id PK
        int user_id FK
        int book_id FK
        int order_id FK
        datetime paid_at
    }
```

## Penjelasan singkat

- `users` menyimpan data pengguna dan role (admin/user).
- `books` menyimpan katalog buku, termasuk cover dan file ebook.
- `carts` dan `cart_items` menyimpan keranjang belanja user.
- `orders` menyimpan transaksi pesanan.
- `order_items` menyimpan detail item per order.
- `library_items` menyimpan kumpulan ebook yang sudah dibeli oleh user.

## Catatan

- Status order biasanya bernilai: `pending`, `paid`, `cancelled`, `refunded`.
- `library_items` dipakai untuk fitur My Library agar user bisa membuka ebook yang telah dibeli.
