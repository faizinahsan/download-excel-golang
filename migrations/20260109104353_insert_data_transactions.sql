-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.transactions (
	id bigserial NOT NULL,
	transaction_amount numeric(30, 5) NULL,
	reference_number varchar(12) NULL,
	transaction_date timestamp(6) NULL,
	feature_name varchar(255) NULL,
	CONSTRAINT t_transaction_pkey PRIMARY KEY (id)
);

-- Create index for better performance on common queries
CREATE INDEX IF NOT EXISTS idx_transactions_date ON public.transactions(transaction_date);
CREATE INDEX IF NOT EXISTS idx_transactions_ref ON public.transactions(reference_number);

-- Insert data in batches with optimized query
WITH data_batch AS (
    SELECT
        1000000.00000 AS transaction_amount,
        'TRX' || lpad(gs::text, 9, '0') AS reference_number,
        now() - (gs % 30) * interval '1 day' AS transaction_date,
        'Transaksi Dummy' AS feature_name
    FROM generate_series(1, 1000000) gs
)
INSERT INTO public.transactions (
    transaction_amount,
    reference_number,
    transaction_date,
    feature_name
)
SELECT * FROM data_batch;

-- Update table statistics for better query planning
ANALYZE public.transactions;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop indexes first
DROP INDEX IF EXISTS idx_transactions_date;
DROP INDEX IF EXISTS idx_transactions_ref;

-- Drop the table and all its data
DROP TABLE IF EXISTS public.transactions;
-- +goose StatementEnd
