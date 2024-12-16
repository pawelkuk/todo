curl -v localhost:8080/tasks/1 \
 --cookie "session_token=${SESSION_TOKEN}; Path=/; Max-Age=1733772315; Secure" \
 | jq