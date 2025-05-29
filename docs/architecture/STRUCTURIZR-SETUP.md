# Structurizr Setup Guide - BMAD-METHOD C4 Model

## Quick Start Options

### Option 1: Structurizr Lite (Recommended for Local Development)

**Step 1: Install Docker**
```bash
# Ensure Docker is installed and running
docker --version
```

**Step 2: Run Structurizr Lite**
```bash
# Navigate to project root
cd /path/to/BMAD-METHOD

# Run Structurizr Lite with our model
docker run -it --rm -p 8080:8080 \
  -v $(pwd)/docs/architecture:/usr/local/structurizr \
  structurizr/lite
```

**Step 3: Access the Model**
1. Open browser to http://localhost:8080
2. The C4 model will load automatically
3. Navigate through different views using the sidebar

### Option 2: Structurizr Cloud (Best for Sharing)

**Step 1: Create Account**
1. Visit https://structurizr.com/
2. Sign up for free account
3. Create a new workspace

**Step 2: Upload Model**
1. Copy content from `docs/architecture/c4-model.dsl`
2. Paste into Structurizr workspace
3. Save and view diagrams

**Step 3: Share**
- Get shareable link from workspace
- Invite team members
- Export diagrams as needed

### Option 3: VS Code Extension (For Editing)

**Step 1: Install Extension**
1. Open VS Code
2. Install "Structurizr DSL" extension
3. Restart VS Code

**Step 2: Open Model**
1. Open `docs/architecture/c4-model.dsl`
2. Use Ctrl+Shift+P → "Structurizr: Preview"
3. Edit and preview in real-time

## Available Views

### System Context Views
- **SystemContext**: High-level system overview
- Shows users, external systems, and main relationships

### Container Views
- **TemplateSystemContainers**: Internal structure of template system
- **GeneratedAppContainers**: Structure of generated applications

### Component Views
- **CLIComponents**: Command-line interface structure
- **TemplateEngineComponents**: Core generation engine
- **TemplateRepoComponents**: Template organization
- **ConfigSystemComponents**: Configuration management
- **TestingFrameworkComponents**: Testing and validation
- **BasicAppComponents**: Basic tier application
- **IntermediateAppComponents**: Intermediate tier application
- **AdvancedAppComponents**: Advanced tier application
- **EnterpriseAppComponents**: Enterprise tier application
- **K8sDeploymentComponents**: Kubernetes deployment
- **TypeScriptClientComponents**: TypeScript SDK

### Dynamic Views
- **TemplateGeneration**: Template generation process flow
- **EnterpriseSecurityFlow**: Security request handling
- **ObservabilityFlow**: Observability and monitoring flow

### Deployment Views
- **Production**: Complete production deployment architecture

## Navigation Tips

### 1. Start with System Context
- Understand overall system purpose
- Identify key users and external systems
- Get high-level relationship overview

### 2. Drill Down to Containers
- See technology choices
- Understand major components
- Identify integration points

### 3. Explore Components
- Understand internal structure
- See detailed responsibilities
- Analyze component interactions

### 4. Review Dynamic Views
- Understand process flows
- See sequence of interactions
- Identify key decision points

### 5. Study Deployment
- Understand infrastructure requirements
- See runtime environment
- Plan deployment strategy

## Customization

### Adding New Views
```dsl
# Add to views section
component newContainer "NewComponentView" {
    include *
    autoLayout
    title "New Component - Component View"
    description "Description of new component structure"
}
```

### Modifying Styles
```dsl
# Add to styles section
element "NewElementType" {
    background #ff0000
    color #ffffff
    shape Box
}
```

### Adding Relationships
```dsl
# Add to model section
componentA -> componentB "Description of relationship"
```

## Export Options

### From Structurizr Cloud
1. **PNG/SVG**: High-quality images
2. **PDF**: Documentation-ready format
3. **PlantUML**: For other tools
4. **Mermaid**: For GitHub/GitLab

### From Structurizr Lite
1. Use browser developer tools
2. Right-click diagrams to save
3. Print to PDF for documentation

### From VS Code
1. Use extension export features
2. Generate images directly
3. Integrate with documentation

## Integration with Documentation

### Embedding in README
```markdown
![System Context](docs/architecture/exports/system-context.png)
```

### Wiki Integration
- Export diagrams as images
- Link to live Structurizr workspace
- Include in architectural documentation

### Presentation Use
- Export high-resolution images
- Use in architecture presentations
- Include in design reviews

## Maintenance

### Keeping Model Updated
1. **Code Changes**: Update components when code structure changes
2. **New Features**: Add new components and relationships
3. **Deployment Changes**: Update deployment views
4. **External Systems**: Update when integrations change

### Version Control
- Model is stored in `docs/architecture/c4-model.dsl`
- Version controlled with code
- Changes tracked in git history
- Review changes in pull requests

### Team Collaboration
1. **Shared Workspace**: Use Structurizr Cloud for team access
2. **Review Process**: Include model updates in code reviews
3. **Documentation**: Keep documentation synchronized
4. **Training**: Ensure team understands C4 methodology

## Troubleshooting

### Common Issues

**Docker Permission Issues**
```bash
# Fix Docker permissions on Linux
sudo usermod -aG docker $USER
# Logout and login again
```

**Port Already in Use**
```bash
# Use different port
docker run -it --rm -p 8081:8080 \
  -v $(pwd)/docs/architecture:/usr/local/structurizr \
  structurizr/lite
```

**Model Syntax Errors**
- Check DSL syntax in VS Code extension
- Validate relationships exist
- Ensure proper nesting structure

**Performance Issues**
- Large models may load slowly
- Consider splitting into multiple workspaces
- Use filtered views for complex models

### Getting Help
1. **Structurizr Documentation**: https://structurizr.com/help
2. **C4 Model Guide**: https://c4model.com/
3. **Community**: Structurizr Slack community
4. **Examples**: Structurizr example repository

## Best Practices

### Model Organization
1. **Logical Grouping**: Group related components
2. **Consistent Naming**: Use clear, consistent names
3. **Appropriate Detail**: Right level of detail for audience
4. **Regular Updates**: Keep model current with code

### Documentation
1. **Clear Descriptions**: Meaningful component descriptions
2. **Relationship Labels**: Clear relationship descriptions
3. **View Purposes**: Explain what each view shows
4. **Context**: Provide architectural context

### Team Usage
1. **Shared Understanding**: Ensure team understands model
2. **Review Process**: Include in architectural reviews
3. **Decision Records**: Link to architectural decisions
4. **Training**: Provide C4 methodology training

This setup guide ensures you can effectively use and maintain the comprehensive C4 model for the BMAD-METHOD template system.
