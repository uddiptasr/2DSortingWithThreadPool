
code is live on Render
curls :

process-single---
"curl --location --request GET 'https://go-sort-api.onrender.com/process-single' \
--header 'Content-Type: application/json' \
--data '{
  "to_sort": [[21, 12, 3], [14, 15, 5], [8, 9, 1]]
}'
"

process-concurrent --
"curl --location --request GET 'https://go-sort-api.onrender.com/process-concurrent' \
--header 'Content-Type: application/json' \
--data '{
  "to_sort": [[21, 12, 3], [14, 15, 5], [8, 9, 1]]
}' "
