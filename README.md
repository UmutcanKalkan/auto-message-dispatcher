# Auto Message Dispatcher

Otomatik mesaj gönderim sistemi.

## Gereksinimler

- Go 1.20+
- MongoDB
- Redis
- Docker & Docker Compose (opsiyonel)

## Kurulum

### Docker ile (Test Edilmedi)

```bash
# Servisleri başlat
docker-compose up -d

# API erişilebilir: http://localhost:8080
# Swagger UI: http://localhost:8080/swagger
```

### Manuel Kurulum

```bash
# Dependencies
go mod download

# .env dosyasını oluştur
cp .env.example .env
# .env dosyasını düzenle ve kendi değerlerini gir

# MongoDB başlat (local)
mongod

# Redis başlat (opsiyonel)
redis-server

# Uygulamayı çalıştır
go run cmd/server/main.go
```

## API Endpoints

### Scheduler Kontrolü
- `POST /scheduler/start` - Otomatik gönderimi başlat
- `POST /scheduler/stop` - Otomatik gönderimi durdur
- `GET /scheduler/status` - Durum kontrolü

### Mesaj İşlemleri
- `GET /messages/sent` - Gönderilen mesajları listele
- `POST /messages` - Yeni mesaj oluştur

## Konfigürasyon

`.env` dosyasını `.env.example` dosyasından kopyalayarak oluşturun ve kendi değerlerinizi girin.

Önemli değişkenler:
- `WEBHOOK_URL`: Mesaj gönderim endpoint'i
- `WEBHOOK_AUTH_KEY`: API authentication key
- `SCHEDULER_INTERVAL`: Mesaj gönderim aralığı (default: 2m)
- `SCHEDULER_BATCH_SIZE`: Her seferde kaç mesaj (default: 2)

## Swagger Dokümantasyonu

API dokümantasyonuna erişim: `http://localhost:8080/swagger`

