I need to understand our e2e testing framework for the serverless controller to implement comprehensive functional tests. Help me analyze and understand:

FRAMEWORK STRUCTURE:
1. How does our BDD-style testing framework work? (Given/When/Then pattern)
2. What are the key testing utilities and helper functions available?
3. How do I set up test environments and namespaces?
4. What's the pattern for uploading test metadata to JFrog?
5. How do I create and manage PaaS API resources in tests?

KUBERNETES INTEGRATION:
1. How do tests interact with real Kubernetes clusters?
2. What utilities exist for waiting on pod states, services, deployments?
3. How do I verify namespace creation/deletion in tests?
4. What's the pattern for checking Knative services and revisions?
5. How do I validate ArgoCD application states?

SERVERLESS CONTROLLER SPECIFIC:
1. How do I test the complete container app lifecycle (create → deploy → update → delete)?
2. What's the pattern for testing NaaS integration and namespace management?
3. How do I test different deployment strategies (blue-green, canary)?
4. How do I simulate service failures (NaaS down, JFrog auth failures)?
5. How do I test multi-region deployments and cloud provider differences?

ERROR TESTING:
1. How do I simulate and test failure scenarios?
2. What's the pattern for testing authorization and permission failures?
3. How do I test resource quota violations and constraint failures?
4. How do I verify proper error messages and status codes?

TEST ORGANIZATION:
1. How should I structure test files for different feature areas?
2. What's the naming convention for test scenarios?
3. How do I share test data and utilities across test files?
4. What's the pattern for test cleanup and resource management?

DEBUGGING AND ANALYSIS:
1. How do I debug failing tests and inspect cluster state?
2. What logging and tracing is available during test execution?
3. How do I capture test artifacts for failure analysis?
4. What tools are available for performance testing and benchmarking?

Show me code examples and patterns from our existing test suite. Focus on practical implementation details I can use to write robust functional tests for the serverless controller.
