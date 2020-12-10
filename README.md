O learn about this sandbox and for instructions on how to run it please head over
to the [envoy docs](https://www.envoyproxy.io/docs/envoy/latest/start/sandboxes/front_proxy.html)

## Step 1: Start all of our containers

```bash
docker-compose build --pull
docker-compose up -d
```

## Step 2: Test Envoy’s routing capabilities

You can now send a request to both services via the `front-envoy`.

For `service1`:

```bash
$ curl -v localhost:8080/service/1
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /service/1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< date: Wed, 09 Dec 2020 17:15:19 GMT
< content-length: 92
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 4
< server: envoy
< x-envoy-decorator-operation: checkAvailabilityService1
< 
Hello from behind Envoy (service 1)! hostname: 135a315e57e5 resolved hostname: [172.18.0.4]
```

For `service2`:

```bash
$ curl -v localhost:8080/service/2
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /service/2 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< date: Wed, 09 Dec 2020 17:16:54 GMT
< content-length: 92
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 0
< server: envoy
< x-envoy-decorator-operation: checkAvailabilityService2
< 
Hello from behind Envoy (service 2)! hostname: f5fff62d6d2f resolved hostname: [172.18.0.5]
```

We can also use `HTTPS` to call services behind the front Envoy. For example, calling `service1`:

```bash
$ curl https://localhost:8443/service/1 -k -v
*   Trying 127.0.0.1:8443...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/certs/ca-certificates.crt
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_256_GCM_SHA384
* ALPN, server did not agree to a protocol
* Server certificate:
*  subject: CN=front-envoy
*  start date: Jul  8 01:31:46 2020 GMT
*  expire date: Jul  6 01:31:46 2030 GMT
*  issuer: CN=front-envoy
*  SSL certificate verify result: self signed certificate (18), continuing anyway.
> GET /service/1 HTTP/1.1
> Host: localhost:8443
> User-Agent: curl/7.68.0
> Accept: */*
> 
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* old SSL session ID is stale, removing
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< date: Wed, 09 Dec 2020 17:22:10 GMT
< content-length: 92
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 3
< server: envoy
< x-envoy-decorator-operation: checkAvailabilityService1
< 
Hello from behind Envoy (service 1)! hostname: 135a315e57e5 resolved hostname: [172.18.0.4]
```

## Step 3: Test Envoy’s load balancing capabilities

Now let’s scale up our `service1` nodes to demonstrate the load balancing abilities of Envoy:

```bash
$ docker-compose scale service1=3
Creating proxy_service1_2 ... done
Creating proxy_service1_3 ... done
```

Now if we send a request to `service1` multiple times, the front Envoy will load balance the requests by doing a round robin of the three  `service1` machines:

```bash
$ curl -v localhost:8080/service/1
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /service/1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< date: Wed, 09 Dec 2020 17:28:35 GMT
< content-length: 92
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 1
< server: envoy
< x-envoy-decorator-operation: checkAvailabilityService1
< 
Hello from behind Envoy (service 1)! hostname: 211f31c036dd resolved hostname: [172.18.0.7]

$ curl -v localhost:8080/service/1
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /service/1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< date: Wed, 09 Dec 2020 17:28:38 GMT
< content-length: 92
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 1
< server: envoy
< x-envoy-decorator-operation: checkAvailabilityService1
< 
Hello from behind Envoy (service 1)! hostname: 135a315e57e5 resolved hostname: [172.18.0.4]

$ curl -v localhost:8080/service/1
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /service/1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< date: Wed, 09 Dec 2020 17:28:39 GMT
< content-length: 92
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 0
< server: envoy
< x-envoy-decorator-operation: checkAvailabilityService1
< 
Hello from behind Envoy (service 1)! hostname: 71e8d5980400 resolved hostname: [172.18.0.6]
```