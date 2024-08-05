# Examples

Below are some example configurations demonstrating different features of `server-faker`.
To use an example run command:
```
server-faker run --file={file_name}
```
Where `{file_name}` is the example you want to run.

## custom-logger.json

This example shows how to use a custom logger plugin. Is uses the wildcard symbol to catch any POST request starting with `/logger/`. 
This loggs all the possible information about the request like payload, headers, path, method, ip and responds with 200 status code. 

**Method:** `POST`

**URL:** `/logger/*`

```sh
curl --request POST \
  --url http://127.0.0.1:8080/logger/test

curl --request POST \
  --url http://127.0.0.1:8080/logger?query=value
```

> **Note:** The custom plugin must be built with name `endpoint_logger.so` and placed in the `../plugins` directory, which is relative to the provided file.


## custom-proto.json

This example demonstrates protobuf usage with a custom plugin generating protobuf response.
This generates a protobuf response for the URL `/proto` and method `GET`. 

**Method:** `GET`

**URL:** `/proto`

```sh
curl --request GET \
  --url http://127.0.0.1:8080/proto
```

> **Note:** The custom plugin must be built with name `protobuf_response.so` and placed in the `../plugins` directory, which is relative to the provided file.

## dynamic-array-with-objects.json

This example returns a dynamic array with JSON objects at the /dynamic URL using the GET method.

**Method:** `POST`

**URL:** `/dynamic-proxy/post/:test_id`

```sh
curl --request GET \
  --url http://127.0.0.1:8080/dynamic
```

## dynamic-proxy.json

This example presents a dynamic proxy setup. 
It proxies a GET request to `https://jsonplaceholder.typicode.com/posts/:proxy_id` from a POST request to `/dynamic-proxy/post/:test_id`.
It showcases how to use the mapping parameter from URL to proxyed URL.

**Method:** `POST`

**URL:** `/dynamic-proxy/post/:test_id`

```sh
curl --request POST \
  --url http://127.0.0.1:8080/dynamic-proxy/post/1 
```

## static-proxy.json

This example forwards the response from GET `/proxy/posts` to `https://jsonplaceholder.typicode.com/posts`.

**Method:** `GET`

**URL:** `/proxy/posts`

```sh
curl --request GET \
  --url http://127.0.0.1:8080/proxy/posts 
```

## static.json

This example responds with static JSON from the file "static-proxy.json" on the GET `"/static"` endpoint.

**Method:** `GET`

**URL:** `/static`

```sh
curl --request GET \
  --url http://127.0.0.1:8080/static
``` 