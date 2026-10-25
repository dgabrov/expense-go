@internal/data/data.go contains data definitions. 
There is a type ConfigData

The app expects an env variable CONFIG_FILE that has the file name of a file that currently is @docs/info/config.json

Please load that file in a ConfigData variable at the beginning of the application. This is the externalized 
config data. Please make it in a function in a different file alongside @internal/proj/start.go because
I do not want to make start.go too bloated
