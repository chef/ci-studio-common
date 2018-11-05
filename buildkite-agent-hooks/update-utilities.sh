#!/bin/bash
echo "Updating 'ci-studio-common'"
sudo ci-studio-common-util update

echo "Updating 'hab'"
sudo install-habitat

echo "Updating 'expeditor' CLI"
hab pkg install chef-es/expeditor-ruby
hab pkg binlink --force chef-es/expeditor-ruby expeditor
