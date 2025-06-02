# Single-Line Installable CLI Tool Packaging Strategy

## Objective
Create a comprehensive packaging and distribution strategy that enables users to install and use the CLI tool with a single command across all major platforms, package managers, and deployment environments.

## Context
This prompt focuses on transforming a development tool into a production-ready, widely distributable CLI application that can be installed with minimal friction across different operating systems, package managers, and deployment scenarios.

## Task Description

### Phase 1: Multi-Platform Binary Distribution
1. **GitHub Releases Automation**
   - Set up GitHub Actions for automated cross-platform builds (Windows, macOS, Linux arm64/amd64)
   - Implement semantic versioning with git tags triggering releases
   - Create checksums and code signing for security verification
   - Generate release notes from commits and pull requests

2. **Smart Installation Script**
   - Create universal installer: `curl -sSL https://install.artmann.foundation?repo=template-health | bash`
   - Implement OS/architecture detection and appropriate binary download
   - Add installation verification, rollback capabilities, and error handling
   - Support custom installation directories and version pinning

3. **Binary Optimization**
   - Statically linked binaries with zero external dependencies
   - Compressed binaries under 50MB with fast startup (<2s)
   - Security features including code signing and checksum verification

### Phase 2: Package Manager Integration
1. **Native Package Managers**
   - **Homebrew**: `brew install vonArtmann/tap/template-health`
   - **Chocolatey**: Windows package manager integration
   - **Snap**: Linux universal package format
   - **APT/YUM**: Debian and RedHat repository setup

2. **Language Ecosystem Integration**
   - **NPM**: `npm install -g @vonArtmann/setup-cli`
   - **Go**: `go install github.com/LarsArtmann/BMAD-METHOD/cmd/template-health@latest`
   - **Nix**: `nix profile install github:LarsArtmann/BMAD-METHOD`
   - **PyPI**: Python wrapper for data science teams

3. **Nix Ecosystem Support**
   - Create Nix derivation with reproducible builds
   - Implement Nix flakes for modern workflow support
   - Build NixOS module for system-wide installation
   - Develop Home Manager module for user configuration

### Phase 3: Container and Cloud Distribution
1. **Container Images**
   - **Docker Hub**: `docker run --rm -v $(pwd):/workspace template-health/cli generate`
   - Multi-architecture builds (arm64/amd64) with Alpine/distroless base
   - Optimized images under 100MB with security scanning
   - Integration with GitHub Container Registry

2. **Cloud-Native Deployment**
   - **Kubernetes Helm Chart**: One-command cluster deployment
   - **Terraform Module**: Infrastructure as code integration
   - **AWS/GCP/Azure Marketplace**: Cloud provider integration
   - **Serverless**: Function-as-a-Service deployment options

### Phase 4: Enterprise and CI/CD Integration
1. **CI/CD Platform Integration**
   - **GitHub Actions**: Marketplace action for workflow integration
   - **GitLab CI**: Component for GitLab ecosystem
   - **Jenkins**: Plugin for enterprise environments
   - **Azure DevOps**: Extension for Microsoft ecosystem

2. **Enterprise Features**
   - **RBAC Integration**: Role-based access control
   - **Audit Logging**: Compliance and security tracking
   - **Air-gapped Installation**: Offline deployment support
   - **LDAP/SSO**: Enterprise authentication integration

## Technical Implementation

### Installation Script Structure
```bash
#!/bin/bash
# Universal installer for template-health
# Usage: curl -sSL https://install.artmann.foundation?repo=template-health | bash

detect_os_arch() {
    # Smart OS and architecture detection
}

download_binary() {
    # Secure binary download with verification
}

verify_installation() {
    # Installation verification and testing
}

main() {
    # Main installation flow with error handling
}
```

### Package Manager Configurations
1. **Homebrew Formula** (`vonArtmann/tap/template-health.rb`)
2. **NPM Package** (`@vonArtmann/setup-cli/package.json`)
3. **Nix Derivation** (`default.nix` and `flake.nix`)
4. **Docker Image** (Multi-stage build with security scanning)

### Distribution Channels
1. **Primary**: Official website (https://oss.artmann.foundation/install?repo=template-health)
2. **GitHub Releases**: Latest versions with release notes
3. **Package Managers**: Platform-specific installers
4. **Container Registries**: Docker Hub and GitHub Container Registry

## Expected Deliverables
1. **Automated Build Pipeline**
   - GitHub Actions workflow for multi-platform builds
   - Release automation triggered by git tags
   - Security scanning and code signing integration

2. **Installation Methods**
   - Universal installation script with smart detection
   - Package manager configurations for major platforms
   - Container images with security best practices

3. **Documentation**
   - Installation guide for all supported methods
   - Troubleshooting guide for common issues
   - Integration examples for CI/CD platforms

4. **Testing Infrastructure**
   - Installation testing across platforms
   - Integration testing for package managers
   - Performance benchmarks for different installation methods

## Success Criteria
- **Installation Success Rate**: >95% across all platforms
- **Installation Time**: <30 seconds for all methods
- **Cross-Platform Compatibility**: >98% success rate
- **User Experience**: From install to first generation <2 minutes
- **Adoption Metrics**: >1,000 downloads/month within 3 months

## Security Considerations
1. **Supply Chain Security**
   - Reproducible builds with verification
   - Signed releases (GPG/cosign)
   - SBOM (Software Bill of Materials)
   - Vulnerability scanning in CI/CD

2. **Runtime Security**
   - Minimal permissions required
   - Secure defaults and input validation
   - Safe template rendering and execution

## Platform Support Matrix
| Platform | Method | Priority | Implementation |
|----------|--------|----------|----------------|
| macOS | Homebrew, Universal Script | High | Week 1 |
| Linux | APT, Snap, Universal Script | High | Week 1 |
| Windows | Chocolatey, Universal Script | High | Week 2 |
| Nix | Derivation, Flakes | Medium | Week 3 |
| Docker | Multi-arch Images | Medium | Week 2 |
| Cloud | Helm, Terraform | Low | Week 4 |

## User Experience Flow
1. **Discovery**: User finds installation instructions
2. **Installation**: Single command execution
3. **Verification**: Automatic installation testing
4. **First Use**: Immediate project generation capability
5. **Updates**: Seamless version management

## Related Files
- CLI implementation in `/cmd/generator/`
- Build configuration in `/.github/workflows/`
- Package configurations in `/packaging/`
- Documentation in `/docs/installation/`