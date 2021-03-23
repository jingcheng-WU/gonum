#!/bin/bash

# Reset the tree to the current commit to handle
# any writes during the build.
git reset --hard

go generate github.com/jingcheng-WU/gonum/blas
go generate github.com/jingcheng-WU/gonum/blas/gonum
go generate github.com/jingcheng-WU/gonum/unit
go generate github.com/jingcheng-WU/gonum/unit/constant
go generate github.com/jingcheng-WU/gonum/graph/formats/dot
go generate github.com/jingcheng-WU/gonum/graph/formats/rdf
go generate github.com/jingcheng-WU/gonum/stat/card

if [ -n "$(git diff)" ]; then	
	git diff
	exit 1
fi

rm -rf *
git reset --hard
