#!/bin/bash
kustomize build --load-restrictor LoadRestrictionsNone | kubectl apply -f -

