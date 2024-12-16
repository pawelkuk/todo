curl -v \
 -H "Content-Type: application/json" \
 -X DELETE \
 --cookie "session_token=${SESSION_TOKEN}; Path=/; Max-Age=1733772315; Secure" \
 http://localhost:8080/periodic-tasks/1