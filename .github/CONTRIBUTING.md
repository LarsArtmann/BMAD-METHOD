# Contributing to BMAD Method

Thank you for your interest in contributing to BMAD Method! 🎉

## Code of Conduct

This project adheres to a Code of Conduct. By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Issues

Before creating an issue, please:
1. Search existing issues to avoid duplicates
2. Use the appropriate issue template
3. Provide clear, detailed information

### Suggesting Features

We love feature suggestions! Please:
1. Use the feature request template
2. Explain the problem you're solving
3. Describe your proposed solution
4. Consider the impact on different tiers

### Code Contributions

#### Development Setup

1. **Prerequisites**
   ```bash
   go version  # 1.21+
   node --version  # 20+
   npm install -g @typespec/compiler
   ```

2. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/BMAD-METHOD.git
   cd BMAD-METHOD
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   npm install
   ```

4. **Run Tests**
   ```bash
   go test ./...
   bash scripts/validate-schemas.sh
   ```

#### Making Changes

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**
   - Follow existing code style
   - Add tests for new functionality
   - Update documentation as needed

3. **Test Your Changes**
   ```bash
   # Run all tests
   go test ./...
   
   # Test generated projects
   go build -o bmad-cli ./cmd/generator
   ./bmad-cli generate --name test-project --tier basic
   cd test-project && go mod tidy && go build ./cmd/server
   ```

4. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

#### Code Style Guidelines

- **Go**: Follow standard Go conventions (`gofmt`, `golint`)
- **TypeScript**: Use Prettier and ESLint configurations
- **Templates**: Keep templates readable and well-commented
- **Documentation**: Use clear, concise language

#### Testing Requirements

- **Unit Tests**: All new functions should have unit tests
- **Integration Tests**: Test generated projects compile and run
- **BDD Tests**: Update feature files for new functionality
- **Manual Testing**: Test CLI interactions manually

#### Pull Request Process

1. **Create PR**: Use the PR template and fill all sections
2. **Review Process**: Address feedback promptly
3. **CI/CD**: Ensure all checks pass
4. **Merge**: Maintainers will merge after approval

## Project Structure

```
BMAD-METHOD/
├── cmd/generator/          # CLI application
├── pkg/                   # Core packages
│   ├── config/           # Configuration handling
│   ├── generator/        # Template generation engine
│   └── schemas/          # TypeSpec schemas
├── template-health/       # Health template files
├── templates/            # Tier templates
├── features/             # BDD test features
├── docs/                 # Documentation
└── .github/              # GitHub templates
```

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):
- **Major (X.0.0)**: Breaking changes
- **Minor (0.X.0)**: New features (backward compatible)
- **Patch (0.0.X)**: Bug fixes

### Release Schedule

- **Patch releases**: As needed for critical bugs
- **Minor releases**: Monthly with new features
- **Major releases**: Quarterly or for significant changes

## Recognition

Contributors are recognized in:
- Release notes
- README contributors section
- Annual contributor spotlight

### Contributor Levels

- **Contributor**: Made valuable contributions
- **Regular Contributor**: Multiple contributions over time
- **Core Contributor**: Significant ongoing contributions
- **Maintainer**: Project maintenance responsibilities

## Getting Help

- **Discord**: Join our community server
- **Discussions**: GitHub Discussions for questions
- **Issues**: GitHub Issues for bugs/features
- **Email**: maintainers@bmad-method.dev

## Areas for Contribution

### High Priority
- Template improvements
- Documentation enhancements
- Bug fixes
- Performance optimizations

### Medium Priority
- New tier features
- CLI enhancements
- Test coverage improvements
- Integration examples

### Good First Issues
- Documentation fixes
- Template typos
- Simple CLI improvements
- Test additions

## Technical Decisions

Major technical decisions are documented in:
- Architecture Decision Records (ADRs)
- Design discussions in GitHub
- Technical RFCs for major changes

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

---

**Happy Contributing! 🚀**

For questions, reach out in our Discord or create a discussion!