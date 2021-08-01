## airvet jwt

### Testing:
Get a JWT token with any of the users:
```
curl -i -d '{ "email": "admin@airvet.com", "password": "Admin-pass" }' localhost:8080/auth

curl -i -d '{ "email": "coolvet@airvet.com", "password": "Cool_pass123" }' localhost:8080/auth
```

Copy the JWT on the follow command and you should get the user profile:
```
curl -i -H 'Authorization: Bearer <COPY-JWT-HERE>' localhost:8080/user
```
