echo "GET / - check if API is up"
curl --silent \
  --header "Content-Type: application/json" \
  -X GET http://localhost:8000 | jq


echo "--------------------------------------"
echo "POST /auth - no token"
curl --silent \
  --header "Content-Type: application/json" \
  -X POST http://localhost:8000/auth \
  --data '{"id": "service1"}' | jq


echo "--------------------------------------"
echo "POST /auth - no token and invalid service"
curl --silent \
  --header "Content-Type: application/json" \
  -X POST http://localhost:8000/auth \
  --data '{"id": "service5"}' | jq


echo "--------------------------------------"
echo "POST /auth - with token"
curl --silent \
  --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTY5NDE3MjQsImlkIjoic2VydmljZTEifQ.Va89vSc1dLcEo4b9T0uOlKuTjo9MMnQDXj00qGi-erg" \
  --header "Content-Type: application/json" \
  --data '{"id": "service2"}' \
  -X POST http://localhost:8000/auth | jq


echo "--------------------------------------"
echo "POST /auth - with token and invalid service"
curl --silent \
  --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTY5NDE3MjQsImlkIjoic2VydmljZTEifQ.Va89vSc1dLcEo4b9T0uOlKuTjo9MMnQDXj00qGi-erg" \
  --header "Content-Type: application/json" \
  --data '{"id": "service3"}' \
  -X POST http://localhost:8000/auth | jq
