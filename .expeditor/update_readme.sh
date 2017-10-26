#!/bin/bash

markdown-toc -i README.md

# This is likely a helper I'll want to pull into Expeditor so that we can use it globally.
perl -0777 -pi -e 's/<!-- stdout "([\w\s\-\.\/]+)" -->.*?<!-- stdout -->/"<!-- stdout \"".$1."\" -->\n```\n".`$1`."```\n<!-- stdout -->"/egs' README.md