# Dynamic Configuration Options

There are several options to choose from to generate the mocked response:
- [static](#static-value)
- [random](#random-value)
- [mapped](#mapped-value)
  - [mapping from payload](#mapping-from-payload)
  - [mapping from URL](#mapping-from-url)
  - [mapping from query](#mapping-from-query)
- [array](#array-value)

## Static value

It is possible to return a static value in the response for a single value. 
This feature allows users to specify a fixed response value for specific keys within the response object. 

### Example Configuration

To return a static value, the JSON configuration should include the `static` key with the desired static value.

```json
{
  "endpoints": [
    {
      "method": "GET",
      "url": "/api/static-value",
      "response": {
        "status": 200,
        "body": {
          "object": [
            {
              "static": {
                "value": "test-string"
              }
            }
          ]
        }
      }
    }
  ]
}
```
### Example Request and Response
Request

```
GET request to /api/static-value
```

Response

```json
"test-string"
```
This example returns just a text `"test-string"` in the response. 

### Use Case

This feature is useful when you want certain fields in the response to always return a consistent value, regardless of any dynamic or randomized data generation for other parts of the response.

## Random value

It's possible to generate random values in the response given some options. 
This allows for more dynamic and varied responses based on specified criteria. 

### Random Value Types

- `string-numeric`: Represents a string that consists only of numeric characters (0-9).
- `string-uppercase`: Represents a string that consists only of uppercase alphabetic characters (A-Z).
- `string-lowercase`: Represents a string that consists only of lowercase alphabetic characters (a-z).
- `string-uppercase-number`: Represents a string that consists of a mix of uppercase alphabetic characters (A-Z) and numeric characters (0-9).
- `string-lowercase-number`: Represents a string that consists of a mix of lowercase alphabetic characters (a-z) and numeric characters (0-9).
- `string-all`: Represents a string that consists of a mix of uppercase and lowercase alphabetic characters (A-Z, a-z) and numeric characters (0-9).
- `integer`: Represents an integer value.
- `float`: Represents a floating-point number.
- `boolean`: Represents a boolean value (true or false).

### Example Configuration

Random values use the `max` and `min` keys to generate the response of a given length or within a specific range.

An example of configuring a random value in the response:

```json
{
  "endpoints": [
    {
      "method": "GET",
      "url": "/api/random-value",
      "response": {
        "status": 200,
        "body": {
          "object": [
            {
              "key": "random",
              "random": {
                "type": "string-all",
                "min": 6,
                "max": 8
              }
            }
          ]
        }
      }
    }
  ]
}
```

### Example Request and Response
Request

```
GET request to /api/random-value
```

Response

```json
{
  "random": "a3BdE7"
}
```

## Mapped value

There are 3 options for mappings:
- [body (payload)](#mapping-from-payload)
- [url](#mapping-from-url)
- [query params](#mapping-from-query)

All of the mappings (query, url, and payload) allow users to convert between data types. Specifically, you can convert:
- From a string to a number
- From a number to a string
More about conversions here 

> It is possible to convert the mapped `string` or `integer` from `url`, `body` and `query` to an `integer` or `string` in the response. 
To so add a `as` keyword to a json specifiing the output. 

### Mapping from Payload

It is possible to map keys from the incoming request body to the response. 
This feature utilizes [JSONPath](https://jsonpath.com/) to locate elements within the request body and return the appropriate value. The value could be a string, array, or object mapped to a certain key in the response.

### How It Works

The `server-faker` allows you to define mappings within your JSON configuration file. 
These mappings specify which elements from the request body should be included in the response. 
The `from` attribute indicates the source of the data, and the `path` attribute specifies the JSONPath expression to locate the desired element in the request body.

### Example Configuration

Here's a minimal example demonstrating how to map a key from the request body to the response:

```json
{
  "endpoints": [
    {
      "method": "POST",
      "url": "/api/map",
      "response": {
        "status": 200,
        "body": {
          "object": [
            {
              "key": "mapped",
              "mapped": {
                "from": "body",
                "path": "$.key"
              }
            }
          ]
        }
      }
    }
  ]
}
```

### Example Request and Response
Request

```json
{
  "key": "value to be mapped"
}
```

Response

```json
{
  "mapped": "value to be mapped"
}
```

In this example, the key from the request body is mapped to the mapped field in the response body. 
The JSONPath expression $.key locates the key in the request body and maps its value to the mapped field in the response.

### Additional Notes

- JSONPath expressions must be correctly formatted to ensure the desired elements are accurately located within the request body.

## Mapping from URL

It is possible to use parts of the URL to be mapped in the response. 
This requires adding a parameter in the URL prefixed with `:` matching the given mapping.

### How It Works

`server-faker` allows you to define URL parameters within your JSON configuration file. 
These parameters can be referenced in the response by specifying the `from` attribute as `url` and providing the corresponding `key`. 
This enables dynamic responses based on the values extracted from the URL.

### Example Configuration

Here's a minimal example demonstrating how to map a part of the URL to the response:

```json
{
  "endpoints": [
    {
      "method": "GET",
      "url": "/test/test/random/:param",
      "response": {
        "status": 200,
        "body": {
          "object": [
            {
              "key": "url-key",
              "mapped": {
                "from": "url",
                "param": "param"
              }
            }
          ]
        }
      }
    }
  ]
}
```

### Example Request and Response
Request

```
GET request to /test/test/random/exampleValue
```

Response

```json
{
  "url-key": "exampleValue"
}
```

In this example, the param from the URL `/test/test/random/exampleValue` is mapped to the mapped field in the response body. 
The URL parameter :param is extracted and its value is used in the response.

### Additional Notes

- URL parameters must be prefixed with `:` in the endpoint URL and match param provided in URL.

## Mapping from Query

It is possible to map query parameters from the incoming request URL to the response object keys.
This allows you to create dynamic responses based on the values of query parameters.
Additionally, you can parse arrays from query parameters by specifying the `index` to locate the desired element.

### How It Works

`server-faker` allows you to define mappings within your JSON configuration file. 
These mappings specify which query parameters should be included in the response. 
The `from` attribute indicates the source of the data as `query`, and the `key` attribute specifies the query parameter to map. 
If the query parameter is an array, you can use the `index` attribute to specify the position of the element you want to map.

### Example Configuration

Here's a minimal example demonstrating how to map a query parameter to the response:

```json
{
  "endpoints": [
    {
      "method": "GET",
      "url": "/api/query",
      "response": {
        "status": 200,
        "body": {
          "object": [
            {
              "key": "query-key",
              "mapped": {
                "from": "query",
                "key": "test"
              }
            }
          ]
        }
      }
    }
  ]
}
```
### Example Request and Response
Request
```
GET request to /api/query?test=exampleValue
```
Response

```json
{
  "query-key": "exampleValue"
}
```

In this example, the test query parameter from the URL `/api/query?test=exampleValue` is mapped to the mapped field in the response body.
The query parameter test is extracted and its value is used in the response.

### Additional Notes

- Query parameters should be specified in the request URL.
- The index attribute is used to specify the position of the element in the array (0-based index).

## Type Conversions

It is possible to convert all types of mappings (`query`, `url`, and `payload`) to a specific type, either `integer` or `string`.
This ensures that the values in the response are correctly formatted as needed.

### How It Works

To add a conversion, the optional key `"as"` needs to be specified with one of the two options: `number` or `string`. This key indicates the desired type to which the mapped value should be converted.

### Example Configuration

Here's an example demonstrating how to apply a type conversion to a query parameter mapping:


```json
{
  "endpoints": [
    {
      "method": "GET",
      "url": "/api/convert",
      "response": {
        "status": 200,
        "body": {
          "object": [
            {
              "key": "converted-key",
              "mapped": {
                "from": "query",
                "key": "value",
                "as": "number"
              }
            }
          ]
        }
      }
    }
  ]
}
```

### Example Request and Response
Request
```
GET request to /api/convert?value=123
```
Response

```json
{
  "converted-key": 123
}
```

In this example, the value query parameter from the URL `/api/convert?value=123` is converted to a number in the response body.

### Handling Conversion Errors

If the conversion fails (for example, if a non-numeric string is provided where a number is expected), a conversion error will be returned. 
This helps ensure that the API mock behaves predictably and provides meaningful feedback when invalid data is encountered.

Request
```
GET request to /api/convert?value=abc
```
Response
```
{
	"title": "Conversion failed",
	"url": "/api/convert",
	"errors": [
		{
			"details": "param key: [number]: not known conversion for 'abc' to number strconv.ParseFloat: parsing \"abc\": invalid syntax: failed converting param"
		}
	]
}
```

### Mapping Arrays from Query Parameters

If the query parameter is an array, you can specify the index attribute to map a particular element from the array.

#### Configuration File for Array Mapping

```json
{
  "endpoints": [
    {
      "method": "GET",
      "url": "/api/query-array",
      "response": {
        "status": 200,
        "body": {
          "object": [
            {
              "key": "query-array-key",
              "mapped": {
                "from": "query",
                "key": "items",
                "index": 1
              }
            }
          ]
        }
      }
    }
  ]
}
```

#### Example Request and Response for Array Mapping
Request
```
GET request to /api/query-array?items=first&items=second&items=third
```
Response

```json
{
  "query-array-key": "second"
}
```
In this example, the items query parameter is an array with values ["first", "second", "third"]. 
The index attribute specifies that the second element ("second") should be mapped to the mapped field in the response body.

## Array value

All mentioned above elements can be used to generate an array. 
The array can consist of any simple type like string, integer or more complex objects. 
Arrays are generated with a length between the specified `min` and `max` values, using the elements specified in the `element` array.

### Example: Generating an Array with Objects

Here is an example of generating an array with integers:

```json
{
  "endpoints": [
    {
      "method": "GET",
      "url": "/api/array-integers",
      "response": {
        "status": 200,
        "body": {
          "object": [
            {
              "array": {
                "min": 3,
                "max": 3,
                "element": [
                  {
                    "random": {
                      "type": "integer",
                      "min": 1,
                      "max": 100
                    }
                  }
                ]
              }
            }
          ]
        }
      }
    }
  ]
}
```
