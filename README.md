# Identity Server

A mini identity server

## How it works

```
Service A                                             IS
               -- (1)Authorization Request ->

               <- (2)Authorization Grant --
                         Token IS

               -- (3)Authorization Request ->
                         to Service B

               <- (4)Authorization Grant -- 
                         Token Service B

```

(1)  The service A request authorization token to use IS.

     body: { id: "serviceAId" }

(2)  If service A is configurated in IS, its responds with 
     token and refreshToken

     body { accessToken: "token", refreshToken: "refreshToken" }

(3)  The service A request authorization to use Serice B 

     body { id: "serviceBid" }

(4)  The IS validate if service A depends on service B, if true, send
     the token to access service B, other wise return error 

     body { token: "serviceBtoken" }



The only token that refresh is the IS token, the others expires in X time and
you need to request the new token to IS

