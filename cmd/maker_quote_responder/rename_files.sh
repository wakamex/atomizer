#!/bin/bash
# Rename the original market_maker.go to preserve it
mv market_maker.go market_maker_original.go
# Rename the refactored file to be the new market_maker.go
mv market_maker_refactored.go market_maker.go