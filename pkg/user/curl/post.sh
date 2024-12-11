curl -v -d '{"email":"paffcio314@gmail.com", "password":"password"}' \
 -H "Content-Type: application/json" \
 -X POST \
 http://localhost:8080/users/ | jq