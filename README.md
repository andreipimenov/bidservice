# Bid Service
Implementation of http bid service.

### Description

Application accepts requests on **/winner** with method **GET**
Query parameters like **s=example.com/fibo** represents external sources
External sorces must return response as JSON array with objects like 
```
[
    {"price": 1},
    {"price": 5}
]
```

Service must response as JSON like 
```
{
    "sorce": "example.com/fibo",
    "price": 1
}
```
Responsed price is the second highest price from all.

#### Details

1. Response time is 100ms or less
2. Using only build-in go libraries
3. Ability to configure TCP port
4. Tests with test sources from <https://github.com/trafficstars/test-job/blob/master/testsources.go>

### Testing

#### Unit testing of main http handler
```
cd cmd/bid
go test -v .
```
```
=== RUN   TestWinnerHandler
--- PASS: TestWinnerHandler (0.08s)
PASS
ok      github.com/andreipimenov/bidservice/cmd/bid     0.465s
```
#### Manual testing
1. Get the app with test sources and run it
```
go get github.com/trafficstars/test-job
cd $GOPATH/src/github.com/trafficstars/test-job
go build -o testsource.exe && testsource.exe
```
2. Run bid service on another port (default port is the same 8080)
```
go build -o bid.exe && bid.exe -p 8000
```
3. Send helth check request
```
curl -X GET 127.0.0.1:8000/ping
```
```
{"message": "pong"}
```
4. Send request to forbidden uri
```
curl -X GET 127.0.0.1:8000/some-other-uri 
```
```
{"code": "Forbidden","message": "This uri is forbidden"}
```
5. Send request to /winner to get different results
```
curl -X GET 127.0.0.1:8000/winner?s=http://127.0.0.1:8080/primes&s=http://127.0.0.1:8080/fibo&s=http://127.0.0.1:8080/fact
```
```
{"uri": "http://127.0.0.1:8080/fact", "price": 6}
```
If no one source responsed correctly (valid json and with less than 100ms), the response is
```
```
{"code": "InternalServerError", "message": "prices not found, there is no winner"}
```
6. Response timeout can being configured by changing const SourceTimeout
