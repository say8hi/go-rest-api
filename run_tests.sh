
#!/bin/bash

docker-compose --profile test up -d

sleep 5

go test -v ./tests/

docker-compose --profile test down -v

if [ $? -eq 0 ]; then
  echo "Тесты успешно завершены!"
else
  echo "Тесты завершились с ошибками!"
  exit 1
fi
