curl -v \
 -d '{"title":"updated title","description":"updated description","schedule":"0 1 * * *"}' \
 -H "Content-Type: application/json" \
 --cookie "session_token=${SESSION_TOKEN}; Path=/; Max-Age=1733772315; Secure" \
 -X PATCH \
 http://localhost:8080/periodic-tasks/1