#!/bin/bash

printf 'What is the new github url? (e.g. github.com/henrywhitaker3/connect-template) '
read repo

repoName=$(echo $repo | sed -e "s/github.com\///g")
baseName=$(echo $repoName | sed -e "s~.*/~~g")

printf 'What is the name of your service? '
read name


# Replace the repo url in the .releaserc file
sed -i "s~repositoryUrl\": \".*\",$~repositoryUrl\": \"https://$repo\",~g" .releaserc
# and in the docker build steps
sed -i "s~ghcr.io/henrywhitaker3/connect-template~ghcr.io/$repoName~g" .github/workflows/test.yaml
sed -i "s~ghcr.io/henrywhitaker3/connect-template~ghcr.io/$repoName~g" .github/workflows/release.yaml
sed -i "s~ghcr.io/henrywhitaker3/connect-template~ghcr.io/$repoName~g" chart/values.yaml

# Update the go.mod file
sed -i "s~module github.com/henrywhitaker3/connect-template~module $repo~g" go.mod
# and now do all the files that need stuff importing...
find . -name '*.go' -print0 | xargs -0 sed -i "s~github.com/henrywhitaker3/connect-template~$repo~g"
find . -name '*.proto' -print0 | xargs -0 sed -i "s~github.com/henrywhitaker3/connect-template~$repo~g"

# and the name in the example config file
sed -i "s~connect-template~$name~g" api.example.yaml
# and the name in the buf template
sed -i "s~connect-template~$name~g" hello/v1/hello.proto

# Now do the default config file location
sed -i "s~connect-template~$baseName~g" chart/values.yaml
sed -i "s~connect-template~$baseName~g" chart/Chart.yaml
sed -i "s~connect-template~$baseName~g" .github/workflows/release.yaml
sed -i "s~connect-template.yaml~$baseName.yaml~g" main.go
sed -i "s~connect-template.yaml~$baseName.yaml~g" docker-compose.yaml
sed -i "s~connect-template.yaml~$baseName.yaml~g" cmd/root/root.go
sed -i "s~*connect-template~*$baseName~g" internal/test/app.go
sed -i "s~connect-template.example.yaml~$baseName.example.yaml~g" internal/test/app.go
mv connect-template.example.yaml "$baseName.example.yaml"
echo "$baseName.yaml" >> .gitignore
