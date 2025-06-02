# Pull Request

## 📋 Description

<!-- Provide a brief description of the changes in this PR -->

## 🔗 Related Issue(s)

<!-- Link to the issue(s) this PR addresses -->
Fixes #<!-- issue number -->
Closes #<!-- issue number -->
Related to #<!-- issue number -->

## 🎯 Type of Change

<!-- Mark the relevant option with an "x" -->

- [ ] 🐛 Bug fix (non-breaking change that fixes an issue)
- [ ] ✨ New feature (non-breaking change that adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📚 Documentation update
- [ ] 🔧 Refactoring (no functional changes)
- [ ] ⚡ Performance improvement
- [ ] 🧪 Test addition or modification
- [ ] 🔐 Security improvement
- [ ] 🎨 UI/UX improvement
- [ ] 🔨 Build/tooling change

## 🎯 Affected Tiers

<!-- Mark all tiers affected by this change -->

- [ ] Basic
- [ ] Intermediate  
- [ ] Advanced
- [ ] Enterprise
- [ ] CLI Tool
- [ ] Documentation

## 🧪 Testing

<!-- Describe the tests you ran and provide instructions for reviewers -->

### Test Environment
- OS: <!-- macOS/Linux/Windows -->
- Go Version: <!-- e.g., 1.21 -->
- Node Version: <!-- e.g., 20 (if applicable) -->

### Tests Performed
- [ ] Unit tests pass (`go test ./...`)
- [ ] Integration tests pass
- [ ] BDD tests pass (`go test ./features/...`)
- [ ] Generated projects compile successfully
- [ ] CLI interactive mode works
- [ ] TypeSpec schemas validate
- [ ] Manual testing performed

### Test Commands
```bash
# Commands used to test the changes
go test ./...
./scripts/validate-schemas.sh
```

## 📋 Checklist

<!-- Mark completed items with an "x" -->

### Code Quality
- [ ] Code follows project style guidelines
- [ ] Self-review of code completed
- [ ] Code is properly commented (especially complex areas)
- [ ] No debugging artifacts left in code
- [ ] Error handling is appropriate

### Testing
- [ ] Tests added/updated for new functionality
- [ ] All existing tests pass
- [ ] Edge cases considered and tested
- [ ] Performance impact assessed

### Documentation
- [ ] Documentation updated (if applicable)
- [ ] README updated (if applicable)
- [ ] API documentation updated (if applicable)
- [ ] Examples updated (if applicable)

### Compatibility
- [ ] Changes are backward compatible
- [ ] Breaking changes documented
- [ ] Migration guide provided (if needed)

## 📸 Screenshots/Examples

<!-- If applicable, add screenshots or examples showing the changes -->

### Before
<!-- Screenshots or examples showing current behavior -->

### After
<!-- Screenshots or examples showing new behavior -->

## 🔄 Generated Project Testing

<!-- For changes affecting code generation -->

### Generated Projects Tested
- [ ] Basic tier project
- [ ] Intermediate tier project
- [ ] Advanced tier project
- [ ] Enterprise tier project

### Features Tested
- [ ] TypeScript client generation
- [ ] Kubernetes manifests
- [ ] Docker configuration
- [ ] Health endpoints
- [ ] OpenTelemetry integration (if applicable)

## 📝 Additional Notes

<!-- Any additional information that reviewers should know -->

### Performance Impact
<!-- Describe any performance implications -->

### Security Considerations
<!-- Describe any security implications -->

### Breaking Changes
<!-- List any breaking changes and migration steps -->

## 🤝 Reviewer Notes

<!-- Information specifically for reviewers -->

### Focus Areas
<!-- Areas where you'd like specific feedback -->

### Known Issues
<!-- Any known issues or limitations -->

---

**Thank you for contributing to BMAD Method! 🙏**

By submitting this PR, you agree that your contributions will be licensed under the project's license.