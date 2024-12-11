curl -v -d '{"title":"title", "description":"description","dueDate":"2024-12-25"}' \
 -H "Content-Type: application/json" \
 -X POST \
 --cookie 'session_token=7d08852c-87a4-40e6-9d70-d3ba969283d8; Path=/; Max-Age=1733772315; Secure' \
 http://localhost:8080/tasks/ | jq