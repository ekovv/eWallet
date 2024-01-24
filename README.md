# eWallet

![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

# 🎲 Service on Go(Gin) for implementing a payment system transaction processing system.

# 📞 Endpoints
```http
POST /api/v1/wallet
- Create wallet
POST /api/v1/wallet/:walletId/send
- transcations
GET /api/v1/wallet/:walletId/history
- history
GET /api/v1/wallet/:walletId
- status
```

# 🏴‍☠️ Flags
```
a - ip for REST -a=host
d - connection string -d=connection string
```

# 🧩 Config

```json
{
  "host": "localhost:8080",
  "dsn": "postgres://bestuser:bestuser@localhost:5432/your_db_name?sslmode=disable"
}
```

# 💎 For working with Docker


