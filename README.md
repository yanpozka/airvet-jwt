## airvet jwt

### Run the server:
```
make run
```


### Testing:
Get a JWT token with any of the users:
```
curl -i -d '{ "email": "admin@airvet.com", "password": "Admin-pass" }' localhost:8080/auth

curl -i -d '{ "email": "coolvet@airvet.com", "password": "Cool_pass123" }' localhost:8080/auth

# if you have jq installed:
export JWT=$(curl -s -d '{ "email": "admin@airvet.com", "password": "Admin-pass" }' localhost:8080/auth | jq -r '.jwt')
```

Copy the JWT on the follow command and you should get the user profile:
```
curl -i -H "Authorization: Bearer $JWT" localhost:8080/user
```


#### Get all JWKS:
```
curl -i localhost:8080/jwks
```
