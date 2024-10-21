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
| Endpoint                               | Method | Description                                      |
|----------------------------------------|--------|--------------------------------------------------|
| `/loans`                               | POST   | Create a new loan                                |
| `/loans/{id}/approve`                  | PUT    | Approve a loan                                   |
| `/loans/{id}/reject`                   | PUT    | Reject a loan                                    |
| `/loans/{id}/invest`                   | PUT    | Invest in a loan                                 |
| `/loans/{id}/disburse`                 | PUT    | Disburse the loan                                |
| `/loans`                               | GET    | List loans                                       |
| `/loans/{id}`                          | GET    | Detail transaction loans                         |
| `/loans/{id}/outstanding?week={n}`     | GET    | Check outstanding amount for a specific week     |
| `/loans/{id}/delinquent?week={n}`      | GET    | Check if a borrower is delinquent                |
| `/loans/{id}/payment?week={n}`         | POST   | Make a payment for a specific week               |
---

## How to Run 🏃‍♂️

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/wsugiri/loansystem.git
   cd loansystem

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
   curl -X GET http://localhost:9001/api/loans
   ```

## Tech Stack 🧑‍💻
  - **Go**: Backend language
  - **RESTful API**: For smooth client-server communication
  - **Postman / Curl**: API testing tools   

## Conclusion
This API provides essential endpoints for managing loans, including state transitions, investments, payments, and delinquency checks. Follow the provided examples to interact with the API effectively.
