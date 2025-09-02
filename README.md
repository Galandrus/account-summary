# Account Summary

Account summary system that processes CSV files and generates email reports.

## Assumptions

- Transactions belong to an **Account**. If the account does not exist, it is created with the provided email when transactions are loaded.
- Transaction dates use the format **YYYY-MM-DD**. I assumed the given example (mm/dd) was just illustrative. If it is actually a requirement, it could be easily formatted while parsing the file.
- Each transaction includes a **name field** for reference. This could eventually be used for classification.
- The home page exists solely to make application testing more user-friendly.
- The project includes only one **CSV file** containing 100 transactions.
  This file could be retrieved from **S3** by implementing the `FileLoaderInterface`, which should locate the `FILE PATH` in S3 and return an `*os.File`. After that, the CSV reader would work normally.

## Improvement Points

- Validate that loaded transactions are **unique** (e.g., by checking their ID) to prevent duplicates.
- Add more **tests** for both services and packages.
- Implement a **structured logger** to improve error readability.
- The API can be built with a web framework such as Fiber or Chi if more complex features are needed.

## Notes

- I used **MongoDB**, since it’s my daily driver and easier for me to work with. However, I understand that for handling transactions, a **relational database** is more suitable due to its stronger guarantees on data integrity and persistence.
  Migrating to a relational DB (e.g., MySQL) would be straightforward by re-implementing the `TransactionRepositoryInterface` and `AccountRepositoryInterface`.
- I implemented an **API** to expose the different features, as this is what I am most familiar with. That said, it could easily be migrated to a **Lambda function**, since the code is modular and just requires proper orchestration of components.
- I only wrote tests for the most important part: the process that analyzes transactions and generates the summary.
- The project is deployed on **Render** and can be accessed at:
  `https://account-summary.onrender.com/`
  (Note: the first load may take some time.)
- The MongoDB database is deployed on **MongoDB Atlas**.
- Both the **Home page** (`frontend.html`) and the **summary page** were generated entirely by an **LLM**, since the focus of the challenge was more on the Backend than the Frontend.
- Emails are sent from an existing Gmail account: `tepegamma@gmail.com`. I prefer using this backup account instead of my main one.
- I did not use AWS services because I do not have an active account and encountered issues when trying to create one.
- I built the API in the most basic way. Since I only needed a simple API, I decided not to use a web framework.

## Page Functionality

You can visit the deployed page at :
[https://account-summary.onrender.com/](https://account-summary.onrender.com/)

(Note: the first load may take some time.)

- **LOAD TRANSACTION**: Loads transactions from the given `FILE PATH` for the `ACCOUNT EMAIL`.
  If the account does not exist, it is created. After saving the transactions in the database, it generates the account summary and stores it.
- **VIEW TRANSACTIONS**: Displays all transactions associated with the given `ACCOUNT EMAIL` in **JSON format**.
- **VIEW SUMMARY**: Displays the summary of the account associated with the given `ACCOUNT EMAIL` in **formatted plain text**.
- **SEND EMAIL**: Sends the account summary of the account associated with the specified `ACCOUNT EMAIL` to that same email address.

## Prerequisites

- Go 1.25 or higher
- MongoDB 7.0 or higher
- Git

## Configuration

### Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# Application configuration
PORT=8080

# MongoDB configuration
MONGO_URI=mongodb://localhost:27017/transaction_summary

# Email configuration
EMAIL_FROM=your-email@gmail.com
EMAIL_PASSWORD=your-application-password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

### Configuring Gmail to Send Emails (EMAIL_PASSWORD)

If you are using Gmail, you need to:

1. **Enable Two-Factor Authentication (2FA)** in your Google account.
2. **Generate an "App Password"** from the Security settings.
3. **Use the generated App Password** instead of your regular Gmail password.

Note: If you just want to try it out, you can test this on the deployed site without setting up Gmail.

For more details, check the official documentation:
[Google Support – Sign in with App Passwords](https://support.google.com/accounts/answer/185833?hl=en)

## Installation and Execution

### Option 1: With Docker (Recommended)

#### Requirements

- Docker
- Docker Compose

#### Steps

1. **Clone the repository**

   ```bash
   git clone git@github.com:Galandrus/account-summary.git
   or
   git clone https://github.com/Galandrus/account-summary.git

   cd account-summary
   ```

2. **Configure environment variables**

   ```bash
   # Copy the example file
   cp .env.example .env

   # Edit .env with your values
   nano .env
   ```

3. **Start the services**

   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Application: http://localhost:8080
   - MongoDB: localhost:27017

#### Useful Docker commands

```bash
# Stop services
docker-compose down

# Rebuild and start
docker-compose up -d --build

# View logs for a specific service
docker-compose logs -f app
docker-compose logs -f mongoDatabase

# Access MongoDB container
docker-compose exec mongoDatabase mongosh
```

### Option 2: Manual Installation

#### Steps

1. **Clone the repository**

   ```bash
   git clone git@github.com:Galandrus/account-summary.git
   or
   git clone https://github.com/Galandrus/account-summary.git

   cd account-summary
   ```

2. **Install Go dependencies**

   ```bash
   go mod download
   ```

3. **Configure MongoDB**

   - Install MongoDB on your system
   - Create the `account_summary` database
   - Ensure MongoDB is running on port 27017

4. **Configure environment variables**

   ```bash
   # Copy the example file
   cp env.example .env

   # Edit .env with your values
   nano .env
   ```

5. **Run the application**

   ```bash
   go run main.go
   ```

6. **Access the application**
   - Application: http://localhost:8080

## Development

### File Structure

```
transactionSummary/
├── assets/                 # Static files (CSV, logos)
├── src/
│   ├── config/            # Application configuration
│   ├── connections/       # Database connections
│   ├── handlers/          # HTTP handlers
│   ├── interfaces/        # Interfaces
│   ├── models/            # Data models
│   ├── pkg/               # Utilities and libraries
│   ├── repository/        # Data access
│   ├── server/            # Server configuration
│   ├── services/          # Business logic
│   └── templates/         # HTML templates
├── docker-compose.yml     # Docker configuration
├── Dockerfile            # Docker image
├── go.mod               # Go dependencies
└── main.go              # Entry point
```
