### Как тестировать локально

```sh
cd bot
cp .env.example .env       # заполнить TELEGRAM_TOKEN, GITHUB_TOKEN, ALLOWED_USER_ID
make local                 # go run ./cmd/local — слушает :8080
```

```sh
# В другом терминале:
ngrok http 8080
# Установить вебхук:
curl -F "url=https://your-ngrok.ngrok.io/api/bot" \
  https://api.telegram.org/bot<TELEGRAM_TOKEN>/setWebhook
```