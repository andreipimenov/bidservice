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