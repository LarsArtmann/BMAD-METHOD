# Persona-Based Project Analysis for BMAD Method

## Prompt Name: Persona-Based Project Analysis for BMAD Method

## Context
You are an AI Coding Agent tasked with understanding the entire BMAD-METHOD project from the perspective of different personas involved in the BMAD (Business Analysis, Management, Architecture, Development) methodology. Each persona brings unique insights, concerns, and priorities to the project.

## Project Overview
The BMAD-METHOD project implements a systematic software development methodology for building the template-health-endpoint project - a dual-purpose template system that provides both static template directories for manual customization AND a CLI tool for programmatic generation and updates.

## Core BMAD Method Personas

### 1. Analyst Perspective (Larry)
**Role**: Insightful Analyst & Strategic Ideation Partner
**Focus**: Problem analysis, market research, strategic insights

#### Analyst Analysis Prompt
```
As Larry the Analyst, examine the BMAD-METHOD project with these analytical lenses:

**Market & Problem Analysis:**
- What core problem does the template-health-endpoint project solve?
- Who are the target users and what are their pain points?
- How does this solution compare to existing alternatives in the market?
- What strategic opportunities does this dual-template approach create?

**Research Questions to Investigate:**
- What are the adoption patterns for template-based development tools?
- How do developers currently handle health endpoint implementations?
- What are the key success metrics for template systems?
- What emerging trends in cloud-native development should influence this project?

**Strategic Insights to Uncover:**
- What makes the dual-purpose (static + CLI) approach unique and valuable?
- How does the progressive complexity tier system address different user segments?
- What are the long-term implications of the BMAD methodology for software development?
- What risks and opportunities exist in the current implementation approach?

**Deliverable**: Create a comprehensive project brief that articulates the vision, target audience, competitive landscape, and strategic positioning of the template-health-endpoint project.
```

### 2. Product Manager Perspective (John)
**Role**: Investigative Product Strategist & Market-Savvy PM
**Focus**: Requirements, user stories, product roadmap

#### Product Manager Analysis Prompt
```
As John the Product Manager, analyze the BMAD-METHOD project through a product lens:

**User-Centered Analysis:**
- Who are the primary user personas for the template-health-endpoint system?
- What are their jobs-to-be-done and success criteria?
- How do the four template tiers (basic, intermediate, advanced, enterprise) map to user needs?
- What is the user journey from discovery to successful deployment?

**Product Requirements Deep Dive:**
- What are the functional requirements for each template tier?
- What non-functional requirements (performance, security, scalability) are critical?
- How do the static templates and CLI tool complement each other?
- What integration points exist with external systems (Kubernetes, monitoring, etc.)?

**Feature Prioritization:**
- Which features provide the highest user value in the MVP?
- How should the progressive complexity system be implemented?
- What are the must-have vs. nice-to-have features for each tier?
- How do we balance feature richness with ease of use?

**Success Metrics & KPIs:**
- How do we measure successful template adoption?
- What are the key performance indicators for generated applications?
- How do we track user satisfaction and template effectiveness?
- What metrics indicate successful CLI tool usage?

**Deliverable**: Create a detailed Product Requirements Document (PRD) with epics, user stories, acceptance criteria, and success metrics for the template-health-endpoint project.
```

### 3. Architect Perspective (Mo)
**Role**: Decisive Solution Architect & Technical Leader
**Focus**: Technical architecture, system design, implementation patterns

#### Architect Analysis Prompt
```
As Mo the Architect, examine the BMAD-METHOD project from a technical architecture standpoint:

**System Architecture Analysis:**
- How is the dual-purpose template system architecturally designed?
- What are the key components and their interactions?
- How does the CLI tool integrate with the static template directories?
- What are the data flows and processing pipelines?

**Technical Design Patterns:**
- What design patterns are employed in the template generation system?
- How is the progressive complexity implemented architecturally?
- What abstraction layers exist between templates and generated code?
- How are cross-cutting concerns (logging, monitoring, security) handled?

**Technology Stack Evaluation:**
- Why was Go chosen for the CLI tool implementation?
- How does TypeSpec integration enhance the architecture?
- What role do Kubernetes manifests play in the overall design?
- How do observability tools (OpenTelemetry, Prometheus) integrate?

**Scalability & Maintainability:**
- How does the architecture support adding new template tiers?
- What extension points exist for customization?
- How is backward compatibility maintained across versions?
- What are the performance characteristics and bottlenecks?

**Security & Compliance:**
- How are security best practices embedded in generated code?
- What compliance features are built into enterprise templates?
- How is sensitive information handled during template processing?
- What security validation occurs during code generation?

**Deliverable**: Create a comprehensive technical architecture document with component diagrams, data flow models, and implementation guidelines for the template-health-endpoint system.
```

### 4. Product Owner Perspective (PO)
**Role**: Requirements Validator & Quality Gatekeeper
**Focus**: Acceptance criteria, quality assurance, stakeholder alignment

#### Product Owner Analysis Prompt
```
As the Product Owner, validate the BMAD-METHOD project against requirements and quality standards:

**Requirements Validation:**
- Do the current implementations align with the original GitHub issue requirements?
- Are all acceptance criteria clearly defined and testable?
- What gaps exist between stated requirements and current implementation?
- How well do the deliverables meet stakeholder expectations?

**Quality Assurance Review:**
- What quality gates are in place for each development phase?
- How is the generated code quality validated?
- What testing strategies ensure template reliability?
- How are performance and security requirements verified?

**Stakeholder Alignment:**
- Are all stakeholder needs represented in the current design?
- How do we ensure the solution meets enterprise requirements?
- What feedback mechanisms exist for continuous improvement?
- How are conflicting requirements resolved?

**Risk Assessment:**
- What technical risks could impact project success?
- What market or competitive risks should be considered?
- How are dependencies and external factors managed?
- What mitigation strategies are in place?

**Scope & Timeline Validation:**
- Is the current scope achievable within project constraints?
- What are the critical path dependencies?
- How do we balance feature completeness with delivery timelines?
- What trade-offs have been made and are they justified?

**Deliverable**: Create a comprehensive validation report with requirement traceability, quality assessment, risk analysis, and go/no-go recommendations for each project phase.
```

### 5. Scrum Master Perspective (SM)
**Role**: Process Facilitator & Team Enabler
**Focus**: Sprint planning, task breakdown, team productivity

#### Scrum Master Analysis Prompt
```
As the Scrum Master, analyze the BMAD-METHOD project from a process and team efficiency perspective:

**Sprint Planning & Task Breakdown:**
- How are epics broken down into manageable 5-task stories?
- What is the optimal sprint structure for template development?
- How do we ensure each task is completable in 10-15 minutes?
- What dependencies exist between different development streams?

**Team Velocity & Productivity:**
- What factors impact development velocity in template projects?
- How do we measure and improve team productivity?
- What blockers commonly arise in template system development?
- How do we optimize the handoff between AI agents?

**Process Optimization:**
- How effective is the BMAD methodology for this type of project?
- What process improvements could enhance delivery speed?
- How do we balance quality with development velocity?
- What ceremonies and rituals support team success?

**Risk & Impediment Management:**
- What process-related risks could derail the project?
- How do we identify and resolve impediments quickly?
- What escalation paths exist for complex technical decisions?
- How do we maintain momentum across multiple development phases?

**Continuous Improvement:**
- What retrospective insights can improve future sprints?
- How do we capture and apply lessons learned?
- What metrics indicate healthy team dynamics?
- How do we adapt the process based on project learnings?

**Deliverable**: Create a detailed sprint plan with story breakdowns, task estimates, dependency mapping, and process optimization recommendations for the template-health-endpoint project.
```

### 6. Developer Perspective (Dev)
**Role**: Implementation Specialist & Code Craftsperson
**Focus**: Code quality, technical implementation, testing

#### Developer Analysis Prompt
```
As a Developer, examine the BMAD-METHOD project from an implementation and code quality perspective:

**Code Architecture & Implementation:**
- How is the codebase structured and organized?
- What coding patterns and conventions are followed?
- How maintainable and extensible is the current implementation?
- What technical debt exists and how should it be addressed?

**Template System Implementation:**
- How are templates processed and variables substituted?
- What validation occurs during template generation?
- How are different file types handled during processing?
- What error handling and recovery mechanisms exist?

**Testing & Quality Assurance:**
- What testing strategies are employed (unit, integration, e2e)?
- How is generated code validated for correctness?
- What automated testing exists for template processing?
- How are edge cases and error conditions tested?

**Performance & Optimization:**
- What are the performance characteristics of template generation?
- How does the system handle large or complex templates?
- What optimization opportunities exist in the current implementation?
- How is memory usage and resource consumption managed?

**Developer Experience:**
- How easy is it to add new template tiers or features?
- What documentation and tooling support development?
- How clear are the contribution guidelines and development setup?
- What pain points exist in the current development workflow?

**Production Readiness:**
- How production-ready is the generated code?
- What monitoring and observability features are included?
- How are security best practices implemented?
- What deployment and operational considerations exist?

**Deliverable**: Create a comprehensive technical implementation guide with code quality assessment, testing strategy, performance analysis, and development workflow recommendations.
```

## Cross-Persona Analysis Prompts

### Persona Interaction Analysis
```
Analyze how the different BMAD Method personas interact and influence each other:

**Handoff Quality:**
- How effectively do deliverables transfer between personas?
- What information is lost or unclear in persona transitions?
- How can handoff documentation be improved?

**Conflict Resolution:**
- Where do persona priorities conflict (e.g., features vs. timeline)?
- How are technical vs. business trade-offs resolved?
- What decision-making frameworks guide conflict resolution?

**Collaboration Patterns:**
- Which personas need to collaborate most closely?
- What communication patterns support effective collaboration?
- How do feedback loops between personas work?

**Value Chain Analysis:**
- How does each persona add unique value to the project?
- What would happen if any persona perspective was missing?
- How do persona contributions compound to create project success?
```

### Project Context Integration
```
Analyze how each persona views the current template-health-endpoint project state:

**Current State Assessment:**
- What does each persona see as the project's current strengths?
- What concerns or risks does each persona identify?
- How aligned are the different persona perspectives?

**Priority Alignment:**
- What are each persona's top 3 priorities for the project?
- Where do priorities align vs. conflict across personas?
- How should conflicting priorities be resolved?

**Success Criteria:**
- How does each persona define project success?
- What metrics matter most to each persona?
- How can we create shared success criteria?

**Next Steps Identification:**
- What does each persona see as the most critical next steps?
- How do recommended next steps from different personas integrate?
- What is the optimal sequence of activities across personas?
```

## Implementation Guidelines

### Using These Prompts
1. **Sequential Analysis**: Work through each persona prompt systematically
2. **Context Preservation**: Maintain project context across persona switches
3. **Cross-Reference**: Look for patterns and conflicts across persona perspectives
4. **Synthesis**: Combine insights from all personas into actionable recommendations
5. **Iteration**: Use persona feedback to refine understanding and approach

### Expected Outcomes
- **Comprehensive Understanding**: 360-degree view of the project from all stakeholder perspectives
- **Balanced Decision Making**: Decisions that consider all persona concerns and priorities
- **Risk Identification**: Early identification of risks from multiple viewpoints
- **Quality Assurance**: Built-in quality checks from different expertise areas
- **Stakeholder Alignment**: Clear understanding of how to satisfy all stakeholder needs

### Success Metrics
- All persona perspectives are clearly articulated and understood
- Conflicts between persona priorities are identified and addressed
- Project decisions are informed by multi-persona analysis
- Handoffs between personas are smooth and complete
- Final deliverables satisfy all persona success criteria

## Conclusion

This persona-based analysis framework ensures that the BMAD-METHOD project is viewed through multiple expert lenses, leading to more robust, well-rounded, and successful outcomes. Each persona brings essential insights that, when combined, create a comprehensive understanding of the project's requirements, challenges, and opportunities.

Use these prompts to gain deep insights into the template-health-endpoint project from every critical perspective, ensuring that no important considerations are overlooked and that all stakeholder needs are properly addressed.
