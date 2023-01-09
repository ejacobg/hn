# Hacker News Exporter

`hn` is a command-line tool used to export a user's favorite or upvoted submissions and comments from Hacker News as
JSON.

## Usage

```
Usage of hn:
hn <favorite|upvoted> <submissions|comments> <username> [flags]
hn <-h|-help>
  -page int
        (Optional) Which page to read from. (default 1)
  -password string
        Password for the given user.
  -token string
        Value of the 'user' cookie from a logged-in session. Takes priority over password.
To view upvoted posts, a password or token is required.
```

## Examples

```
hn favorite submissions ejacobg
[
    {
        "id": "32540883",
        "title": "Hacker News Official API",
        "url": "https://github.com/HackerNews/API",
        "discussion": "https://news.ycombinator.com/item?id=32540883"
    },
    {
        "id": "22788236",
        "title": "Show HN: Export HN Favorites to a CSV File",
        "url": "item?id=22788236",
        "discussion": "https://news.ycombinator.com/item?id=22788236"
    }
]
```

```
hn favorite comments ejacobg
[
    {
        "id": "32543023",
        "parent": "Hacker News Official API",
        "text": [
            "In response to people's complaints about the usability of the HN Firebase API: yes. We're going to eventually have a new API that returns a simple JSON version of any HN URL. At that point we'll phase out the Firebase API, with a generous deprecation period. I'd be curious to hear people's thoughts about what a generous deprecation period might be."
        ],
        "context": "https://news.ycombinator.com/context?id=32543023",
        "discussion": "https://news.ycombinator.com/item?id=32540883"
    }
]
```

```
hn upvoted submissions ejacobg -password=<password>
[
    {
        "id": "8863",
        "title": "My YC app: Dropbox - Throw away your USB drive",
        "url": "http://www.getdropbox.com/u/2/screencast.html",
        "discussion": "https://news.ycombinator.com/item?id=8863"
    },
    {
        "id": "121003",
        "title": "Ask HN: The Arc Effect",
        "url": "item?id=121003",
        "discussion": "https://news.ycombinator.com/item?id=121003"
    }
]
```
