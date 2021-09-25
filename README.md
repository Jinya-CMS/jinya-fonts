# Jinya Fonts
Jinya Fonts is a simple dropin replacement for Google Fonts that doesn't track your users. The public instance of Jinya Fonts is available under [fonts.jinya.de](https://fonts.jinya.de). 

## How to run
Jinya Fonts is a Go application that has two commands. The first command serves the font on port 8090. The second command syncs the fonts from Google Fonts to your hard drive.

### Configuration
The configuration options of Jinya Fonts are rather simple, this is the structure:

    api_key: <api token>
    font_file_folder: ./data
    filter_by_name:
        - <font name>

The `filter_by_name` option is optional and allows you to only sync the fonts you want.

### Serve
To serve the Go application on port 8090, simply run the compiled Go binary with the command serve and the config file provided. This is how the command looks:

    ./jinya-fonts serve -config-file=./config.yaml

### Sync
To sync the fonts from Google Fonts run the application with the command sync and the config file provided, this looks like follows:

    ./jinya-fonts sync -config-file=./config.yaml

## Using Jinya Fonts
Jinya Fonts is API compatible to Google Fonts. The advantage is, that you can simply choose your font set from [fonts.google.com](https://fonts.google.com) and the host from Google with your host. That is all. For example 

    https://fonts.googleapis.com/css2?family=Space+Mono:ital,wght@0,400;0,700;1,400;1,700&display=swap

Turns into:

    https://fonts.jinya.de/css2?family=Space+Mono:ital,wght@0,400;0,700;1,400;1,700&display=swap

## Custom fonts
Currently, custom fonts are supported, but needs a deeper understanding of YAML.

To create a custom font, you need to grab the WOFF2 files from your provider and place it in a folder with its name in your data directory.

The file name must follow the following pattern:

    <Fontname>.<Character subset>.<weight as number>.woff2

You also need to provide a configuration file in the data directory. Here is an example for the font Lato. The config options category and unicode_range are optional.

    name: Lato
    fonts:
      - path: Lato.latin-ext.100.woff2
        subset: latin-ext
        variant: "100"
        unicode_range: U+0100-024F, U+0259, U+1E00-1EFF, U+2020, U+20A0-20AB, U+20AD-20CF, U+2113, U+2C60-2C7F, U+A720-A7FF
        weight: "100"
        style: normal
        category: sans-serif
      - path: Lato.latin.100.woff2
        subset: latin
        variant: "100"
        unicode_range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD
        weight: "100"
        style: normal
        category: sans-serif

## Run with Docker
The most convenient way to run Jinya Fonts is with our prebuilt docker images. Simple run the following command:

For serving the Jinya Fonts endpoints:

    docker run -d -p 8090:8090 --name jinya-fonts -v /your/jinya/fonts/dir:/data jinyacms/jinya-fonts /app/jinya-fonts -config-file=/data/config.yaml serve

For syncing from Google Fonts:

    docker run -d --name jinya-fonts-sync -v /your/jinya/fonts/dir:/data jinyacms/jinya-fonts /app/jinya-fonts -config-file=/data/config.yaml sync

## Why should I use Jinya Fonts?
Jinya Fonts doesn't track your users and therefore you don't need to mention it in your data protection page. We also set no cookies or similar.

## Planned features
On our roadmap are currently two more features. The first is, that we want to provide you a nice web UI to look through the fonts and generate a link to your custom font collection.

The other feature planned is, that you can support custom fonts and upload them. Currently, custom fonts are supported. You only need to provide a few files, as discussed above.

## Found a bug?
If you found a bug feel free to create an issue on Github or on my personal Taiga instance: https://taiga.imanuel.dev/project/jinya-fonts/

## License
Jinya Fonts is licensed under the MIT License.