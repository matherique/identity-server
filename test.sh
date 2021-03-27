echo "GET"

curl \
  --header "Content-Type: application/json" \
  -X GET http://localhost:8000

echo "POST"

curl \
  --header "Content-Type: application/json" \
  -X POST http://localhost:8000/auth \
  --data '{"id": "service1"}' \


