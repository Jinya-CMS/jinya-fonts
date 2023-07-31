---
permalink: index.html
layout: "default"
---

# Jinya Fonts

Jinya Fonts is a simple dropin replacement for Google Fonts that doesn't track your users. The public instance of Jinya
Fonts is available under [fonts.jinya.de](https://fonts.jinya.de).

## How to run

Jinya Fonts is a Go application that has two commands. The first command serves the font on port 8090. The second
command syncs the fonts from Google Fonts to your hard drive.

### Configuration

The configuration options of Jinya Fonts are rather simple, this is the structure:

```yaml
api_key: <api token>
font_file_folder: ./data
filter_by_name:
  - <font name>
admin_password: <secure admin password>
serve_website: true
```

The `api token` is has to be a valid Google Fonts API token, you can generate
one [here](https://console.developers.google.com/apis/credentials).   
The `filter_by_name` option is optional and allows you to only sync the fonts you want.  
The `admin_password` option, should be a secure random token. This password can be used to login into the admin
dashboard. If you don't set it, the admin dashboard will be disabled.  
If the `serve_website` option is set to true the frontend of Jinya Fonts will be served. 

### Serve

To serve the Go application on port 8090, simply run the compiled Go binary with the command serve and the config file
provided. This is how the command looks:

    ./jinya-fonts serve -config-file=./config.yaml

### Sync

To sync the fonts from Google Fonts run the application with the command sync and the config file provided, this looks
like follows:

    ./jinya-fonts sync -config-file=./config.yaml

### Run with Docker

The most convenient way to run Jinya Fonts is with our prebuilt docker images. Simple run the following command:

For serving the Jinya Fonts endpoints:

    docker run -d -p 8090:8090 --name jinya-fonts -v /your/jinya/fonts/dir:/data jinyacms/jinya-fonts /app/jinya-fonts -config-file=/data/config.yaml serve

For syncing from Google Fonts:

    docker run -d --name jinya-fonts-sync -v /your/jinya/fonts/dir:/data jinyacms/jinya-fonts /app/jinya-fonts -config-file=/data/config.yaml sync

## Using Jinya Fonts

Jinya Fonts is API compatible to Google Fonts. The advantage is, that you can simply choose your font set
from [fonts.google.com](https://fonts.google.com). Then you replace host from `fonts.googleapis.com` with your host. And
that is all. For example

    https://fonts.googleapis.com/css2?family=Space+Mono:ital,wght@0,400;0,700;1,400;1,700&display=swap

Turns into:

    https://fonts.jinya.de/css2?family=Space+Mono:ital,wght@0,400;0,700;1,400;1,700&display=swap

Jinya Fonts also has a web interface where you can grab filter for fonts and select the weights and styles you want.
Just check fonts.jinya.de. If you host your own Jinya Fonts instance, point the browser to your own instance.

## Custom fonts

Adding custom fonts is rather simple, just access the admin dashboard under https://\<jinya-fonts-instance>/admin and
enter the admin password. After that you can create a new font and add font files to it.

## Why should I use Jinya Fonts?

Jinya Fonts doesn't track your users, and therefore you don't need to mention it in your data protection page. We also
set no cookies or similar. Apart from that we also disabled

## Found a bug?

If you found a bug feel free to create an issue on Github or on my personal Taiga
instance: https://taiga.imanuel.dev/project/jinya-fonts/

## License

Jinya Fonts is licensed under the MIT License.
