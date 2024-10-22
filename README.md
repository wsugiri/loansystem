# Loan Management API 🚀

This repository contains a **Loan Management API** built with **Go**, designed to handle key operations such as loan creation, approval, investment, disbursement, payment tracking, and delinquency status management. It follows a RESTful approach, ensuring scalability and simplicity.

---

## Key Features 🛠️

- **Create Loan**: Initiate new loan applications.  
- **Approve Loan**: Approve loans after field validation.  
- **Reject Loan**: Decline loans that don't meet requirements.  
- **Invest in Loan**: Enable investors to fund loans.  
- **Disburse Loan**: Transfer loan funds to borrowers.  
- **Outstanding Loan**: Track remaining loan amounts to be repaid.  
- **Check Delinquent Status**: Identify borrowers with missed payments.  
- **Make Payment**: Process weekly repayments from borrowers.  
- **Loan Lifecycle Management**: Oversee the complete lifecycle of loans.  
---

## Base URL
```
{base_url} = http://localhost:9001/api
```

## API Endpoints 📡
| No | Endpoint                         | Method | Description                                  |
|----|----------------------------------|--------|----------------------------------------------|
|  1 | `/loans`                         | POST   | Create a new loan                            |
|  2 | `/loans/{id}/approve`            | PUT    | Approve a loan                               |
|  3 | `/loans/{id}/reject`             | PUT    | Reject a loan                                |
|  4 | `/loans/{id}/invest`             | PUT    | Invest in a loan                             |
|  5 | `/loans/{id}/disburse`           | PUT    | Disburse the loan                            |
|  6 | `/loans`                         | GET    | List loans                                   |
|  7 | `/loans/{id}`                    | GET    | Detail transaction loans                     |
|  8 | `/loans/{id}/{week}/outstanding` | GET    | Check outstanding amount for a specific week |
|  9 | `/loans/{id}/{week}/delinquent`  | GET    | Check if a borrower is delinquent            |
| 10 | `/loans/{id}/{week}/payment`     | POST   | Make a payment for a specific week           |
---

### 1. Create a new loan
Create a new loan.
#### Request
```
POST {base_url}/loans
Content-Type: application/json
```

#### Request Body
```json
{
  "borrower_id": 1,
  "principal_amount": 5000000,
  "interest_rate": 10,
  "loan_duration_weeks": 50,
  "agreement_url": "https://image-upload.io/loans/borrower_1.jpg"
}
```

#### Response
```json
{
   "id": 9,
   "status": "proposed",
   "data": {
      "duration_weeks": 50,
      "instalment": 110000,
      "total_loan": 5500000
   }
}
```

#### Sample Response Error
```json
{
   "error": "unregistered_borrower"
}
```


### 2. Approve a loan
#### Request
```
PUT {base_url}/loans/9/approve
Content-Type: application/json
```

#### Request Body
```json
{
  "employee_id": 4,
  "approval_date": "2024-12-01 14:20",
  "validator_photo": "string"  // URL or base64 image
}
```

#### Response
```json
{
   "id": 1,
   "status": "approved",
   "data": {
      "duration_weeks": 50,
      "total_loan": 5500000,
      "instalments": [
         {
            "Week": 1,
            "Amount": 110000,
            "DueDate": "2024-12-08"
         },
         {
            "Week": 2,
            "Amount": 110000,
            "DueDate": "2024-12-15"
         },
         {
            "Week": 3,
            "Amount": 110000,
            "DueDate": "2024-12-22"
         },
         ...
      ]
   }
}
```

#### Sample Response Error
```json
{
   "error": "unregistered_approver",
}
```


### 3. Reject a loan
#### Request
```
PUT {base_url}/loans/7/reject
Content-Type: application/json
```

#### Request Body
```json
{
  "employee_id": 5,
  "rejection_date": "2024-02-01",
  "rejection_message": "string"  // reason message 
}
```

#### Response
```json
{
   "id": 7,
   "status": "rejected",
   "rejection_date": "2024-02-01"
}
```

#### Sample Response Error
```json
{
   "error": "unregistered_rejector",
}
```


## How to Run 🏃‍♂️

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/wsugiri/loansystem.git
   cd loansystem

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set Up the Database Schema:**
   - Run the SQL script to create the necessary tables:
   ```bash
   mysql -u your_username -p your_database_name < 001_create_tables.sql
   ```
   - Run the SQL script to insert sample initial data:
   ```bash
   mysql -u your_username -p your_database_name < 002_insert_sample_users_data.sql
   ```

4. **Rename the Environment File:**
   - Rename .env.example to .env

5. **Update Environment Variables:**
   - Open the .env file in your favorite text editor and change the values as needed for your configuration.

6. **Run the server:**
   ```bash
   go run main.go
   ```

7. **Test the endpoints using Postman or Curl:**
   Example:
   ```
   curl -X GET http://localhost:9001/api/loans
   ```

## Tech Stack 🧑‍💻
  - **Go**: Backend language
  - **github.com/gofiber/fiber/v2**: Framework for creating RESTful APIs
  - **RESTful API**: For smooth client-server communication
  - **MySQL**: Database for data storage
  - **Postman / Curl**: API testing tools   

## Conclusion
This API provides essential endpoints for managing loans, including state transitions, investments, payments, and delinquency checks. Follow the provided examples to interact with the API effectively.
