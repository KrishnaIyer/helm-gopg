#!/bin/bash

# This script is used to bump the version of the project.
# The version is stored in the VERSION file.
# The version is also updated in the plugin.yaml file.
# The user can choose to bump either the minor or the patch version.
# The script will then update the VERSION file and the plugin.yaml file.

current_version=$(cat VERSION)
echo "Current version: $current_version"

type=$1

if [ "$type" == "minor" ]; then
    new_version=$(echo $current_version | awk -F. '{$(NF-1) = $(NF-1) + 1;} 1' | sed 's/ /./g' )
elif [ "$type" == "patch" ]; then
    new_version=$(echo $current_version | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')
else
    echo "Invalid version type. Please use 'minor' or 'patch'."
    exit 1
fi

echo "New version: $new_version"

echo "Updating VERSION file..."
echo $new_version > VERSION

echo "Updating plugin.yaml file..."
sed -i '' "s/$current_version/$new_version/g" plugin.yaml

echo "Commit changes and tag"
git add VERSION plugin.yaml
git commit -m "Bump version to v$new_version"
git tag v$new_version
