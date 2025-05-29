#!/bin/bash

# Convert generated projects to template format
# This script converts hardcoded values to template variables

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Function to convert a file to template format
convert_file_to_template() {
    local file=$1
    local temp_file="${file}.tmp"
    
    # Skip binary files and directories
    if [[ -d "$file" ]] || ! [[ -f "$file" ]]; then
        return 0
    fi
    
    # Skip if file doesn't exist
    if [[ ! -f "$file" ]]; then
        return 0
    fi
    
    log_info "Converting $file to template format..."
    
    # Create template version with variables
    sed \
        -e 's/template-basic/{{.Config.Name}}/g' \
        -e 's/github\.com\/template\/basic/{{.Config.GoModule}}/g' \
        -e 's/Basic health endpoint service/{{.Config.Description}}/g' \
        -e 's/Template Health Endpoint Generator v1\.0\.0/Template Health Endpoint Generator v{{.Version}}/g' \
        -e 's/Generated at: [0-9T:-]*Z/Generated at: {{.Timestamp}}/g' \
        "$file" > "$temp_file"
    
    # Replace original with template version
    mv "$temp_file" "$file"
}

# Function to add template file extension where needed
add_template_extension() {
    local file=$1
    local base_name=$(basename "$file")
    
    # Add .tmpl extension to files that need templating
    case "$base_name" in
        "go.mod"|"README.md"|"package.json"|"Dockerfile"|"docker-compose.yml")
            if [[ ! "$file" =~ \.tmpl$ ]]; then
                mv "$file" "${file}.tmpl"
                log_success "Renamed $file to ${file}.tmpl"
            fi
            ;;
    esac
}

# Function to process a template directory
process_template_dir() {
    local template_dir=$1
    local tier_name=$2
    
    log_info "Processing $tier_name template directory: $template_dir"
    
    if [[ ! -d "$template_dir" ]]; then
        log_info "Directory $template_dir does not exist, skipping..."
        return 0
    fi
    
    # Find all files and convert them
    find "$template_dir" -type f | while read -r file; do
        convert_file_to_template "$file"
        add_template_extension "$file"
    done
    
    # Create template metadata
    cat > "$template_dir/template.yaml" << EOF
name: $tier_name
description: ${tier_name^} tier health endpoint template
tier: $tier_name
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: $([ "$tier_name" != "basic" ] && echo "true" || echo "false")
  cloudevents: $([ "$tier_name" = "advanced" ] || [ "$tier_name" = "enterprise" ] && echo "true" || echo "false")
version: "1.0.0"
EOF
    
    log_success "Created template metadata for $tier_name tier"
}

# Main execution
main() {
    log_info "Converting generated projects to template format..."
    
    # Process each tier
    for tier in basic intermediate advanced enterprise; do
        if [[ -d "templates/$tier" ]]; then
            process_template_dir "templates/$tier" "$tier"
        else
            log_info "Template directory templates/$tier does not exist, skipping..."
        fi
    done
    
    log_success "Template conversion completed!"
}

# Run main function
main "$@"
