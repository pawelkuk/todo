curl -v \
 -d '{"title":"updated title","description":"updated description","dueDate":"2024-12-25"}' \
 -H "Content-Type: application/json" \
 -X PATCH \
 --cookie "session_token=${SESSION_TOKEN}; Path=/; Max-Age=1733772315; Secure" \
 http://localhost:8080/tasks/1 | jq