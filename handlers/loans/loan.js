exports.createLoan = async (req) => {
    try {
        const payload = req.body || {};
        const total_loan = + payload.principal_amount * +payload.interest_rate * 0.01;
        const instalment = total_amount / loan_duration_weeks;

        return {
            data: {
                total_loan,
                loan_duration_weeks,
                instalment,
            }
        }
    } catch (err) {
        throw err;
    }
}