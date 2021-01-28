# Test Data

There seems to be some problems with linking the `TestAssembly.asmdef` with
the `Google.Protobuf` namespace. We're using this hack where test data is
technically defined in the app space, and referencing it in our tests.

Clients should not use this data.