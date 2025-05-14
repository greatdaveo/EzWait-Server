# EZWait API (Golang Backend)

EZWait is a mobile-first waitlist and appointment management system built for service providers such as barbers, hair stylist, salons, and similar service-based businesses. This backend is written in Go using the Fiber framework and connects to a PostgreSQL database.

It powers user registration, booking appointments, managing profiles (stylist and customer), reminders, and more.

---

## Features

### Authentication & Users Profile Management
- Register/Login using JWT tokens
- Supports 2 roles: `stylist` and `customer`
- Update profile (name, email, location, image)
- Change and reset password
- Delete account
- Toggle appointment reminder setting

### Stylists
- Create/update stylist profile
- Add services, service images, and available time slots
- Track total bookings and current customers
- Set auto-confirmation for new bookings
- List or filter stylists by service or rating

### Bookings
- Create booking with `start_time`, `end_time`, and `booking_day`
- Customers and stylists can view their bookings
- View a single booking’s details
- Filter bookings by status
- Cron job to auto-mark past bookings as completed

### Notifications (WIP)
- Push notification toggle (`isReminderOn`)
- Email verification & OTP logic coming soon

---

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: [Fiber](https://gofiber.io)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Hosting**: Render
- **Push**: Expo (client-side + backend support planned)
- **Email (planned)**: Resend/SMTP for OTPs

---

## Project Structure


---

## API Endpoints

### Auth
| Method | Endpoint              | Description         |
|--------|-----------------------|---------------------|
| POST   | `/api/auth/register`  | Register new user   |
| POST   | `/api/auth/login`     | Login user (JWT)    |

### Users
| Method | Endpoint                      | Description                 |
|--------|-------------------------------|-----------------------------|
| GET    | `/api/user/me`                | Get current user profile    |
| PUT    | `/api/user/update`            | Update profile              |
| PUT    | `/api/user/change-password`   | Change password             |
| DELETE | `/api/user/delete`            | Delete account              |

### Stylists
| Method | Endpoint                  | Description                     |
|--------|---------------------------|---------------------------------|
| POST   | `/api/stylist/profile`    | Create stylist profile          |
| GET    | `/api/stylist/all`        | List stylists (with filters)    |
| GET    | `/api/stylist/:id`        | Get single stylist profile      |

### Bookings
| Method | Endpoint                   | Description                      |
|--------|----------------------------|----------------------------------|
| POST   | `/api/booking`             | Create new booking               |
| GET    | `/api/bookings`           | View all bookings (filtered)     |
| GET    | `/api/bookings/:id`       | View single booking              |

---

## .env Configuration

```env
APP_PORT=3000
DB_HOST=your_host
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db
DB_PORT=5432
JWT_SECRET=your_jwt_secret

# Clone the repo
git clone https://github.com/greatdaveo/EzWait-Server
cd EzWait-Server

# Set up environment
cp .env.example .env

# Install dependencies
go mod tidy

# Run the server (hot reload)
air

services:
  - type: web
    name: ezwait-api
    runtime: go
    plan: free
    buildCommand: go build -tags netgo -ldflags="-s -w" -o app ./cmd/server
    startCommand: ./app
    envVars:
      - key: DB_HOST
        sync: false
      - key: DB_USER
        sync: false
      - key: DB_PASSWORD
        sync: false
      - key: DB_NAME
        sync: false
      - key: DB_PORT
        value: "5432"
      - key: JWT_SECRET
        sync: false
```
### Future Enhancements
- Email verification (OTP)
- Push notification logic from backend
- Booking reminders with background scheduling
- Admin dashboard
- Analytics and reports

---

## 🧑‍💻 Author
Made with ❤️ by David Olowomeye – [GitHub](https://github.com/greatdaveo)