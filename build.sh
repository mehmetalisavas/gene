#! /bin/bash
set -o errexit

geneservices=(
  github.com/cihangir/gene/cmd/gene
  github.com/cihangir/gene/plugins/gene-models
  github.com/cihangir/gene/plugins/gene-rows
  github.com/cihangir/gene/plugins/gene-errors
  github.com/cihangir/gene/plugins/gene-kit
  github.com/cihangir/gene/plugins/gene-dockerfiles
  github.com/cihangir/geneddl/cmd/gene-ddl
)

`which go` install -v "${geneservices[@]}"
