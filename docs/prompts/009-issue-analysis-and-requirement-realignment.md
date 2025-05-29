# Issue Analysis and Requirement Realignment

## Prompt Name
**Issue Analysis and Requirement Realignment**

## Context
When working on a project, it's critical to periodically step back and verify that the current implementation aligns with the original requirements. This prompt helps identify misalignments and course-correct before significant effort is wasted.

## Objective
Analyze the original issue/requirements against current implementation to identify gaps, misalignments, or scope drift, then realign the project to meet the actual requirements.

## Key Steps

### 1. Original Requirement Analysis
- **Read the original issue/ticket completely** including all comments and related issues
- **Extract the core requirements** and expected deliverables
- **Identify the target users** and their specific needs
- **Document the success criteria** as defined in the original issue

### 2. Current State Assessment
- **Audit the current implementation** against original requirements
- **Identify what has been built** vs. what was requested
- **Document any scope drift** or architectural decisions that deviated from requirements
- **Assess the percentage completion** relative to original goals

### 3. Gap Analysis
- **List missing features** that were in the original requirements
- **Identify over-engineered components** that weren't requested
- **Document architectural mismatches** (e.g., CLI tool vs. template repository)
- **Prioritize gaps** by impact and effort required

### 4. Realignment Strategy
- **Determine if current work can be adapted** to meet original requirements
- **Plan the minimal changes** needed to align with requirements
- **Identify what needs to be restructured** vs. what can be kept
- **Create a clear path forward** that satisfies original requirements

## Example Application
In the template-health-endpoint project:

**Original Issue #127 Required:**
- Template repository following template-* pattern
- Static template directories users can copy/fork
- 4 tiers: Basic → Intermediate → Advanced → Enterprise
- TypeSpec-first API definitions

**Initial Implementation:**
- CLI tool with embedded templates
- No static template directories
- Only basic tier implemented
- CLI-focused approach

**Realignment:**
- Restructure to dual-purpose system (static templates + CLI)
- Create actual template directories in `/templates/`
- Implement all 4 tiers as requested
- Maintain CLI functionality while adding static templates

## Success Criteria
- [ ] Original requirements are fully understood and documented
- [ ] Current implementation gaps are clearly identified
- [ ] Realignment strategy addresses all original requirements
- [ ] Path forward is clear and achievable
- [ ] Stakeholder expectations are reset if needed

## Best Practices
- **Read everything**: Don't assume you understand requirements without reading all details
- **Use available tools**: GitHub API, CLI tools to fetch original issues
- **Document gaps clearly**: Be specific about what's missing vs. what's extra
- **Preserve valuable work**: Find ways to adapt current work rather than starting over
- **Communicate changes**: Update stakeholders on any scope or timeline adjustments

## Related Patterns
- Requirements analysis and validation
- Project scope management
- Technical debt assessment
- Stakeholder communication
- Agile course correction
