# `server-faker` Documentation

`server-faker` is a program designed to mock API responses for prepared endpoints specified in a JSON file. It operates in two modes: parser and server.

## Table of Contents
- [Quickstart](#Quickstart)
  - [Parser mode](#parser-mode)
  - [Server mode](#server-mode)
  - [Creating the first endpoint](#creating-first-endpoint)
- [Serve Content and Examples](#serve-content-and-examples)
  - [Serve Static Content](#serve-static-content)
  - [Serve Dynamic Content](#serve-dynamic-content)
    - [Dynamic Configuration Options](dynamic_configuration.md)
  - [Serve Custom Content](#serve-custom-content)
- [Proxy the request](#creating-proxy-endpoint)

---
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/vimekgo)

## Quickstart

### Docker image

It's possible to download a Docker image from [docker.io](https://hub.docker.com/r/vimekgo/server-faker) or build it and run the following command:

```sh
# build and run: 
docker build  . -t server-faker
docker run -d -p 8080:8080 server-faker

# pull and run
docker pull vimekgo/server-faker
docker run -d -p 8080:8080 server-faker
```

The default configuration runs a static JSON example.

### Parse a server file from JSON

```sh
docker run -v ./api-dir:/fake_api server-faker ./server-faker parse --file=fake_api/{file-name}.json server-faker
```

This prints to console the content of the file for further modifications. 
You can save the output in JSON and run the file created. 


#### Run your server-file

To run your server-file, use the following command:

```sh
docker run -d -p 8080:8080 -v {server-file-dir}:/fake_api server-faker ./server-faker run --file=fake_api/{server-file}
```

Replace `{server-file-dir}` with the directory containing your server-file and `{server-file}` with the name of your server-file.

Now you can access your mocked API at http://localhost:8080.

## Docker Compose Setup

If you prefer using Docker Compose for managing your containers, you can set up `server-faker` easily. Below is an example `docker-compose.yml` file:

```yaml
services:
  server-faker:
    image: vimekgo/server-faker
    ports:
      - "8080:8080"
    volumes:
      - ./fake/api:/fake/api
    command: server-faker run /fake/api/{file-name}
```

```sh
docker-compose up -d
```
This will start server-faker using Docker Compose, making your mocked API available at http://localhost:8080.

Make sure to adapt the paths and configurations in docker-compose.yml to match your project setup and server-file location.


## Definitions 
There are 2 different files that this documentation is referring to.

`server-file` is a JSON file that is used to run the server. It defines all the endpoints.

`input-file` is a file with an example JSON object that could be transformed with `parser` to a `server-file`. 

## Parser mode

In this mode, `server-faker` parses the given JSON file(`input-file`) and outputs a new JSON file(`server-file`) suitable for running the server. 
The output is displayed in the console by default. It can be redirected to a file.


```sh
server-faker parse --file=./{input-file}.json --type=dynamic
```

Arguments

    -f, --file: Specifies the path to the JSON file to be parsed.
    -t, --type: Specifies the parsing type. Options are:
        dynamic: Generates random values for given fields based on their original type.
        static: Outputs the same JSON as provided.
    -u, --url: Specyfies the endpoint url, where the content is served. 

Example

To parse `test-ev.json` with dynamic values and save the output to `output.json`:

```sh
server-faker parse --file=./test-ev.json --type=dynamic > output.json
```

## Server mode

In this mode, the `server-faker` runs a mocked API server on a specified port using the provided JSON file(`server-file`) to define endpoints and responses.


```sh
server-faker run --file=./test-api.json --port=8080
```

Arguments

    -f, --file: Specifies the path to the JSON file that defines the API endpoints and responses.
    -p, --port: Specifies the port on which the server will run.

Example

To run the server on port `8080` using the `test-api.json` configuration:

```sh
server-faker run --file=./test-api.json --port=8080
```

## Creating first endpoint

### URL Structure

- The URL given for an endpoint must start with `/`.
- For matching paths with dynamic sections like `product/{id}/category`, you can set the endpoint with a parameter prefixed with `:`: `/product/:id/category`.
- It is also possible to use a wildcard argument `*` to match any URL. The `*` must be at the end of URL. 
For example, the endpoint `/test/*` matches all paths that start with `test` regardless of their depth, such as `test/1`, `test/category/1`, and so on.

#### Note on wildcard endpoints

> Wildcard endpoints (e.g., `/*` and `/test/*`) cannot overlap because there is no way to determine which endpoint should handle requests like `/test/1/category`. 
Ensure that your wildcard endpoint patterns are distinct to avoid conflicts. 
Overlapping wildcard endpoints throws an error at startup.

> URLs with `*` in the middle are invalid. To match this kind of URL use the params with `:`. 

### Required Fields

- `url`: The URL path for the endpoint, starting with `/`.
- `method`: The HTTP method for the endpoint (e.g., `GET`, `POST`).
- `response`: The response configuration.
  - `status`: The HTTP status code to be returned by the endpoint.
  - `type`: The type of endpoint, which can be `static`, `dynamic`, or `custom`.

### Example: Static Endpoint

This is the smallest endpoint example serving a static file to the user for the `GET` method on the endpoint `test/test`.

```json
{
  "endpoints": [
    {
      "method": "GET",
      "url": "/test/test",
      "response": {
        "status": 200,
        "type": "static",
        "file": "test-test.json",
        "format": "json"
      }
    }
  ]
}
```

# Serve content and examples

## Serve static content 

For static content, users can provide a file that will be served. 
The path to the file is relative to the JSON file that the server consumes. 
Below are examples of how to serve a JSON file and a PNG image.

### Serve a JSON File

```json
{
  "url": "/test/test",
  "method": "GET",
  "response": {
    "status": 200,
    "type": "static",
    "file": "test-test.json",
    "format": "json", // "bytes"
    "content_type": "application/json"
  }
}
```

### Serve a PNG Image

```json
{
  "url": "/test/png",
  "method": "GET",
  "response": {
    "status": 200,
    "type": "static",
    "file": "dashboard.png",
    "format": "bytes",
    "content_type": "image/png"
  }
}
```

Any type of static file can be served this way by specifying the appropriate `file`, `format`, and `content_type` in the JSON configuration.
The `format` field defines how the file should be read. This means that JSON could be served with `format` `json` and `byts` the only difference is the validation. 
When `json` is specified the file is parsed to `json`, and validated at server startup. If the `json` is invalid the server shows the error at startup. 
Using `bytes` allows for the malformated `json` to be returned.


## Serve dynamic content

To serve dynamic content like random values or mapped values, the `server-file` needs to be defined. 
The suggested approach is to generate a `server-file` with a parser.

### Input JSON for Parser
 
The input JSON file can be any valid JSON object. For example:

```json
{
  "id": 1,
  "name": "Test Item",
  "price": 19.99,
  "inStock": true,
  "tags": ["test", "item"]
}
```

### Output JSON for Server

#### The output is meant to be modified by the user with mappings or other dynamically changing objects.

For static content to be served, you can just provide a file to be served. [More info](#serve-static-content)

After parsing, the output JSON will be formatted to define a mock endpoint. 
For dynamic type parsing, the fields will contain randomly generated values based on their original types. 
For static-type parsing, the fields will retain their original values.
> By default arrays generate 3 elements

> Note: For arrays the static value will use the first value of the array

Detailed information about generating dyamic content is [here](dynamic_configuration.md)

## Serve Custom Content

`server-faker` allows you to serve custom responses using a plugin-based system. 
Users can define custom logic and compile it into a plugin to be used by `server-faker`.

### Plugin Interface Requirement

To create a custom plugin, it must fulfil the following interface requirements:

```go
type PlgHandler interface {
	Respond(c *gin.Context)
}
```

### Building the Plugin

After defining the plugin, it needs to be built as a plugin using the Go compiler:

```sh
go build -buildmode=plugin -o {plugin-name}.so .
```

### Using the Plugin in Endpoints

To use the custom plugin, specify the path to the plugin file relative to the `server-faker` configuration file.

#### Example: Custom Logger Plugin

Below is an example of a custom logger plugin. 
This plugin logs information such as IP address, HTTP method, payload, and query parameters.

### Custom Logger Plugin Code

The code for a logger is in the [folder](plugins/logger). 

Follow to this directory and build the plugin:

```sh
go build -buildmode=plugin -o endpoint_logger.so logger_plugin.go
```

### Using the Plugin in the Endpoint Configuration

Specify the custom plugin in your endpoint configuration:

```json
{
  "endpoints": [
    {
      "method": "POST",
      "url": "/*",
      "response": {
        "type": "custom",
        "file": "endpoint_logger.so"
      }
    }
  ]
}
```

### Example Endpoint

When accessing the endpoint with a `POST` request to any URL path, the server will load and execute the custom plugin `endpoint_logger.so`. 
This plugin will log the request information (IP address, HTTP method, payload, and query parameters) and respond with a 200 status code confirming that the request has been logged.


# Creating Proxy Endpoint

`server-faker` allows you to create an endpoint that proxies requests to a specified URL using a given HTTP method.
You can specify URL parameters, query parameters, and payload using the available dynamic configurations.

### Example

Below is a complete example of creating a proxy endpoint:

```json
{
  "endpoints": [
    {
      "url": "/test/proxy",
      "method": "GET",
      "proxy": {
        "url": "http://localhost:8080/:url-key",
        "method": "POST",
        "type": "dynamic",
        "url_params": [
          {
            "key": "url-key",
            "random": {
              "type": "string-all",
              "min": 2,
              "max": 4
            }
          }
        ],
        "query_params": [
          {
            "key": "query-param",
            "random": {
              "type": "string-all",
              "min": 2,
              "max": 4
            }
          },
        ],
        "content_type": "application/json",
        "headers": {
          "test": "test"
        }
      }
    }
  ]
}
```

### Explanation

- `url`: Specifies the URL path for the endpoint in `server-faker`.
- `method`: Specifies the HTTP method for the endpoint (`GET` in this example).
- `proxy`: Defines the proxy configuration.
  - `url`: The URL to which the request will be proxied. You can include dynamic URL parameters prefixed with `:` (e.g., `:key1`).
  - `method`: The HTTP method to be used for the proxied request (`POST` in this example).
  - `type`: Specifies the type of proxy configuration. Here, it's `dynamic`.
  - `url_params`: Defines the dynamic URL parameters to be included in the proxied URL.
  - `query_params`: Defines the dynamic query parameters to be included in the proxied URL.
  - `content_type`: Specifies the content type of the request payload (e.g., `application/json`).
  - `headers`: Specifies additional headers to be included in the proxied request.

### Example Request Flow

1. **Client Request:** The client sends a `GET` request to `/test/proxy`.
2. **Proxy Logic:** The `server-faker` server processes the request and proxies it to `http://localhost:8080/:key1` using the `POST` method.
3. **Dynamic Values:** The `key1` parameter in the URL and the `key1` and `key2` query parameters are dynamically generated based on the specified random configurations.
4. **Headers and Content-Type:** The request includes the specified headers and content type.

### Generated Request Example

If the dynamic values are generated as follows:
- `url-key` URL parameter: `AB`
- `query-key` query parameter: `CD`

The proxied request would look like this:

- **URL:** `http://localhost:8080/AB?key1=CD`
- **Method:** `POST`
- **Headers:** `{"test": "test"}`
- **Content-Type:** `application/json`


## Next Steps

The work on `server-faker` is still in progress. 
Here are some plans for further development. This list is not in any particular order:

- **Use Multiple Files for Reading `server-faker` Endpoints:** Enhance the functionality to allow reading and combining multiple JSON files for endpoint configurations.
- **Incorporate `gofakeit`:** Integrate the `gofakeit` library to generate more meaningful and varied random values for mocked responses.
- **Enable Preparing `server-faker` Configuration from Swagger Documentation:** Implement functionality to automatically generate `server-faker` configuration files from existing Swagger documentation, simplifying the setup process for users.
- **Add Payload and URL Validators with Custom Responses:** Introduce validators to check the payload and URL parameters, allowing for custom responses based on validation results.
- **Add a System to Support Customizable Global API Events:** Develop a system to support global API events such as random response time delays or occasionally throwing errors, providing a more realistic testing environment.

---

Stay tuned for these and other exciting features in future releases!
If you like the project, consider supporting.

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/vimekgo) 