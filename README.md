# Loan Management API 🚀

This repository contains a **Loan Management API** built with **Go**, designed to handle key operations such as loan creation, approval, investment, disbursement, payment tracking, and delinquency status management. It follows a RESTful approach, ensuring scalability and simplicity.

---

## Features 🛠️

- **Loan Lifecycle Management**: Create, approve, reject, invest, and disburse loans.  
- **Payment Handling**: Track outstanding amounts and process weekly repayments.  
- **Delinquency Checks**: Identify delinquent borrowers based on missed payments.  
- **RESTful Endpoints**: Simple and consistent API design for ease of integration.  
- **Error Handling**: Meaningful HTTP status codes for easy troubleshooting.

---

## API Endpoints 📡

| Endpoint                               | Method | Description                                      |
|----------------------------------------|--------|--------------------------------------------------|
| `/loans`                               | GET    | List all loans                                   |
| `/loans`                               | POST   | Create a new loan                                |
| `/loans/{id}/approve`                  | PUT    | Approve a loan                                   |
| `/loans/{id}/reject`                   | PUT    | Reject a loan                                    |
| `/loans/{id}/invest`                   | PUT    | Invest in a loan                                 |
| `/loans/{id}/disburse`                 | PUT    | Disburse the loan                                |
| `/loans/{id}/outstanding?week={n}`     | GET    | Check outstanding amount for a specific week     |
| `/loans/{id}/delinquent?week={n}`      | GET    | Check if a borrower is delinquent                |
| `/loans/{id}/payment?week={n}`         | POST   | Make a payment for a specific week               |

---

## How to Run 🏃‍♂️

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/wsugiri/loansystem.git
   cd loan-management-api

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

4. **Test the endpoints using Postman or Curl:**
   Example:
   ```
   curl -X GET http://localhost:9002/api/loans
   ```

## Tech Stack 🧑‍💻
  - **Go**: Backend language
  - **RESTful API**: For smooth client-server communication
  - **Postman / Curl**: API testing tools   

## Conclusion
This API provides essential endpoints for managing loans, including state transitions, investments, payments, and delinquency checks. Follow the provided examples to interact with the API effectively.
