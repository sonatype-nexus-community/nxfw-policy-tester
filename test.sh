#!/bin/bash

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Read credentials from environment variables
NEXUS_USERNAME="${NEXUS_USERNAME}"
NEXUS_PASSWORD="${NEXUS_PASSWORD}"

# Check if credentials are set
if [ -z "$NEXUS_USERNAME" ] || [ -z "$NEXUS_PASSWORD" ]; then
    echo -e "${RED}Error: NEXUS_USERNAME and NEXUS_PASSWORD environment variables must be set.${NC}"
    echo "Example: export NEXUS_USERNAME='your_username'"
    echo "         export NEXUS_PASSWORD='your_password'"
    exit 1
fi

# Prompt for package format
echo -e "${YELLOW}Select package format:${NC}"
echo "1) NPM"
echo "2) Maven"
read -p "Enter choice (1 or 2): " -r FORMAT_CHOICE
echo ""

case $FORMAT_CHOICE in
    1)
        PACKAGE_FORMAT="npm"
        ;;
    2)
        PACKAGE_FORMAT="maven"
        ;;
    *)
        echo -e "${RED}Error: Invalid choice. Please enter 1 or 2.${NC}"
        exit 1
        ;;
esac

# Prompt for Nexus Repository URL
echo -e "${YELLOW}Enter your Sonatype Nexus Repository URL:${NC}"
echo "(Example: https://nexus.example.com)"
read -r NEXUS_URL

# Remove trailing slash if present
NEXUS_URL="${NEXUS_URL%/}"

# Validate URL format
if [[ ! "$NEXUS_URL" =~ ^https?:// ]]; then
    echo -e "${RED}Error: Invalid URL format. URL must start with http:// or https://${NC}"
    exit 1
fi

# Prompt for repository name based on format
echo ""
if [ "$PACKAGE_FORMAT" = "npm" ]; then
    echo -e "${YELLOW}Enter the NPM proxy repository name:${NC}"
    echo "(Example: npm-proxy)"
else
    echo -e "${YELLOW}Enter the Maven proxy repository name:${NC}"
    echo "(Example: maven-proxy)"
fi
read -r REPO_NAME

# Validate repository name is not empty
if [ -z "$REPO_NAME" ]; then
    echo -e "${RED}Error: Repository name cannot be empty.${NC}"
    exit 1
fi

# Define packages based on selected format
if [ "$PACKAGE_FORMAT" = "npm" ]; then
    # NPM packages
    # Format: "name:version:label"
    declare -a PACKAGES=(
        "bson:1.0.9:Security-Critical"
        "braces:1.8.5:Security-High"
        "cookie:0.3.1:Security-Medium"
    )
else
    # Maven packages
    # Format: "group/artifact:version:label:extension"
    declare -a PACKAGES=(
        "com.amazonaws/aws-android-sdk-core:2.75.0:Security-Critical:aar"
        "org.jsoup/jsoup:1.13.1:Security-High:jar"
        "ant/ant:1.6.5:Security-Medium:jar"
    )
fi

# Arrays to track results (parallel to PACKAGES array)
declare -a PACKAGE_RESULTS=()

# Function to get color for label
get_label_color() {
    case "$1" in
        "Security-Critical")
            echo "$RED"
            ;;
        "Security-High")
            echo "$MAGENTA"
            ;;
        "Security-Medium")
            echo "$YELLOW"
            ;;
        *)
            echo "$NC"
            ;;
    esac
}

# Function to construct Maven URL
# Maven path format: group/artifact/version/artifact-version.extension
construct_maven_url() {
    local group_artifact="$1"
    local version="$2"
    local extension="$3"
    
    # Split group/artifact
    IFS='/' read -r group artifact <<< "$group_artifact"
    
    # Convert group dots to slashes (e.g., com.amazonaws -> com/amazonaws)
    local group_path="${group//./\/}"
    
    # Construct the Maven path
    echo "$NEXUS_URL/repository/$REPO_NAME/$group_path/$artifact/$version/$artifact-$version.$extension"
}

# Function to construct NPM URL
construct_npm_url() {
    local name="$1"
    local version="$2"
    
    echo "$NEXUS_URL/repository/$REPO_NAME/$name/-/$name-$version.tgz"
}

# Display summary
echo ""
echo -e "${YELLOW}=== Configuration Summary ===${NC}"
echo "Nexus URL: $NEXUS_URL"
echo "Username: $NEXUS_USERNAME"
if [ "$PACKAGE_FORMAT" = "npm" ]; then
    echo "Format: ${CYAN}NPM${NC}"
    echo "Repository: $REPO_NAME"
else
    echo "Format: ${BLUE}Maven${NC}"
    echo "Repository: $REPO_NAME"
fi
echo ""
echo "Packages to check:"

if [ "$PACKAGE_FORMAT" = "npm" ]; then
    for pkg in "${PACKAGES[@]}"; do
        IFS=':' read -r name version label <<< "$pkg"
        label_color=$(get_label_color "$label")
        echo -e "  - $name@$version ${label_color}[$label]${NC}"
    done
else
    for pkg in "${PACKAGES[@]}"; do
        IFS=':' read -r name version label extension <<< "$pkg"
        label_color=$(get_label_color "$label")
        echo -e "  - $name@$version (.$extension) ${label_color}[$label]${NC}"
    done
fi
echo ""

# Prompt for confirmation
read -p "Proceed with checking packages? (y/n): " -n 1 -r
echo ""

if [[ ! $REPONSE =~ ^[Yy]$ ]] && [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Check cancelled."
    exit 0
fi

echo ""
echo -e "${YELLOW}=== Checking Package Availability ===${NC}"
echo ""

# Check each package
SUCCESS_COUNT=0
FAIL_COUNT=0

if [ "$PACKAGE_FORMAT" = "npm" ]; then
    # Check NPM packages
    for pkg in "${PACKAGES[@]}"; do
        IFS=':' read -r name version label <<< "$pkg"
        label_color=$(get_label_color "$label")
        
        PACKAGE_URL=$(construct_npm_url "$name" "$version")
        display_name="$name@$version"
        
        echo -e "Checking $display_name ${label_color}[$label]${NC}..."
        
        # Perform curl request with authentication (HEAD request to not download)
        HTTP_CODE=$(curl -u "$NEXUS_USERNAME:$NEXUS_PASSWORD" \
                         -w "%{http_code}" \
                         -o /dev/null \
                         -s \
                         --head \
                         "$PACKAGE_URL")
        
        if [ "$HTTP_CODE" -eq 200 ]; then
            echo -e "${GREEN}✓ Package available: $display_name ${label_color}[$label]${NC}"
            PACKAGE_RESULTS+=("success")
            ((SUCCESS_COUNT++))
        else
            echo -e "${RED}✗ Package not available: $display_name ${label_color}[$label]${NC} (HTTP $HTTP_CODE)"
            PACKAGE_RESULTS+=("failed")
            ((FAIL_COUNT++))
        fi
        echo ""
    done
else
    # Check Maven packages
    for pkg in "${PACKAGES[@]}"; do
        IFS=':' read -r name version label extension <<< "$pkg"
        label_color=$(get_label_color "$label")
        
        PACKAGE_URL=$(construct_maven_url "$name" "$version" "$extension")
        display_name="$name@$version (.$extension)"
        
        echo -e "Checking $display_name ${label_color}[$label]${NC}..."
        
        # Perform curl request with authentication (HEAD request to not download)
        HTTP_CODE=$(curl -u "$NEXUS_USERNAME:$NEXUS_PASSWORD" \
                         -w "%{http_code}" \
                         -o /dev/null \
                         -s \
                         --head \
                         "$PACKAGE_URL")
        
        if [ "$HTTP_CODE" -eq 200 ]; then
            echo -e "${GREEN}✓ Package available: $display_name ${label_color}[$label]${NC}"
            PACKAGE_RESULTS+=("success")
            ((SUCCESS_COUNT++))
        else
            echo -e "${RED}✗ Package not available: $display_name ${label_color}[$label]${NC} (HTTP $HTTP_CODE)"
            PACKAGE_RESULTS+=("failed")
            ((FAIL_COUNT++))
        fi
        echo ""
    done
fi

# Summary
echo -e "${YELLOW}=== Check Summary ===${NC}"
echo -e "${GREEN}Available: $SUCCESS_COUNT${NC}"
echo -e "${RED}Not Available: $FAIL_COUNT${NC}"

# Display packages by security level
echo ""
echo -e "${YELLOW}=== Security Level Breakdown ===${NC}"
index=0

if [ "$PACKAGE_FORMAT" = "npm" ]; then
    for pkg in "${PACKAGES[@]}"; do
        IFS=':' read -r name version label <<< "$pkg"
        label_color=$(get_label_color "$label")
        display_name="$name@$version"
        
        if [ "${PACKAGE_RESULTS[$index]}" = "success" ]; then
            status="${GREEN}[Available]${NC}"
        else
            status="${RED}[Not Available]${NC}"
        fi
        
        echo -e "${label_color}$label${NC}: $display_name $status"
        ((index++))
    done
else
    for pkg in "${PACKAGES[@]}"; do
        IFS=':' read -r name version label extension <<< "$pkg"
        label_color=$(get_label_color "$label")
        display_name="$name@$version (.$extension)"
        
        if [ "${PACKAGE_RESULTS[$index]}" = "success" ]; then
            status="${GREEN}[Available]${NC}"
        else
            status="${RED}[Not Available]${NC}"
        fi
        
        echo -e "${label_color}$label${NC}: $display_name $status"
        ((index++))
    done
fi

exit 0
