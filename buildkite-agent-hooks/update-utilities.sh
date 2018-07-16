#!/bin/bash
echo "Updating 'ci-studio-common'"
ci-studio-common-util update

echo "Updating 'hab'"
sudo install-habitat
