-- Create users table to represent borrowers and staff
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role ENUM('borrower', 'staff', 'investor') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create loans table to store loan information
CREATE TABLE loans (
    id INT AUTO_INCREMENT PRIMARY KEY,
    borrower_id INT NOT NULL,
    principal_amount DECIMAL(15, 2) NOT NULL,
    rate DECIMAL(5, 2) NOT NULL, -- Interest rate
    total_loan DECIMAL(15, 2) NOT NULL,
    instalment DECIMAL(15, 2) NOT NULL,
    status ENUM('proposed', 'approved', 'rejected', 'invested', 'disbursed') DEFAULT 'proposed',
    agreement_url VARCHAR(500), -- Link to agreement letter (PDF)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (borrower_id) REFERENCES users(id)
);

-- Create approvals table to store approval details
CREATE TABLE approvals (
    id INT AUTO_INCREMENT PRIMARY KEY,
    loan_id INT NOT NULL,
    picture_proof_url VARCHAR(500),
    approval_date DATE NOT NULL,
    approval_by INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (approval_by) REFERENCES users(id)
);

-- Create rejections table to store rejection details
CREATE TABLE rejections (
    id INT AUTO_INCREMENT PRIMARY KEY,
    loan_id INT NOT NULL,
    rejection_reason VARCHAR(255) NOT NULL,
    rejection_date DATE NOT NULL,
    rejected_by INT NOT NULL,
    comments TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (rejected_by) REFERENCES users(id)
);

-- Create investments table to track investor contributions
CREATE TABLE investments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    loan_id INT NOT NULL,
    investor_id INT NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (investor_id) REFERENCES users(id)
);

-- Create disbursements table to store loan disbursement details
CREATE TABLE disbursements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    loan_id INT NOT NULL,
    field_officer_id INT NOT NULL,
    disbursement_date DATE NOT NULL,
    agreement_signed_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (field_officer_id) REFERENCES users(id)
);

-- Create payments table to track repayments
CREATE TABLE payments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    loan_id INT NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    payment_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id)
);

-- Create a view to calculate outstanding loan balance
CREATE OR REPLACE VIEW loan_outstanding AS
SELECT 
    l.id AS loan_id,
    l.principal_amount - IFNULL(SUM(p.amount), 0) AS outstanding_amount
FROM loans l
LEFT JOIN payments p ON l.id = p.loan_id
GROUP BY l.id;
