# Aplikasi E-Commerce Mobile

Aplikasi e-commerce mobile dengan backend Go dan frontend React Native yang modern dan mudah digunakan.

## Fitur Utama

- 🔐 Autentikasi pengguna (login dan registrasi)
- 🛍️ Katalog produk dengan pencarian dan filter
- 🛒 Keranjang belanja
- 💳 Proses checkout dan pembayaran
- 📋 Riwayat dan detail pesanan
- 👨‍💼 Panel admin untuk pengelolaan produk dan pesanan

## Teknologi

### Frontend

- React Native dengan Expo
- React Navigation untuk navigasi
- Context API untuk state management
- Axios untuk HTTP requests
- AsyncStorage untuk penyimpanan lokal

### Backend

- Golang dengan Gin framework
- JWT untuk autentikasi
- GORM sebagai ORM
- MySQL sebagai database
- REST API

## Prasyarat

Sebelum menjalankan aplikasi ini, pastikan Anda telah menginstal:

- **Go** (v1.16+)
- **Node.js** (v14+) dan npm/yarn
- **MySQL** (v5.7+)
- **Git**

## Struktur Proyek

```
├── ecom-be/           # Backend Go
│   ├── config/        # Konfigurasi database dan env
│   ├── controllers/   # Handler logika bisnis
│   ├── middleware/    # Middleware (auth, dll)
│   ├── models/        # Model data
│   └── routes/        # Definisi API endpoint
│
└── ecomfe/            # Frontend React Native
    ├── src/
    │   ├── components/  # Komponen reusable
    │   ├── contexts/    # Context API
    │   ├── navigation/  # React Navigation
    │   ├── screens/     # Screen aplikasi
    │   ├── services/    # API Services
    │   ├── types/       # Type definitions
    │   └── utils/       # Utility functions
    └── assets/        # Gambar, font, dll
```

## Setup dan Instalasi

### 1. Clone Repository

```bash
git clone <repo-url>
cd ecommers-mobile
```

### 2. Setup Backend

#### Konfigurasi Database

1. Pastikan MySQL sudah terinstal dan berjalan
2. Buat database baru:

```sql
CREATE DATABASE ecommerce CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

3. Pastikan user MySQL memiliki akses ke database:

```sql
GRANT ALL PRIVILEGES ON ecommerce.* TO 'root'@'localhost';
FLUSH PRIVILEGES;
```

4. Sesuaikan file `.env` jika diperlukan:

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=ecommerce
JWT_SECRET=my-super-secret-jwt-token-for-ecommerce-app
PORT=8080
```

#### Instal Dependensi dan Jalankan Backend

```bash
cd ecom-be
go mod tidy
go run main.go
```

Backend akan berjalan di `http://localhost:8080`

### 3. Setup Frontend

```bash
cd ecomfe
npm install
# atau
yarn install
```

### 4. Jalankan Frontend

```bash
npx expo start
```

## Migrasi Database

Saat pertama kali menjalankan aplikasi, tabel-tabel akan otomatis dibuat oleh GORM. Jika Anda ingin menambahkan data awal untuk testing, Anda dapat mengeksekusi SQL berikut:

```sql
-- Tambahkan user admin (password: password)
INSERT INTO users (name, email, password, role, created_at)
VALUES ('Admin', 'admin@example.com', '$2a$10$NENgSd8F7zCoJuxU6dZLs.CpLXKnN08VuFkGhjYqFJ9xI5i74gKH6', 'admin', NOW());

-- Tambahkan user biasa (password: password)
INSERT INTO users (name, email, password, role, created_at)
VALUES ('User', 'user@example.com', '$2a$10$NENgSd8F7zCoJuxU6dZLs.CpLXKnN08VuFkGhjYqFJ9xI5i74gKH6', 'user', NOW());

-- Tambahkan produk contoh
INSERT INTO products (name, description, price, stock, created_at)
VALUES
('Smartphone XYZ', 'Smartphone canggih dengan fitur terbaru', 2500000, 50, NOW()),
('Laptop ABC', 'Laptop ringan dengan performa tinggi', 8000000, 20, NOW()),
('Headphone Premium', 'Headphone dengan kualitas suara terbaik', 1200000, 100, NOW());
```

## Endpoint API

### Autentikasi

- `POST /register` - Registrasi user baru
- `POST /login` - Login user

### Produk (Publik)

- `GET /products` - Daftar semua produk
- `GET /products/:id` - Detail produk

### Keranjang (Perlu Autentikasi)

- `GET /api/cart` - Lihat keranjang
- `POST /api/cart` - Tambah produk ke keranjang
- `PUT /api/cart/:id` - Update item keranjang
- `DELETE /api/cart/:id` - Hapus item dari keranjang
- `DELETE /api/cart` - Kosongkan keranjang

### Pesanan (Perlu Autentikasi)

- `POST /api/orders` - Buat pesanan baru
- `GET /api/orders` - Daftar pesanan
- `GET /api/orders/:id` - Detail pesanan
- `PUT /api/orders/:id/cancel` - Batalkan pesanan

### Admin (Perlu Role Admin)

- `POST /admin/products` - Tambah produk baru
- `PUT /admin/products/:id` - Update produk
- `DELETE /admin/products/:id` - Hapus produk
- `GET /admin/orders` - Daftar semua pesanan
- `PUT /admin/orders/:id/status` - Update status pesanan

## Kredensial Default

### Pengguna

- Email: user@example.com
- Password: password

### Admin

- Email: admin@example.com
- Password: password

## Troubleshooting

### Masalah Database

- **Error koneksi MySQL**: Pastikan MySQL berjalan dan kredensial di file `.env` benar
- **Error migrasi**: Periksa bahwa database kosong atau tidak memiliki konflik dengan skema yang ada
- **Error "unknown driver"**: Pastikan `go mod tidy` sudah dijalankan untuk menginstal semua dependensi

### Masalah Frontend

- **Error module not found**: Pastikan semua dependensi sudah terinstal dengan `npm install`
- **Error koneksi ke API**: Pastikan backend berjalan dan URL API benar di `src/utils/config.ts`
- **Masalah emulator**: Jika menggunakan emulator Android, gunakan `10.0.2.2:8080` sebagai host API

## Lisensi

[MIT License](LICENSE)

## Kontributor

- [Nama Kontributor](https://github.com/username)

## Kontak

Untuk pertanyaan atau saran, hubungi kami di [email@example.com](mailto:email@example.com)
