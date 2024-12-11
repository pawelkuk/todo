curl -v localhost:8080/tasks/ \
 --cookie 'session_token=7d08852c-87a4-40e6-9d70-d3ba969283d8; Path=/; Max-Age=1733772315; Secure' \
 | jq