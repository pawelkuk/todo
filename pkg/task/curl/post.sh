curl -v -d '{"title":"title", "description":"description","dueDate":"2024-12-25"}' \
 -H "Content-Type: application/json" \
 -X POST \
 --cookie "session_token=${SESSION_TOKEN}; Path=/; Max-Age=1733772315; Secure" \
 http://localhost:8080/tasks/ | jq