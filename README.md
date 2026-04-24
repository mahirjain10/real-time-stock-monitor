# Real-Time Stock Monitor

A stock alert app that actually tells you when to pay attention. Set a target price, pick your condition, and let it run in the background — you'll get a push notification the moment your stock hits the mark.

Built with a Go backend and Flutter frontend, communicating over **Server-Sent Events (SSE)** so prices stream live without hammering your battery or bandwidth.

---

## What it does

You set an alert like *"tell me when AAPL goes above $200"*, and the app:

1. Connects to a live stock feed
2. Checks the price every 2 seconds via SSE
3. Evaluates your condition (`>`, `>=`, `==`, `<`, `<=`)
4. Fires a Firebase push notification the moment it triggers

You can have multiple alerts running at once, pause them, update the target, or delete them entirely — all through the Flutter app.

---

## Architecture

```
Flutter App  ←──SSE stream──  Go Backend  ──►  MySQL
                                   │
                              Redis Pub/Sub
                                   │
                           External Stock API
                                   │
                         Firebase Cloud Messaging
```

The backend is the brain. It fetches live prices, evaluates alert conditions in the background, manages SSE connections per ticker, and pushes notifications when conditions are met. MySQL handles persistence; Redis handles caching and real-time pub/sub between goroutines.

---

## Tech Stack

**Backend**
- Go 1.22.4 + Gin (HTTP framework)
- MySQL 8.0 (persistence)
- Redis 6.0 (caching + pub/sub)
- Server-Sent Events (real-time price streaming)
- Firebase Cloud Messaging (push notifications)
- JWT (authentication)

**Frontend**
- Flutter 3.4.3+
- Material Design 3
- Poppins + WorkSans fonts

---

## Prerequisites

- Go 1.22.4+
- Flutter SDK 3.4.3+
- MySQL 8.0+
- Redis 6.0+
- A Firebase project (for push notifications)

---

## Getting Started

### 1. Clone

```bash
git clone https://github.com/mahirjain10/real-time-stock-monitor.git
cd real-time-stock-monitor
```

### 2. Backend

Create `backend/.env`:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_NAME=stock_alert_db

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=your_jwt_secret_key
```

Set up the database:

```bash
mysql -u root -p
```

```sql
CREATE DATABASE stock_alert_db;
USE stock_alert_db;
```

Run the backend:

```bash
cd backend
go mod download
go run main.go
```

Server starts at `http://localhost:8000`.

### 3. Frontend

```bash
cd frontend
flutter pub get
flutter run
```

---

## API Reference

### Auth

| Method | Endpoint | Body |
|--------|----------|------|
| POST | `/api/auth/register` | `{ name, email, password }` |
| POST | `/api/auth/login` | `{ email, password }` |

### Alerts

| Method | Endpoint | Body |
|--------|----------|------|
| POST | `/api/alert/get-current-price` | `{ ticker_to_monitor }` |
| POST | `/api/alert/create-stock-alert` | `{ user_id, alert_name, ticker_to_monitor, alert_condition, alert_price }` |
| PUT | `/api/alert/update-stock-alert` | `{ user_id, id, alert_name, alert_condition, alert_price }` |
| PUT | `/api/alert/update-stock-alert-status` | `{ user_id, id, active }` |
| DELETE | `/api/alert/delete-stock-alert` | `{ user_id, id }` |

### SSE

| Endpoint | Description | Query Params |
|----------|-------------|--------------|
| `GET /events` | Start streaming live price for a ticker | `alertID`, `ticker` |
| `GET /disconnect` | Stop monitoring and close SSE connection | `alertID`, `ticker` |

**SSE event format:**
```
data: {"price": 194.50, "time": "2026-04-24T09:15:00Z"}
```

**Standard API response format:**
```json
{
  "statusCode": 200,
  "message": "Alert created",
  "data": {},
  "error": null,
  "success": true
}
```

---

## Database Schema

```sql
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,         -- bcrypt hashed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE stock_alert (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    ticker VARCHAR(20) NOT NULL,
    alert_name VARCHAR(255) NOT NULL,
    current_fetched_price DECIMAL(10,2),
    current_fetched_time TIMESTAMP,
    alert_condition VARCHAR(5) NOT NULL,    -- >, >=, ==, <, <=
    alert_price DECIMAL(10,2) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE monitor_stock (
    id VARCHAR(36) PRIMARY KEY,
    alert_id VARCHAR(36) NOT NULL,
    ticker VARCHAR(20) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (alert_id) REFERENCES stock_alert(id) ON DELETE CASCADE
);
```

**Redis structure:**
```
Key: {alert_id}
Fields: user_id, ticker, alert_price, alert_condition, active

Key: monitor_stock:{monitor_id}
Fields: id, alert_id, ticker, is_active
```

---

## Project Structure

```
real-time-stock-monitor/
├── backend/
│   ├── internal/
│   │   ├── app/          # Application services
│   │   ├── database/     # DB connections
│   │   ├── events/       # SSE event handling
│   │   ├── helpers/      # Utility functions
│   │   ├── middleware/   # Auth middleware (JWT)
│   │   ├── models/       # Data models
│   │   ├── sse/          # SSE server and client management
│   │   ├── types/        # Type definitions
│   │   ├── utils/        # Helpers
│   │   └── validator/    # Input validation
│   ├── web/              # HTTP handlers and routes
│   ├── main.go
│   └── go.mod
├── frontend/
│   ├── lib/
│   │   ├── screens/      # Alert form, history screen
│   │   └── widgets/      # Reusable components
│   ├── assets/           # Icons, fonts, images
│   └── pubspec.yaml
└── README.md
```

---

## Development Commands

**Backend:**
```bash
# Hot reload (requires air)
go install github.com/air-verse/air@latest
air

# Run tests
go test ./...

# Test SSE manually
curl -N "http://localhost:8000/events?alertID=test-uuid&ticker=AAPL"

# Build for production
go build -o stock-alert-backend main.go
```

**Frontend:**
```bash
flutter run -d chrome        # Web
flutter run -d android       # Android
flutter run -d ios           # iOS

flutter build apk --release  # Android APK
flutter build web --release  # Web
flutter analyze              # Lint
flutter test                 # Tests
```

---

## Current Status

**Backend** — fully implemented:
- User auth (register/login with bcrypt + JWT)
- Live stock price fetching via external API
- Full CRUD for alerts
- SSE hub managing concurrent connections per ticker
- Background monitoring goroutines with condition evaluation
- Redis pub/sub for multi-channel notifications
- Firebase Cloud Messaging integration

**Frontend** — UI complete, API integration in progress:
- ✅ Material Design 3 setup
- ✅ Bottom tab navigation
- ✅ Alert creation form
- ✅ Condition selector (segmented slider)
- ✅ Alert history screen with search
- 🔄 Backend API calls
- 🔄 SSE connection and live price updates
- 🔄 Alert history data fetching

---

## Why SSE instead of WebSockets?

For stock price feeds, SSE is a better fit than WebSockets. Price data flows one way — server to client — and SSE handles that natively. You get automatic reconnection built into the browser, HTTP/2 compatibility, and no custom connection management code. It also works cleanly through most corporate firewalls that block WebSocket upgrades.

WebSockets make sense when you need two-way communication. For a monitoring feed where the server is doing all the talking, SSE keeps things simpler.

---

## Planned

- Email notifications when alerts trigger
- Interactive candlestick charts
- Portfolio tracking across multiple stocks
- Historical price data and alert analytics
- App store release (iOS + Android)

---

## Contributing

1. Fork the repo
2. Create a branch: `git checkout -b feature/your-feature`
3. Commit: `git commit -m 'Add your feature'`
4. Push: `git push origin feature/your-feature`
5. Open a pull request

---

## License

MIT — see [LICENSE](LICENSE) for details.

---

Built by [Mahir Jain](https://github.com/mahirjain10)
