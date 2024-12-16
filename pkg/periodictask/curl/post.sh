curl -v -d '{"title":"title", "description":"description","schedule":"* * * * *"}' \
 -H "Content-Type: application/json" \
 --cookie "session_token=${SESSION_TOKEN}; Path=/; Max-Age=1733772315; Secure" \
 -X POST \
 http://localhost:8080/periodic-tasks/