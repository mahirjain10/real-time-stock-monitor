# 📈 Stock Alert App

A real-time stock monitoring and alert application built with Go backend and Flutter frontend. Get instant notifications when your stocks hit target prices!

## 🚀 Features

### Backend (Go) - ✅ **FULLY IMPLEMENTED**
- **🔐 User Authentication System**: Complete registration and login with bcrypt password hashing
- **📊 Real-time Stock Monitoring**: WebSocket-based live price updates every 2 seconds
- **🔔 Advanced Alert System**: Create, update, delete, and toggle stock alerts with multiple conditions
- **⚡ Redis Pub/Sub Integration**: Multi-topic messaging system for instant notifications
- **💾 Dual Database Storage**: MySQL for persistence + Redis for caching and real-time data
- **🌐 RESTful API**: 7 fully functional endpoints with proper error handling
- **📈 External Stock API Integration**: Real-time price fetching from external stock data provider
- **🎯 Smart Condition Evaluation**: Support for >, >=, ==, <, <= operators
- **🔄 Automatic Alert Processing**: Background monitoring with automatic trigger notifications
- **🛡️ Input Validation**: Comprehensive request validation and sanitization

### Frontend (Flutter) - 🔄 **IN DEVELOPMENT**
- **🎨 Modern Material Design 3**: Clean, professional UI with custom color scheme
- **📱 Responsive Design**: Adaptive layout for different screen sizes
- **🔔 Alert Creation Form**: Complete form with 4 input fields (Alert Name, Stock Name, Current Price, Alert Price)
- **⚙️ Condition Selector**: Interactive segmented slider for choosing alert conditions (>, >=, ==, <, <=)
- **📋 Alert History Screen**: Search functionality for managing existing alerts
- **🧭 Bottom Navigation**: Tab-based navigation between Alert Screen and History
- **🎭 Custom Typography**: Poppins and WorkSans fonts for enhanced readability
- **🖼️ Asset Integration**: Custom icons and images (bell, history, notification icons)
- **📱 Cross-platform Ready**: Configured for Android, iOS, Web, Windows, macOS, and Linux

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Flutter App   │◄──►│   Go Backend    │◄──►│   MySQL DB      │
│   (Frontend)    │    │   (REST API)    │    │   (Data Store)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   Redis Cache   │
                       │   (Pub/Sub)     │
                       └─────────────────┘
```

## 🛠️ Tech Stack

### Backend
- **Go 1.22.4** - Core backend language
- **Gin** - HTTP web framework
- **MySQL** - Primary database
- **Redis** - Caching and pub/sub messaging
- **WebSocket** - Real-time communication
- **Gorilla WebSocket** - WebSocket implementation

### Frontend
- **Flutter 3.4.3+** - Cross-platform UI framework
- **Material Design 3** - Modern UI components
- **Custom Fonts** - Poppins and WorkSans typography
- **SVG Icons** - Scalable vector graphics

## 📋 Prerequisites

Before running the application, ensure you have:

- **Go 1.22.4+** installed
- **Flutter SDK 3.4.3+** installed
- **MySQL 8.0+** running
- **Redis 6.0+** running
- **Git** for version control

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd stock-alert-app-main
```

### 2. Backend Setup

#### Environment Configuration
Create a `.env` file in the `backend` directory:

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

#### Database Setup
```bash
cd backend
mysql -u root -p
```

```sql
CREATE DATABASE stock_alert_db;
USE stock_alert_db;
```

#### Install Dependencies and Run
```bash
cd backend
go mod download
go run main.go
```

The backend will start on `http://localhost:8080`

### 3. Frontend Setup

```bash
cd frontend
flutter pub get
flutter run
```

## 📚 API Documentation

### 🔐 Authentication Endpoints

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/api/auth/register` | Register a new user | `{"name": "string", "email": "string", "password": "string"}` |
| POST | `/api/auth/login` | User login | `{"email": "string", "password": "string"}` |

### 📊 Alert Management Endpoints

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/api/alert/get-current-price` | Get real-time stock price | `{"ticker_to_monitor": "string"}` |
| POST | `/api/alert/create-stock-alert` | Create a new stock alert | `{"user_id": "string", "alert_name": "string", "ticker_to_monitor": "string", "alert_condition": "string", "alert_price": "float64"}` |
| PUT | `/api/alert/update-stock-alert` | Update existing alert | `{"user_id": "string", "id": "string", "alert_name": "string", "alert_condition": "string", "alert_price": "float64"}` |
| PUT | `/api/alert/update-stock-alert-status` | Toggle alert active status | `{"user_id": "string", "id": "string", "active": "boolean"}` |
| DELETE | `/api/alert/delete-stock-alert` | Delete an alert | `{"user_id": "string", "id": "string"}` |
| POST | `/api/alert/alert-notification` | Send alert notification (internal) | `{"user_id": "string", "id": "string", "active": "boolean"}` |

### 🔌 WebSocket Endpoints

| Endpoint | Description | Message Format |
|----------|-------------|----------------|
| `/ws/get-stock-price-socket` | Real-time stock price updates | `{"ticker_to_monitor": "string", "alert_id": "string", "is_active": "boolean"}` |

### 📝 **API Response Format**
```json
{
  "statusCode": 200,
  "message": "Success message",
  "data": { /* response data */ },
  "error": null,
  "success": true
}
```

## 🗄️ Database Schema

### 👤 Users Table
```sql
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,                    -- UUID for user identification
    name VARCHAR(255) NOT NULL,                    -- User's full name
    email VARCHAR(255) UNIQUE NOT NULL,            -- Unique email address
    password VARCHAR(255) NOT NULL,                -- Bcrypt hashed password
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 🔔 Stock Alert Table
```sql
CREATE TABLE stock_alert (
    id VARCHAR(36) PRIMARY KEY,                    -- UUID for alert identification
    user_id VARCHAR(36) NOT NULL,                  -- Reference to users table
    ticker VARCHAR(20) NOT NULL,                   -- Stock symbol (e.g., "AAPL", "GOOGL")
    alert_name VARCHAR(255) NOT NULL,              -- Custom name for the alert
    current_fetched_price DECIMAL(10,2),           -- Latest fetched price
    current_fetched_time TIMESTAMP,                -- When price was last fetched
    alert_condition VARCHAR(5) NOT NULL,           -- Condition: >, >=, ==, <, <=
    alert_price DECIMAL(10,2) NOT NULL,            -- Target price for alert
    is_active BOOLEAN DEFAULT TRUE,                -- Alert status (active/inactive)
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 📊 Monitor Stock Table
```sql
CREATE TABLE monitor_stock (
    id VARCHAR(36) PRIMARY KEY,                    -- UUID for monitoring record
    alert_id VARCHAR(36) NOT NULL,                 -- Reference to stock_alert table
    ticker VARCHAR(20) NOT NULL,                   -- Stock symbol being monitored
    is_active BOOLEAN DEFAULT TRUE,                -- Monitoring status
    FOREIGN KEY (alert_id) REFERENCES stock_alert(id) ON DELETE CASCADE
);
```

### 🔄 Redis Data Structure
```
Key: alert_id (e.g., "123e4567-e89b-12d3-a456-426614174000")
Hash Fields:
- user_id: "user-uuid"
- ticker: "AAPL"
- alert_price: "150.50"
- alert_condition: ">"
- active: "true"

Key: monitor_stock:monitor_id
Hash Fields:
- id: "monitor-uuid"
- alert_id: "alert-uuid"
- ticker: "AAPL"
- is_active: "true"
```

## 🔧 Development

### 🚀 Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Create .env file with your configuration
cp .env.example .env

# Run the application
go run main.go

# Alternative: Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...

# Build for production
go build -o stock-alert-backend main.go
```

### 📱 Frontend Development

```bash
cd frontend

# Install dependencies
flutter pub get

# Run in debug mode (choose your platform)
flutter run                    # Interactive platform selection
flutter run -d chrome         # Run on web browser
flutter run -d android        # Run on Android device/emulator
flutter run -d ios            # Run on iOS simulator

# Build for production
flutter build apk --release           # Android APK
flutter build ios --release           # iOS (requires macOS)
flutter build web --release           # Web deployment
flutter build windows --release       # Windows executable
flutter build macos --release         # macOS application
flutter build linux --release         # Linux application
```

### 🔧 **Development Tools & Commands**

#### Backend Utilities
```bash
# Database migration
go run internal/models/init__tables_model.go

# Redis connection test
redis-cli ping

# API testing with curl
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'
```

#### Frontend Utilities
```bash
# Clean build cache
flutter clean
flutter pub get

# Analyze code quality
flutter analyze

# Run tests
flutter test

# Generate app icons
flutter packages pub run flutter_launcher_icons:main
```

### Project Structure

```
stock-alert-app-main/
├── backend/                 # Go backend application
│   ├── internal/           # Private application code
│   │   ├── app/           # Application services
│   │   ├── database/      # Database connections
│   │   ├── helpers/       # Utility functions
│   │   ├── models/        # Data models
│   │   ├── types/         # Type definitions
│   │   ├── utils/         # Helper utilities
│   │   ├── validator/     # Input validation
│   │   └── websocket/     # WebSocket handling
│   ├── web/               # HTTP handlers and routes
│   ├── main.go           # Application entry point
│   └── go.mod            # Go module definition
├── frontend/              # Flutter frontend application
│   ├── lib/              # Dart source code
│   │   ├── screens/      # UI screens
│   │   └── widgets/      # Reusable widgets
│   ├── assets/           # Images, fonts, icons
│   └── pubspec.yaml      # Flutter dependencies
└── README.md             # This file
```

## 🚧 Current Status

### ✅ **Backend - COMPLETELY IMPLEMENTED**
- **🔐 Authentication**: Full user registration/login with bcrypt hashing
- **📊 Stock Price API**: Real-time price fetching from external API every 2 seconds
- **🔔 Alert Management**: Complete CRUD operations for stock alerts
- **⚡ WebSocket Hub**: Multi-client real-time communication system
- **🔄 Redis Pub/Sub**: Dual-channel messaging (monitor + alert-topic)
- **💾 Database Models**: 3 tables (users, stock_alert, monitor_stock)
- **🎯 Condition Logic**: Smart price comparison with 5 operators
- **📡 API Endpoints**: 7 fully functional REST endpoints
- **🛡️ Error Handling**: Comprehensive validation and error responses
- **⚙️ Auto-Monitoring**: Background processes for continuous price tracking

### 🔄 **Frontend - IN DEVELOPMENT**
- **✅ UI Framework**: Complete Material Design 3 setup
- **✅ Navigation**: Bottom tab navigation with 2 screens
- **✅ Alert Form**: 4-input form with validation styling
- **✅ Condition Selector**: Interactive segmented button slider
- **✅ Custom Assets**: Icons, fonts, and images integrated
- **✅ Responsive Layout**: Adaptive design for all screen sizes
- **🔄 API Integration**: Backend connection in progress
- **🔄 WebSocket Connection**: Real-time updates integration pending
- **🔄 Alert History**: Data fetching and display logic pending
- **🔄 Form Submission**: Backend API calls pending

### 📋 **Planned Enhancements**
- 📧 Email notifications for triggered alerts
- 📱 Push notifications for mobile devices
- 📊 Interactive stock charts and graphs
- 💼 Portfolio tracking and management
- 🔍 Advanced search and filtering
- 📈 Historical price data visualization
- 🌍 Multi-language support
- 🏪 App store deployment (iOS/Android)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Mahir Jain** - *Initial work* - [mahirjain_10](https://github.com/mahirjain_10)

## 🙏 Acknowledgments

- Stock price data providers
- Open source contributors
- Flutter and Go communities

## 📞 Support

If you have any questions or need help, please:
- Open an issue on GitHub
- Contact the maintainer
- Check the documentation

## 🎯 **What You've Built - Technical Highlights**

### 🏗️ **Backend Architecture Excellence**
- **Microservices-Ready**: Clean separation of concerns with modular design
- **Real-time Engine**: WebSocket hub managing multiple concurrent connections
- **Smart Monitoring**: Automatic price fetching every 2 seconds with condition evaluation
- **Dual Storage Strategy**: MySQL for ACID compliance + Redis for lightning-fast caching
- **Pub/Sub Messaging**: Multi-topic Redis channels for scalable notifications
- **Production-Ready**: Comprehensive error handling, validation, and logging

### 🎨 **Frontend Design Innovation**
- **Material Design 3**: Latest Google design system implementation
- **Custom Components**: Segmented button slider for condition selection
- **Responsive UI**: Adaptive layouts for all screen sizes
- **Professional Typography**: Custom font integration (Poppins + WorkSans)
- **Asset Management**: Optimized images and icons for all platforms

### 🔧 **Advanced Technical Features**
- **UUID Generation**: Secure, collision-resistant ID system
- **Bcrypt Security**: Industry-standard password hashing
- **Context Management**: Proper Go context handling for cancellation
- **Connection Pooling**: Efficient database connection management
- **Graceful Shutdowns**: Proper resource cleanup on application exit

### 📊 **Real-time Data Flow**
```
External Stock API → Price Fetching → WebSocket Hub → Client Updates
                ↓
Redis Pub/Sub → Condition Evaluation → Alert Triggering → Notifications
```

---

**🚀 Status**: Backend is production-ready with comprehensive testing. Frontend UI is complete, API integration in progress.

**Happy Trading! 📈**
