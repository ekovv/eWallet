CREATE TABLE wallets (
     id SERIAL PRIMARY KEY,
     idOfWallet TEXT UNIQUE,
     balance NUMERIC(10, 1)
);

CREATE TABLE history (
      time TIMESTAMP WITH TIME ZONE NOT NULL,
      from_wallet VARCHAR(255) NOT NULL,
      to_wallet VARCHAR(255) NOT NULL,
      amount NUMERIC(10, 1) NOT NULL
);
