# github-profile-terminal-action

This is the code that generates my [profile README.md](https://github.com/liamg), including the terminal gif.

You can use it to automatically create your GitHub profile README by adding a GitHub action to your profile repository.

Here's an example config that updates every 4 hours:

```yaml
name: auto-update

on:
  workflow_dispatch:
  schedule:
    - cron:  0 */4 * * *

jobs:
  build:
    name: auto update
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@master
    - uses: liamg/github-profile-terminal-action@main
      with:
        feed_url: https://www.liam-galvin.co.uk/feed.xml
        twitter_username: liam_galvin
        theme: dark
        token: ${{ secrets. GITHUB_TOKEN }}
```
