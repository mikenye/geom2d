# Contributing to Geom2D

Thank you for considering contributing to Geom2D! We welcome contributions of all kinds, including bug fixes, new features, documentation improvements, and more.

## How to Contribute

There are several ways to contribute to Geom2D:

- **Report Bugs**: If you find a bug, please [open an issue](https://github.com/mikenye/geom2d/issues/new/choose).
- **Suggest Features**: Have an idea for a new feature? Before starting any work, please [open a discussion](https://github.com/mikenye/geom2d/discussions/new) or comment on an existing issue to ensure the feature aligns with the project goals.
- **Improve Documentation**: Spot a typo or think the docs could be clearer? Feel free to submit a pull request.
- **Write Code**: Implement new features, fix bugs, or improve the codebase.

## Code of Conduct

Geom2D has adopted the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/) to foster an open, welcoming, and inclusive community.  
Please review the code of conduct and adhere to its guidelines when interacting with the project and its community.

## Development Setup

To get started, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/mikenye/geom2d.git
   ```
2. Navigate to the project directory:
   ```bash
   cd geom2d
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run tests to verify the setup:
   ```bash
   go test ./...
   ```

You're all set!

## Submitting Changes

1. **Fork the repository** and create a branch:
   ```bash
   git checkout -b your-feature-branch
   ```
2. **Make your changes** and commit them:
   ```bash
   git commit -m "Add a meaningful commit message"
   ```
3. **Push your branch** to GitHub:
   ```bash
   git push origin your-feature-branch
   ```
4. Open a pull request (PR) against the `main` branch.

Please ensure your changes:
- Pass all tests (`go test ./...`) to ensure your changes don’t break existing functionality.
- Include updates to documentation if necessary.
- Are formatted using `go fmt`.

## Style Guide

- Follow idiomatic Go practices.
- Use `go fmt` to format your code before submitting.
- Write clear and concise commit messages.
- Add tests for any new features or bug fixes.
- Any new public functions should include an `Example` function to demonstrate their usage.

## Thank You

Thank you for taking the time to contribute to Geom2D! Your efforts make the project better for everyone.
