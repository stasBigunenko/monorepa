# monorepa

To run: 
make up

To use:
- Generate token:  http GET http://127.0.0.1:8080/login "user"="bob" "password"="123123"
- To get items: http GET http://127.0.0.1:8081/items 'Authorization: bearer <token>' "Name"="Peter"
