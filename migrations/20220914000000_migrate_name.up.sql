CREATE TABLE IF NOT EXISTS player(
    player_name VARCHAR(255) NOT NULL PRIMARY KEY,
    free_rounds integer CHECK (free_rounds >= 0)
);

CREATE TABLE IF NOT EXISTS balance(
    player_name VARCHAR(255) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    balance integer CHECK (balance.balance >= 0),
    CONSTRAINT player_currency PRIMARY KEY(player_name, currency)
);

CREATE TABLE IF NOT EXISTS transaction(
    transaction_ref VARCHAR(255) NOT NULL PRIMARY KEY,
    player_name VARCHAR(255),
    withdraw integer CHECK (withdraw >= 0),
    deposit integer CHECK (deposit >= 0),
    currency VARCHAR(3),
    charge_free_rounds integer CHECK (charge_free_rounds >= 0),
    status integer CHECK (status IN(0, 1))
);
