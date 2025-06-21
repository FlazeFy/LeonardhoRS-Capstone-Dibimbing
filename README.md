# 🚀 PELITA

> **PELITA** is designed as an **inventory and asset management** system for offices or warehouses, allowing users to store all asset data, monitor assets, create routine maintenance schedules, and perform daily checklists. Additionally, the application is integrated with **Telegram Chat**, enabling users to receive broadcasts easily and monitor data with just a simple click in the Telegram chat room.

> This is my **Final Project Application** when I take 6-Month **Backend Development Bootcamp** at **Dibimbing.id**, the courses itself specializing in JavaScript and **Go**. It also learning about ERD Design, JavaScript and Go fundamentals, OOP in Go, SQL (DDL, DML, DCL), NoSQL with MongoDB, Git Flow, REST API development, Postman testing, file handling, Docker deployment, microservices, and professional ethics training.

## 📋 Basic Information

If you want to see the project detail documentation, you can read my software documentation document. 

1. **Pitch Deck**
https://docs.google.com/presentation/d/14qVLPVMjAUH41Lu2ywoHwUac-U5Uh-wPo_noPZ2D-5k/edit?usp=sharing 

2. **Diagrams**
https://drive.google.com/file/d/1SMg47iUcwOPuq941d1vfzisYLV2CVHsn/view?usp=drive_link 

3. **Software Requirement Specification**
https://docs.google.com/document/d/1YrI2zIeso2RFS84KFrrcj1FxbeeBgmRild-x5L9fr84/edit?usp=sharing 

4. **User Manual**
https://docs.google.com/presentation/d/1Opa2PSrKvuIC9vVz8GkdJSfeyZrwzvT88W78NUWCe1Y/edit?usp=sharing

5. **Presentation Slide**
https://docs.google.com/presentation/d/1Opa2PSrKvuIC9vVz8GkdJSfeyZrwzvT88W78NUWCe1Y/edit?usp=sharing 

6. **Test Cases**
https://docs.google.com/spreadsheets/d/1UZ97iSDHLe6XYc5qb1qnLwP3TCwy98211GIggZ0KX4w/edit?usp=sharing 

### 🌐 Deployment URL

Backend (Swagger Docs) : https://pelita.leonardhors.com/swagger/index.html#/
Telegram : https://t.me/pelita_app_bot

### 🎥 Presentation Materials (PDF)

https://docs.google.com/presentation/d/1c4EuNaVFpk-2JuQReHT9npc9y-ClZuMqEj78EG08LEE/edit?usp=drive_link 

### 📱 Demo video on the actual device

[URL]

---

## 🎯 Product Overview

- **Asset and Inventory Management**
Users can add, update, categorize, and track office or warehouse assets and inventory in real-time. Each asset record includes details like location, condition, purchase date, and assigned personnel, making it easy to organize and audit resources.

- **Maintenance Scheduling**
PELITA allows users to set up routine maintenance schedules for equipment or assets. Notifications and reminders help ensure timely inspections and servicing, reducing the risk of equipment failure and prolonging asset lifespan.

- **Daily Checklist and Monitoring**
The system supports customizable daily checklists for different departments or asset types. Staff can log their checks directly in the system, helping ensure operational standards are met and providing a digital trail of accountability.

- **Telegram Bot Chat Integration**
This app is integrated with Telegram to provide instant access to system updates, maintenance alerts, and asset status. Users can receive broadcasts or even interact with the system via Telegram, simplifying communication and monitoring outside the app.

- **PDF Report Export and Auto Delivery**
PELITA can automatically generate and export comprehensive asset and activity reports in PDF format. These documents are sent directly to designated users, streamlining reporting processes and supporting audit readiness and decision-making.

## 🚀 Target Users

1. **Office Managers / Warehouse Managers (Admin Role)**
Personel responsible for manage asset use, inventory levels, and maintenance activities.

2. **Maintenance Staff / Technicians (Technician Role)**
Staff responsible for checking, repairing, or maintaining assets regularly.

3. **Executives or Department Heads (Guest Role)**
Higher-level stakeholders who need quick overviews of company assets.

3. **Office Staff & Office Guest (Guest Role)**
Other employees who require quick access to facility overviews or wish to submit feedback regarding assets they come across.

## 🧠 Problem to Solve

1. Offices and warehouses struggle to efficiently manage and track assets, leading to disorganized inventory, misplaced items, and time-consuming audits.
2. Routine maintenance is often forgotten or delayed, resulting in equipment failure, higher repair costs, and reduced asset lifespan.
3. Communication about asset status or maintenance updates is fragmented, causing delays and lack of visibility among teams.
4. Compiling reports for audits, management reviews, or compliance is time-consuming and prone to human error.

## 💡 Solution

1. Provide a feature that allows users to **add, update, categorize, and track assets**. Complete with location, condition, purchase date, and assigned personnel making inventory well-organized and easy to audit.
2. Enable admins to set **automated maintenance schedules** with reminders and technician assignments, ensuring timely inspections and prolonging asset durability.
3. Integrate with **Telegram to deliver real-time updates**, maintenance alerts, and allow users to monitor or respond to asset activities from within their chat room.
4. Allow users to automatically gather asset and inventory data, export it into a well-formatted **PDF report**, and send it directly to designated users—saving time and ensuring accuracy in reporting.

## 🔗 Features

- ✍️ Asset and Inventory Management
- ⏰ Maintenance Scheduling
- ✅ Daily Checklist and Monitoring
- 🤖 Telegram Bot Chat Integration
- 📄 PDF Report Export and Auto Delivery

---

## 🛠️ Tech Stack

### Backend

- Golang Gin
- Golang - Telegram Bot


### Database

- MySQL

### Others Data Storage

- Firebase Storage (Cloud Storage for Asset File)
- Redis (In-Memory Storage for Sign Out Schema)

### Infrastructure & Deployment

- Cpanel (Deployment)
- Github (Code Repository)
- Firebase (External Services)
- Docker (Deployment)

### Other Tools & APIs

- Postman
- Swagger Docs

---

## 🏗️ Architecture
### Structure

### 📁 Project Structure

| Directory/File      | Purpose                                                                 |
|---------------------|-------------------------------------------------------------------------|
| `config/`           | Application configuration files (e.g., environment, database, externa services, auth, and in-memory storage) and const variabels           |
| `controller/`       | Handles incoming HTTP requests and sends responses                      |
| `docs/`             | Swagger API documentation setup                                         |
| `entity/`           | Core data structures and model definitions                              |
| `factory/`          | Generates dummy/test data for development or testing                    |
| `middleware/`       | Custom Gin middleware (e.g., authentication, logging, role access management)                   |
| `reports/`          | PDF generation and reporting logic                                      |
| `repository/`       | Data access logic (repository pattern for DB abstraction)               |
| `routes/`           | Defines dependency, API routes and maps them to controllers, also scheduler and seeder init                         |
| `scheduler/`        | Background jobs or scheduled tasks (e.g., cron-like functions)          |
| `seeder/`           | Database seeding logic for initial or test data                         |
| `service/`          | Business logic layer reused across controllers                          |
| `tests/`            | End-to-End, Unit and integration tests                                              |
| `utils/`            | Utility and helper functions                                            |
| `.env`              | Environment variable configuration                                      |
| `.gitignore`        | Specifies intentionally untracked files to ignore by Git                |
| `go.mod`            | Go module definition (dependencies and module path)                     |
| `go.sum`            | Checksums for module dependencies                                       |
| `main.go`           | Entry point of the application                                          |
| `README.md`         | Project simple documentation                                                   |

---

### 🧾 Environment Variables

To set up the environment variables, just create the `.env` file in the root level directory.

| Variable Name                     | Description                                                    |
|----------------------------------|----------------------------------------------------------------|
| `DB_HOST`                        | Database host (e.g., `localhost`)                              |
| `DB_PORT`                        | Database port (e.g. `3306`)                            |
| `DB_USER`                        | Database username                                              |
| `DB_PASSWORD`                    | Database password                                              |
| `DB_NAME`                        | Name of the primary database                                   |
| `TEST_DB_HOST`                   | Host for the test database                                     |
| `TEST_DB_PORT`                   | Port for the test database                                     |
| `TEST_DB_USER`                   | Username for the test database                                 |
| `TEST_DB_PASSWORD`               | Password for the test database                                 |
| `TEST_DB_NAME`                   | Name of the test database                                      |
| `JWT_SECRET_KEY`                 | Secret key used for JWT authentication                         |
| `JWT_EXPIRES_IN`                 | JWT token expiration duration (e.g., `1h`, `24h`)              |
| `FIREBASE_BUCKET_NAME`           | Firebase Storage bucket name for handling file uploads         |
| `GOOGLE_APPLICATION_CREDENTIALS`| Path to Firebase service account JSON file                     |
| `TELEGRAM_BOT_TOKEN`             | Telegram bot token for chat integration                        |
| `PORT`                           | Port on which the application will run (e.g., `9000`)          |

---

## 🗓️ Development Process

### Technical Challenges

- **Daily Limitation** for data transaction in Firebase Storage
- Not all **utils (helpers)** can be tested in **automation testing**
- Feature that return the **output in Telegram Chat or Exported File** must be **tested manually** 

---

## 🚀 Setup & Installation

### Prerequisites

#### 🔧 General
- Git installed
- A GitHub account
- Basic knowledge of Golang, Software Testing, Firebase Service,  and SQL Databases
- Code Editor
- Telegram Account
- Postman

#### 🧠 Backend
- Go version 1.21 or higher
- Git for cloning the repository.
- MySQL database.
- Make (optional), if your project includes a Makefile to simplify common commands.
- Swagger CLI to generate Swagger API docs.
- Firebase service account JSON file or Google App Credential.
- Telegram Bot token, you can get it from **Bot Father** `@BotFather`
- Telegram User ID for testing the scheduler chat in your Telegram Account. You can get it from **IDBot** `@username_to_id_bot`
- Internet access from the hosting server (for Telegram webhook polling or long-polling)

### Installation Steps

**Local Init & Run**
1. Download this Codebase as ZIP or Clone to your Git
2. Set Up Environment Variables `.env` at the root level directory. You can see all the variable name to prepare at the **Project Structure** before or `.env.example`
3. Install Dependencies using `go mod tidy`
4. DB Migration will be running everytime your start the app. If you want to disabled it, just commented the code `MigrateAll(db)` that labelled with `Connect DB`, you can find the file at `main.go`. 
5. Seeder also will be running everytime your start the app. If you want to disabled it, just commented the code that labelled with `Seeder & Factories`, you can find the file at `routes/dependency.go`. 
6. Same like Task Scheduler. If you want to disabled it, just commented the code that labelled with `Task Scheduler`, you can find the file at `routes/dependency.go`. 
7. Run the Go Gin using `go run main.go`

**CPanel Deployment**
1. Source code uploaded to CPanel
2. ...

---

## 👥 Team Information

| Role     | Name                    | GitHub                                     | Responsibility |
| -------- | ----------------------- | ------------------------------------------ | -------------- |
| Backend Developer  | Leonardho R. Sitanggang | [@FlazeFy](https://github.com/FlazeFy)     | Manage Backend and Telegram Bot Codebase         |
| System Analyst  | Leonardho R. Sitanggang | [@FlazeFy](https://github.com/FlazeFy)     | Manage Diagram & Software Docs         |
| Quality Assurance  | Leonardho R. Sitanggang | [@FlazeFy](https://github.com/FlazeFy)     | Manage Testing & Documented The API         |

---

## 📝 Notes & Limitations

### ⚠️ Precautions When Using the Service
- Ensure API endpoints requiring authentication are protected with proper middleware.
- Do not expose sensitive environment variables (e.g., API keys, database credentials) in public repositories.
- Avoid using seeded dummy data in production environments.
- Avoid using seeded dummy data with large seed at the same time.

### 🧱 Technical Limitations
- Telegram bot polling may cause delays or downtime if the server experiences high load

### 🐞 Known Issues
- Limitation when using Firebase Storage for free plan Firebase Service, upgrade to Blaze Plan to use more.

---

## 🏆 Appeal Points

- ✅ **Asset Tracking**: Easy to management all asset data across company.
- 🤖 **Bot Interaction**: Realtime reminder notifications and audit via Telegram Chat.
- 🧪 **Test-Driven**: All APIs are tested using End-to-End (E2E) testing, repositories are tested using integration tests, and utilities/helpers are tested using unit tests.
- 📄 **Auto-Generated PDF Reports**: Automatically compile and send structured inventory reports.
- 💡 **Practical Use Case**: Solves a real-world problem, making it relevant and valuable.
- ⭐️ **Clean and Modular Code Structure**: Easy to extend, test, and maintain.

---

_Made with ❤️ by Leonardho R. Sitanggang_