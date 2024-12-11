curl -v \
 -d '{"email":"paffcio315@gmail.com","password":"password2"}' \
 -H "Content-Type: application/json" \
 -X PATCH \
 http://localhost:8080/users/1 | jq