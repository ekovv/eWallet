# eWallet

![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

# 🎲 Service on Go(Gin) for implementing a payment system transaction processing system.

# 📞 Endpoints
```http
POST /api/v1/wallet
- сreate wallet
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
d - db connection string -d=connection string
s - salt -s=your_salt
```

# 🧩 Config

```json
{
  "host": "localhost:8080",
  "dsn": "postgres://bestuser:bestuser@localhost:5432/your_db_name?sslmode=disable",
  "salt": "your_salt"
}
```

# 💎 For working with Docker
-Build postgres
```
docker run --name ewallet-pg -p 4999:5432 -e POSTGRES_USER=username -e POSTGRES_PASSWORD=password -e POSTGRES_DB=db_name -d postgres:13.3
```
-Run
```
docker run --name your-app-name -p host-port:container-port -d your-app-image
```
