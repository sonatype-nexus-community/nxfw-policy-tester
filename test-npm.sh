#!/bin/bash

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
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

# Prompt for repository name
echo ""
echo -e "${YELLOW}Enter the NPM proxy repository name:${NC}"
echo "(Example: npm-proxy)"
read -r REPO_NAME

# Validate repository name is not empty
if [ -z "$REPO_NAME" ]; then
    echo -e "${RED}Error: Repository name cannot be empty.${NC}"
    exit 1
fi

# Define NPM packages to download with labels
# Format: "name:version:label"
declare -a PACKAGES=(
    "bson:1.0.9:Security-Critical"
    "braces:1.8.5:Security-High"
    "cookie:0.3.1:Security-Medium"
)

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

# Display summary
echo ""
echo -e "${YELLOW}=== Configuration Summary ===${NC}"
echo "Nexus URL: $NEXUS_URL"
echo "Username: $NEXUS_USERNAME"
echo "Repository: $REPO_NAME"
echo ""
echo "Packages to check:"
for pkg in "${PACKAGES[@]}"; do
    IFS=':' read -r name version label <<< "$pkg"
    label_color=$(get_label_color "$label")
    echo -e "  - $name@$version ${label_color}[$label]${NC}"
done
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

for pkg in "${PACKAGES[@]}"; do
    IFS=':' read -r name version label <<< "$pkg"
    label_color=$(get_label_color "$label")
    
    # Construct the Nexus URL for the package
    # Format: /repository/[repo-name]/[package]/-/[package]-[version].tgz
    PACKAGE_URL="$NEXUS_URL/repository/$REPO_NAME/$name/-/$name-$version.tgz"
    
    echo -e "Checking $name@$version ${label_color}[$label]${NC}..."
    
    # Perform curl request with authentication (HEAD request to not download)
    HTTP_CODE=$(curl -u "$NEXUS_USERNAME:$NEXUS_PASSWORD" \
                     -w "%{http_code}" \
                     -o /dev/null \
                     -s \
                     --head \
                     "$PACKAGE_URL")
    
    if [ "$HTTP_CODE" -eq 200 ]; then
        echo -e "${GREEN}✓ Package available: $name@$version ${label_color}[$label]${NC}"
        PACKAGE_RESULTS+=("success")
        ((SUCCESS_COUNT++))
    else
        echo -e "${RED}✗ Package not available: $name@$version ${label_color}[$label]${NC} (HTTP $HTTP_CODE)"
        PACKAGE_RESULTS+=("failed")
        ((FAIL_COUNT++))
    fi
    echo ""
done

# Summary
echo -e "${YELLOW}=== Check Summary ===${NC}"
echo -e "${GREEN}Available: $SUCCESS_COUNT${NC}"
echo -e "${RED}Not Available: $FAIL_COUNT${NC}"

# Display packages by security level
echo ""
echo -e "${YELLOW}=== Security Level Breakdown ===${NC}"
index=0
for pkg in "${PACKAGES[@]}"; do
    IFS=':' read -r name version label <<< "$pkg"
    label_color=$(get_label_color "$label")
    
    if [ "${PACKAGE_RESULTS[$index]}" = "success" ]; then
        status="${GREEN}[Available]${NC}"
    else
        status="${RED}[Not Available]${NC}"
    fi
    
    echo -e "${label_color}$label${NC}: $name@$version $status"
    ((index++))
done

exit 0
