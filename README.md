## airvet jwt

### Run the server:
```
make run

# also we can build and exec:
# make build
# ./server
```


### Usage:
Get a JWT token with any of the users, using the `POST /auth` endpoint:
if you have jq installed (easy way):

```
export JWT=$(curl -s -d '{ "email": "admin@airvet.com", "password": "Admin-pass" }' localhost:8080/auth | jq -r '.jwt')
```

or just:
```
curl -i -d '{ "email": "admin@airvet.com", "password": "Admin-pass" }' localhost:8080/auth

curl -i -d '{ "email": "coolvet@airvet.com", "password": "Cool_pass123" }' localhost:8080/auth
```

Copy the JWT and paste it in the follow command and you should get the user profile, using the `GET /user` endpoint:
```
curl -i -H "Authorization: Bearer $JWT" localhost:8080/user
```

### JWKs:

The JWK (private and public keys) are generated when we run the server 


#### Get all JWKS:
Hit the `GET /jwks` endpoint
```
curl -i localhost:8080/.well-known/jwks.json
```

#### Rotate keys:
```
make rotate
```

then call `curl -i localhost:8080/.well-known/jwks.json` to get the new JWK

