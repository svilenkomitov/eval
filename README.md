# eval

Implement a web service that will evaluate simple math word problems, with the
following endpoints:

## POST /evaluate

- Parse and evaluate a problem expression returning the answer as an integer.
- Expressions with no operations simply evaluate to the number given:
  - "What is 5?" should return 5.
- Simple arithmetic:
  - Add two numbers together
    - What is 5 plus 13?
  - Subtraction
    - What is 7 minus 5?
  - Multiplication
    - What is 6 multiplied by 4?
  - Division
    - What is 25 divided by 5?
- Handle a set of operations, in sequence.
  - Since these are verbal word problems, evaluate the expression from
    left-to-right, ignoring the typical order of operations.
    - What is 3 plus 2 multiplied by 3? 15 (i.e. not 9)
- The service should return error on:
  - Unsupported operations ("What is 52 cubed?")
  - Non-math questions ("Who is the President of the United States?")
  - Expressions with invalid syntax ("What is 1 plus plus 2?")
- The endpoint should accept json data in the following format: `{"expression":"<a
  simple math problem>"}`
- The endpoint should return json data in the following format: `{"result":<the
  expression's result>}`

## POST /validate

- Validate an expression, return if valid or not.
- The endpoint should accept json data in the following format: `{"expression":"<a
  simple math problem>"}`
- The endpoint should return json data in the following format:
  - If the expression is not valid:
  
  ```
  {
  "valid":false,
  "reason":"<the reason why the expression is invalid>"
  }
  ```
  - If the expression is valid: `{"valid":true}`

## GET /errors

- Return all errors occurred with their frequency (e.g. if the same unsupported
  expression is submitted many times), URL of the endpoint on which they have
  occurred (i.e. /evaluate or /validate) and the type of the error (e.g. Unsupported
  operations).
- The endpoint should return json data in the following format:
  ```
  [
    {
      "expression": "<a simple math problem>",
      "endpoint": "<endpoint URL>",
      "frequency": <number of times the expression failed>,
      "type": "<error type>"
    }
  ]
  ```
  
Please note that for the purpose of the task it is fine for the web service to
persist all data in memory.

### Optional
- Implement CLI client that will provide the same set of functionality and will use
  the web service as its backend.

## Installation

* add env values in ```docker-compose.yaml```
* install the service
    ```bash
    $ docker-compose up --build -d
    ```