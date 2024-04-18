# skeleton-go-api

This is my take of skeleton project for a Go API application. It is intended to be used as a starting point for new API applications.

## Usage

To use this project as a starting point for a new API application, follow these steps:
1. Clone this repository.
2. Rename the `skeleton-go-api` directory to the name of your new application.
3. Update the `go.mod` file to reflect the new module name.
4. Replace `skeleton-go-api` with the new application name in all files.
5. Implement the desired functionality for the new API application.

You can basically replace the `skeleton-go-api` with your new application name and start building your API application.

## What's Included

This project includes the following components:
- A basic API application structure.
- A simple command with a subcommand.
- Makefile for Build and test scripts.
- A `Makefile` with common tasks.
- A `Dockerfile` for building a Docker image.
- A GitHub Actions workflow for versioning
- A GitHub Actions workflow for PR checks. It will check linting, test coverage and build.
- golangci-lint for linting and static analysis.

## Sample Service, Get Photos from jsonplaceholder.typicode.com

This project includes a sample service that fetches photos from jsonplaceholder.typicode.com. The service is implemented in the `photos` package and is used by the `get` command to fetch.
```bash
make build
./skeleton-go-api
```

Now run `curl http://localhost:8080/photos/1` will return

```json
{"albumId":1,"id":1,"title":"accusamus beatae ad facilis cum similique qui sunt","url":"https://via.placeholder.com/600/92c952","thumbnailUrl":"https://via.placeholder.com/150/92c952"}
```

## Go Implementation Guidelines 

### TL;DR: Enhance flexibility and maintainability by:
- Testability: Use dependency injection and factory patterns for easier mocking and isolated testing. Aim for test coverage above 80%.
- Adhering to SOLID Principles: Extend functionalities without altering existing code to preserve functionality and ensure stability.
- Embracing DRY: Centralize logic to avoid duplications and inconsistencies.
- Separating Concerns: Keep modules focused on their intended functionality to simplify maintenance and reduce conflicts.
- Maintaining Code Hygiene: Standardize coding practices, keep methods concise, and reduce conditional complexity for readability and quicker onboarding.
By following these key strategies, we can build a codebase that's easier to manage, extend, and debug.

### Quick overview of implementation approach we want to take, basically summary of software engineering best practices 

#### Testability through Dependency Injection / Factory Pattern:
Dependency injection and the use of factory patterns are key to creating testable code. They allow you to easily mock dependencies in your tests, ensuring that units can be tested in isolation without relying on external systems or complex setup. This approach improves the reliability of your tests and makes them easier to write and understand.
- Action Point
    - Setup test coverage check. 80% is ideal.
    - Create a service and related methods, inject it to another service and accepts other services.
    - Accept interfaces, return structs - Go proverb.
    - Define interface at the consumer side, not at the provider side.

#### Preserving Existing Functionality (Open/Closed Principle):
The Open/Closed Principle, part of SOLID principles, suggests that software entities (classes, modules, functions, etc.) should be open for extension but closed for modification. This means you can add new functionality without changing existing code, thereby reducing the risk of introducing bugs into existing features. It promotes modularization and extensibility.
    The sub-points about reusability and avoiding modifications to integrate with new modules emphasize the importance of designing components that are flexible and can be easily reused, which is fundamental for reducing duplication and fostering a modular architecture.
- Action Point
    - Follow Point 1, 3, 4, it will be automatically given.
    - See if we can reuse module, when implementing a new package, rather than putting ad hoc changes in between other unrelated module.

#### DRY - Don't Repeat Yourself:
Avoiding repetitive logic is crucial for maintaining a codebase that's easy to update and debug. When changes are needed, having a single source of truth for any piece of logic means only needing to make that change in one place, reducing the risk of inconsistencies and bugs.
    
- Action Point
  - Check for duplicate logic e.g. if we have multiple places where we try to do same thing, it can be all done where it matters and in one place. It will help avoid unexpected values causing bugs.
      
#### Separation of Concerns:
Ensuring that a package or module only contains what its name suggests is about adhering to the principle of separation of concerns. Each component should have a single responsibility and not take on more than it should, which simplifies understanding and maintaining the code.
Moving code or configurations to more relevant areas when the context changes is part of keeping the codebase organized and intuitive. This adaptability is key to maintaining clarity and efficiency as the project evolves.
Immediate benefit will be less merge/rebase conflicts.
We can avoid unrelated logic mixed up in one place, related logics scattered around in different places and same checks being repeated. It was observed that when we add more code than needed in one spot, we end up repeating the same checks or validations in different places, because it easily made it confusing to figure out which part of the code is responsible for what, as our code grows.
- Action Point
    - Follow idiomatic Go patterns.
    - Create a package that does only what its names says. 
    - Let’s break nested package directory into flat structure as much as possible.

#### Code hygiene.
Following code standard will help all developers getting onboard easier to read code, thus speed up development process.
- Action point
    - setup golangci-lint in task test
    - Keep functions/methods small; If it does more than 3 operation you can make it into a function. There is no hard limit for how long. But if we have to set a line, it shouldn’t be longer than a screenful, which is around 30 lines. If it’s less than 30 lines, probably still better to split but it might be ok cases like struct assignments which might occupy more space, can be considered as a single logical line of code or simple initialisations. If a function contains distinctly different logic although no more than 30 lines, better to split. When a function gets long, this could easily mean it encompasses more logic than what the function name suggests.
    - Try to minimize the use of if statements as much as possible; having fewer code paths always reduces bugs.
