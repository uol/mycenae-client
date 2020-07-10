# Mycenae Client

A client built to call the API services from the timeseries database [Mycenae](https://github.com/uol/mycenae).

### How to build:
```sh
  $ go mod vendor
  $ go build
```

### How to execute the tests:
All tests are located in the directory "tests". They can be executed in the main suite file called "suite_test.go" or executing the script "run-tests.sh" located in the project root.
Ex:
```sh
  $ ./run-tests.sh
```