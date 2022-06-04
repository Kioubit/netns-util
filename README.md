# netns-util

## Small setuid utility around network namespaces

Setup:

    go build -o netns-util .; sudo chown root:root netns-util; sudo chmod u+s netns-util

Using:

    Usage:
    ./netns-util enter [netns name] [command]
    ./netns-util enter-netadmin [netns name] [command]
    ./netns-util moveto [netns name] [link]
    ./netns-util netadmin [command]