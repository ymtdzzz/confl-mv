# confl-mv

`confl-mv` is a command line tool to move Confluence pages quickly.

It is motivated by the fact that Confluence has a limit of up to 100 page movement APIs and Web UI.

## How to use

Install

```sh
go install github.com/ymtdzzz/confl-mv
```

Set environment variables

```sh
export CONFLUENCE_DOMAIN="your-domain.atlassian.net"
export CONFLUENCE_USERNAME="your_name@example.com"
export CONFLUENCE_APIKEY="your api key..."
```

First, attempt to navigate to the target page.

```sh
confl-mv movethis <target_page_id> <destination_page_id> -d ${CONFLUENCE_DOMAIN} -u ${CONFLUENCE_USERNAME} -a ${CONFLUENCE_APIKEY}
```

If this fails, a temporary page must be created and evacuate the child pages because the page contains more than 100 pages

```sh
# fail!
Failed to move page [<target_page_id>] to [<destination_page_id>]
because it has over 99 child pages.
Move it's child pages first

# create temporary page on Confluence API or WebUI (currently this tool doesn't has this feature 🙏)

# move it's child page to temporary page
confl-mv movechild <target_page_id> <temporary_page_id> -d ${CONFLUENCE_DOMAIN} -u ${CONFLUENCE_USERNAME} -a ${CONFLUENCE_APIKEY}
```

Then move the target page.

```sh
# try again!
confl-mv movethis <target_page_id> <destination_page_id> -d ${CONFLUENCE_DOMAIN} -u ${CONFLUENCE_USERNAME} -a ${CONFLUENCE_APIKEY}
```

Finally, move the evacuated pages to the target page.

```sh
confl-mv movechild <temporary_page_id> <target_page_id> -d ${CONFLUENCE_DOMAIN} -u ${CONFLUENCE_USERNAME} -a ${CONFLUENCE_APIKEY}
```
