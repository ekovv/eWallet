CREATE TABLE wallets (
     id SERIAL PRIMARY KEY,
     idOfWallet TEXT UNIQUE,
     balance NUMERIC(10, 1)
);