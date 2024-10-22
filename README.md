# Loan Management API 🚀

This repository contains a **Loan Management API** built with **Go**, designed to handle key operations such as loan creation, approval, investment, disbursement, payment tracking, and delinquency status management. It follows a RESTful approach, ensuring scalability and simplicity.

---

## Key Features 🛠️

- **Create Loan**: Initiate new loan applications.  
- **Approve Loan**: Approve loans after field validation.  
- **Reject Loan**: Decline loans that don't meet requirements.  
- **Invest in Loan**: Enable investors to fund loans.  
- **Disburse Loan**: Transfer loan funds to borrowers.  
- **List Loan**: List available loans with status.  
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
---

### 1. Create a new loan
The create new loan feature allows authorized staff to initiate and set up a loan application within the system. This process captures essential details about the loan and the borrower, facilitating the lending workflow.

   When creating a new loan, the following information must be provided:
   - **Borower Id**: A unique identifier for the borrower in the system.
   - **Principal Amount**: The total amount of money being borrowed.
   - **Interest Rate**: The percentage of interest that will be charged on the principal amount.
   - **Loan Duration**: The number of weeks over which the loan will be repaid.
   - **Agreement URL**: A link to the loan agreement that outlines the terms and conditions.

#### Request
```
POST {base_url}/loans
Content-Type: application/json
```

#### Request Body
```json5
{
  "borrower_id": 1,
  "principal_amount": 5000000,
  "interest_rate": 10,
  "loan_duration_weeks": 50,
  "agreement_url": "https://image-upload.io/loans/borrower_1.jpg"
}
```

#### Response
```json5
{
  "status": "success",
  "message": "Loan successfully created",
  "data": {
    "loan_id": 1001,
    "borrower_id": 1,
    "principal_amount": 5000000,
    "interest_rate": 10,
    "loan_duration_weeks": 50,
    "total_loan": 5500000,  // Principal + interest (principal * interest_rate * 0.01)
    "agreement_url": "https://image-upload.io/loans/borrower_1.jpg"
  }
}
```

#### Sample Response Error
```json5
{
   "status": "error",
   "message": "invalid borrower_id"
}
```


### 2. Approve a loan
The loan approval process is a crucial step in the lending workflow, ensuring that loans are thoroughly vetted before being offered to investors or lenders.
1. **Approval Requirements**: 
   Each loan approval must include the following information:
   - **Proof of Visit**: A picture confirming that a field validator has visited the borrower.
   - **Employee ID**: The identification number of the field validator who conducted the visit.
   - **Date of Approval**: The date when the loan was approved.

2. **Irreversible Approval**: 
   Once a loan is approved, it cannot revert back to the "proposed" state. This ensures the integrity of the approval process and maintains a clear workflow.

3. **Investor Readiness**: 
   After approval, the loan is ready to be offered to investors or lenders, enabling the funding process to begin.


#### Request
```
PUT {base_url}/loans/9/approve
Content-Type: application/json
```

#### Request Body
```json5
{
  "employee_id": 4,
  "approval_date": "2024-12-01 14:20",
  "validator_photo": "string"  // URL or base64 image
}
```

#### Response
```json5
{
   "status": "success",
   "message": "Loan successfully approved",
   "data": {
      "loan_id": 1001,
      "borrower_id": 1,
      "approval_date": "2024-12-01",
      "employee_id": 5,
      "validator_photo": "https://example.com/validator_photo.jpg",
      "loan_status": "approved"
   }
}
```

#### Sample Response Error
```json5
{
   "status": "error",
   "message": "invalid employee_id"
}
```


### 3. Reject a loan
The loan rejection process allows staff to deny loan applications that do not meet the necessary criteria or pose potential risks. This step is crucial for maintaining the integrity of the lending system.

For instance, if a loan application is rejected due to insufficient income verification, the staff member will document the reason and notify the borrower. The rejected loan will remain in the system for record-keeping but will not proceed to the approval phase.

#### Request
```
PUT {base_url}/loans/7/reject
Content-Type: application/json
```

#### Request Body
```json5
{
  "employee_id": 5,
  "rejection_date": "2024-02-01",
  "rejection_message": "string"  // reason message 
}
```

#### Response
```json5
{
   "status": "success",
   "message": "Loan successfully rejected",
   "data": {
      "loan_id": 9,
      "borrower_id": 3,
      "employee_id": 4,
      "loan_status": "rejected",
      "rejection_date": "2024-02-01",
      "rejection_message": "tidak memenuhi syarat"
   }
}
```

#### Sample Response Error
```json5
{
   "status": "error",
   "message": "invalid employee_id"
}
```

### 4. Invest in a loan
This feature allows multiple investors to contribute to a single loan. Each investor can invest a unique amount towards the loan, enabling flexible funding options.

**Key Points**
- **Multiple Investors**: A loan can have multiple investors, with each investor contributing their own amount. This allows for diverse participation in the funding process.
  
- **Investment Limit**: The total amount invested by all investors combined cannot exceed the loan's principal amount. This ensures that the funding is balanced and aligns with the loan's initial value.

#### Request
```
PUT {base_url}/loans/2/invest
Content-Type: application/json
```

#### Request Body
```json5
{
  "investor_id": 6,
  "amount": 500000
}
```

#### Sample Response
```json5
{
   "status": "success",
   "message": "Investment successfully made",
   "data": {
      "loan_id": 4,
      "borrower_id": 3,
      "investor_id": 7,
      "investment_amount": 150000,
      "loan_principal_amount": 5000000,
      "remaining_amount": 4600000,
      "total_invested": 400000
   }
}
```

#### Sample Response Error
```json5
{
   "status": "error",
   "message": "cannot invest more than 450000"
}
```


### 5. Disburse the loan
The disbursement of a loan marks the final step in the loan process, where the approved loan amount is transferred to the borrower. This step ensures that the loan is officially activated and funds are made available to the borrower.

#### Request
```
PUT {base_url}/loans/2/disburse
Content-Type: application/json
```

#### Request Body
```json5
{
  "employee_id": "string",
  "disbursement_date": "YYYY-MM-DD",
  "agreement_letter": "string"               // URL or base64 image
}
```

#### Sample Response
```json5
{
  "status": "success",
  "message": "Loan successfully disbursed",
  "data": {
    "loan_id": 1001,
    "borrower_id": 12345,
    "disbursement_date": "2024-11-01",
    "disbursed_amount": 500000,
    "employee_id": 128,
    "agreement_letter": "https://example.com/agreement/1001"
  }
}
```

#### Sample Response Error
```json5
{
  "status": "error",
  "message": "already disburse",
}
```
```json5
{
  "status": "error",
  "message": "Loan is not in an approved state",
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
