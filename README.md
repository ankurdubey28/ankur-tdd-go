## What is TDD (Test Driven Development)?

**TDD (Test Driven Development)** is a software development process where **tests are written before the actual implementation code**.  
The goal is to ensure that the implemented code satisfies the required behavior through **short feedback loops**.

TDD follows agile development methodology , where feedback loop is integrated fast and development is done iteratively.

### TDD Cycle: Red → Green → Refactor

1. **Red**
   - Write a test that defines a desired feature or behavior.
   - The test should fail initially because the feature is not yet implemented.

2. **Green**
   - Write the **minimum amount of code** required to make the test pass.

3. **Refactor**
   - Improve and optimize the code structure without changing its behavior.
   - Ensure that all tests **continue to pass** after refactoring.

<img alt="img.png" height="500" src="img.png" width="500"/>


## Best Practices for TDD
1. Write atomic tests - each test should focus on single functionality / behavioural aspect. So each test should be kept small.
2. Write simplest tests first -  begin by writing simplest test case which will fail.
3. Write test case for edge cases - consider boundary conditions of input and add test case for them, often bugs are caught here.
4. Refractor regularly - after test case passes, take time to refractor code and improve its design without change in behaviour.
5. Automation - use test automation tools to fast track the process for testing.
6. Follow Red-Green-Refractory cycle.
7. Maintain fast feedback loop - test suite should execute immediately to receive feedback and accordingly code can be refactored, allowing faster development.
8. Continuous test - integrate test with CI/CD pipelines to automatically execute test whenever code changes are made to ensure early bug detection and feedback addressal.


# Unit Test in Go
1. It needs to be in a file with a name like xxx_test.go
2. The test function must start with the word Test
3. The test function takes one argument only t *testing.T
4. To use the *testing.T type, you need to import "testing", like we did with fmt in the other file


# Side Effects and Randomness In Unit tests
1. Testing in presence of side effects and randomness is not advisable , because outputs are non deterministic
2. Side effects means function changes or depends on external state (api / input / databases) outside functions scope.
3. Solution? Dependency Injection + Determinism

# Dependency Injection (DI)
1. Programming technique where a construct is passed (injected) to another construct that depends on it.
2. DI also helps in decoupling code , allowing for better unit test.


# MOCKING
1. Technique used to isolate a unit being tested by replacing dependencies with objects we can control. This can be done using DI (Dependency Injection).
2. Dependency here being used by the unit can be anything , but generally its the a module imported by unit.
3. Why Mocking?
   - Individual units depends on external apis , db , services , whose behaviours is not in our control, making testing tough.
   - Thus mocking introduces determinism , speed (fast testing) , isolation (important for individual modules to test) and give control over scenario.

# Test Doubles
1. A test double is to test objects , just how stunt double is to actor. Any object which can stand in place of actual object being tested is test double.

- **Dummy**
   - Used only to satisfy interface.
   - Mostly used as a placeholder to make code compile and compiler happy.

- **Stub**
   - Returns hardcoded data. Used to control the output.

- **Mocks**
   - Very similar to stubs, but they are interaction heavy compared to stubs.
   - Means we do not just expect mocks to return data, but also assume specific order of steps to be executed.

- **Fake**
   - These objects have complete working implementation in them, but the implementation provided is not exactly production grade, just something enough to allow test to pass.
   - Example → In-memory database.

- **Spies**
   - Spies are stubs which also record the flow of action and information.