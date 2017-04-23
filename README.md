# Testing External Services

This is sample implementation of a service/tool that uses the Google Search API.

The intention here is not to explain how to use the API but:
* how to implement a package that provides a good abstraction of an external service 
* how to test how your code interacts with the external service
* how to ensure the external service does what you need it to do

## Where to start?

* If you're interest take a quick look at the "server" [here](main.go).
* There is not much to it, it just takes a single search parameter and use it to call an external service (Google in this case)
* Other than it's ugliness, there is no much of significance.

## Directory Structure

Personally I prefer to keep all "external services" under 1 parent package called "[external](external/)".

This is totally not necessary but I find it feels organised to me.

Under this external package I will make 1 package per service.

In this example I have also used the ["internal"](https://golang.org/cmd/go/#hdr-Internal_Directories) so that none of 
the external service's formats and constants leak from this package.

### API
[Code](external/search/api.go)

This file contains:
* The API that is exported by the package
* The standard errors returned by the package

Points of interest:
* I am exporting an interface instead of struct
* I am intentionally not allowing any of the implementation details of the external service to leak from this package

## Implementation
[Code](external/search/implementation.go)

Many people recommend not using constructors, exporting the implementation type and leaving the mocking to the
user.  
I do not.  
I find that exporting only the API and providing mocks, I both force myself to use the interface
(which results in better encapsulation) and providing the mocks reduces duplication.

This file contains the default "live" implementation of the client.

## Unit tests
[Code](external/search/implementation_unit_test.go)

In this example, I have not added any unit tests.
Most, if not all of the implementation could be considered "too simple to break" and all of it is pretty well
covered in UAT tests.

This is not always the case and you are welcome to disagree with my definition of simple.
I always recommend you add enough tests to make yourself feel comfortable (and not more).
Remember, too many tests does have a (maintenance) cost too.

## UAT tests
[Code](external/search/implementation_uat_test.go)

These tests verify our API contract with our users.

They are designed to ensure that the code of this package is doing what we want it to independent of the external
service.

They go a long way to preventing regression.

We do not consider them E2E tests as they do not make calls to the external service.

UAT tests like these are great for testing failures of the external service that are otherwise difficult/impossible to 
test.

These tests also provide good contrast with the E2E tests.
This is important because if these tests are working and the E2E tests break, then this most likely means the
external service is no longer doing what you need it too (i.e. their API contract changed) or the configuration is
incorrect.

Personally I find this results in faster identification/debugging of issues.

## E2E tests
[Code](external/search/implementation_e2e_test.go)

These tests ensure that this package provides the functionality it intends to.

They make "real" calls to the Google API as such they require and API key and custom search engine id to be
supplied as environment variables ("KEY" and "CX" respectively)

These tests indirectly confirm that the external service performs in a manner expected by this package.

## Mock
[Code](external/search/mock_client.go)

This is a generated mock implementation of the client (generated with [Mockery](https://github.com/vektra/mockery)).

This is provided to make it easy for users of this package to test their usage of this package.

## Notes:
* For the tests to work you need to obtain an Google Search API key from https://developers.google.com/custom-search/json-api/v1/introduction and you will need to set environment vars for:
	* KEY - this is the API key
	* CX - this is a custom search engine ID (I created '010104945395156267508:uhan0jxej6g' I am not sure this will work for you)
* To run the sample type `KEY=xxx CX=xxx go run main.go` then open your browser to `http://localhost:8080/?q=mysql`
