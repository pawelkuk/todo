curl -v localhost:8080/tasks/ \
 --cookie "session_token=${SESSION_TOKEN}; Path=/; Max-Age=1733772315; Secure" \
 | jq

curl -v "localhost:8080/tasks/?before=2024-12-31&after=2024-12-01&completed=true&dueDate=2024-12-15&title=abc&description=abc" \
 --cookie "session_token=${SESSION_TOKEN}; Path=/; Max-Age=1733772315; Secure" \
 | jq