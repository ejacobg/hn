# Hacker News Exporter

`hn` is a command-line tool used to export and organize a user's favorite or upvoted submissions and comments from Hacker News as JSON.

## Usage

```
Usage of hn:
        hn <export|sort|update> <favorite|upvoted> <submissions|comments> [args...] [flags...]
        hn <-h|-help>
Subcommands:
        hn export - returns a user's saved items as JSON
        hn sort - applies categories to exported items
        hn update - organizes export files, pulls unexported items
```

### Export

```
Usage of export:
export <favorite|upvoted> <submissions|comments> <username> [flags]
export <-h|-help>
  -page int
        (Optional) Which page to read from. (default 1)
  -password string
        Password for the given user.
  -token string
        Value of the 'user' cookie from a logged-in session. Takes priority over password.
To view upvoted posts, a password or token is required.
```

The `hn export` tool is used to transform a page of a user's saved submissions or comments, and returns them as JSON.

`hn export` will print its results to the standard output. If you wish to save them to a file, just redirect the output.

### Sort

```
Usage of sort:
sort <favorite|upvoted> <submissions|comments> [flags]
sort <-h|-help>
  -directory string
        (Optional) Custom directory to read from. Must be structured correctly.
```

The `hn sort` tool allows you to apply categories and reading statuses to your exported items. The `hn sort` tool will go through your exported items and prompt you to apply a new category for all items that don't already have one:

```
$ hn sort favorite submissions
Title: React Common Tools and Practices: State Management Overview
URL: https://react-community-tools-practices-cheatsheet.netlify.app/state-management/overview/
Discussion: https://news.ycombinator.com/item?id=26385984
1 - Advice
2 - Career
3 - Hobby
4 - Money
5 - Opinion
6 - Product
7 - Project
8 - Reference
9 - Skill
0 - Tip
f - Finish
```

There are only 10 preset categories. An item can only have 1 category.

When an item is given a category, an entry will be appended to the appropriate category file (`<category>.md`). All entries have the same form:

```
[<status>] [<details>](<id>)
```

You are free to add any other content to these category files. Only lines that are of the above format will be processed by the program.

The `status` field can be anything. For my use case, it is just a 1-word description describing an action taken on that item. For example, I might use `notes` for items I have written notes on, `seen` for items I've read but haven't taken notes on, and `skip` for items I feel aren't worth reading. If the `status` field is blank, then the item is considered unread.

The purpose of these category files is to provide a more convenient interface into your saved items. You can use these files to keep track of which items you have read, and organize them into some logical format. A typical use case is to add headings for subcategories within each category file. For example, a `money.md` category file might look something like this:

```markdown
# Accounting/Budgeting
[] [Accounting For Developers, Part I](32495724)
[] [Accounting for Developers, Part II](32580016)
[] [Ask HN: How do you record your personal finances?](31605741)

# Investing
[] [Ask HN: How to invest savings?](31563463)
[] [Ask HN: How to save/invest/deal with money?](33668398)
[] [Ask HN: How do you invest your money?](31718948)

# Uncategorized
...
```

After the `hn sort` tool is run, the program will try to reconcile the state of the category files with the state of the export files (more on export files later). It will apply any changes made to the category files onto its respective item. For example, if I've made notes on one of the articles from above, I would update it in the `money.md` file:

```markdown
# Accounting/Budgeting
[notes] [Accounting For Developers, Part I](32495724)
...
```

After running the `hn sort` tool, the respective item would be updated in the export files:

```json
{
  "id": "32495724",
  "category": "money",
  "state": "notes",
  "title": "Accounting For Developers, Part I",
  "url": "https://www.moderntreasury.com/journal/accounting-for-developers-part-i",
  "discussion": "https://news.ycombinator.com/item?id=32495724"
}
```

Your main interaction with the saved items is through the category files. The purpose of the export files is mainly to obtain the item's relevant URLs since they aren't stored in the category files.

### Update

```
Usage of update:
update <favorite|upvoted> <submissions|comments> <username> [flags]
update <favorite|upvoted> <submissions|comments> -shuffle
update <-h|-help>
  -directory string
        (Optional) Directory to be updated. (default "./<favorite|upvoted>/<submissions|comments>")
  -password string
        Password for the given user.
  -shuffle
        Shuffle items in the given directory.
  -token string
        Value of the 'user' cookie from a logged-in session. Takes priority over password.
To view upvoted posts, a password or token is required.
```

The `hn update` command acts as a smarter version of `hn export`. It will continuously pull the saved items in your profile until it has caught up to the most recently saved ones. If you have partially exported some items already, the `hn update` command will ensure that it will not save any items already present in an export file.

The `hn update` command saves its output to a separate directory, `updated/`. You should confirm that the output and filenames are good before moving these files to the `exported/` directory.

HN displays your saved items 30 at a time. To enforce an upper limit of 30 items across all your export files, the `-shuffle` flag is provided. This will reorder all the items in your export files such that they will all contain 30 items except for the "last" one. The "last" file should contain your most recently saved items.

## Directory Structure

Here is the default directory structure assumed by the program:

```
ejacobg/
├── favorite/
│   ├── comments/
│   │   ├── exported/
│   │   │   ├── updated/
│   │   │   │   ├── 21.json
│   │   │   │   ├── 22.json
│   │   │   │   └── ...
│   │   │   ├── 01.json
│   │   │   ├── 02.json
│   │   │   ├── ...
│   │   │   └── 20.json
│   │   ├── advice.md
│   │   ├── career.md
│   │   └── ...
│   └── submissions/
│       └── ...
└── upvoted/
    └── ...
```

By default, the `hn` tool assumes that it is being run in the root directory. The root directory (in this case `ejacobg/`) can be anything. It is assumed that only 1 user's data is saved in this directory. Underneath this directory are the `favorite/` and `upvoted/` directories, which hold their respective items. The `favorite/` and `upvoted/` directories have the same structure.

Underneath the `favorite/` and `updated/` directories are the `comments/` and `submissions/` directories. These contain the category files (`<category>.md`) which serve as an index into your exported items.

The export files are found under the `exported/` directory. Under this scheme, the filenames are an incrementing sequence of integers, with 0's added to the start so that they are ordered correctly. The first file in the sequence (`01.json`) contains the oldest saved items (i.e. the first 30 items that you saved), while the last file in the sequence (`20.json`) contains the most recently saved items.

Files generated by the `hn update` tool are placed in the `updated/` directory. They are named so that they continue the sequence described above. You may need to add leading 0's to the filenames since the `hn update` tool will not do this for you.

## Walkthrough

Inside your root directory, create the directories used by the program:

```shell
mkdir -p favorite/comments/exported/updated
mkdir -p favorite/submissions/exported/updated
mkdir -p upvoted/comments/exported/updated
mkdir -p upvoted/submissions/exported/updated
```

While still inside the root directory, pull all of your account's items (in this case favorite submissions):

```shell
hn update favorite submissions ejacobg
```

The above command will populate the `favorite/submissions/exported/updated/` directory with a `1.json` file:

```json
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

Verify the contents of your exported files, then move them to the `favorite/submissions/exported/` directory, renaming them as needed. Run the `hn update` command with the `-shuffle` flag to redistribute the items across all the export files.

```shell
hn update favorite submissions -shuffle
```

In this case, a shuffle will not change anything, but it is good practice to perform a shuffle after an update.

Once everything has been moved into the `favorite/submissions/exported/` directory, start the sorting process using `hn sort`:

```shell
hn sort favorite submissions
```

This will display the sorting screen for all items.

```
Title: Hacker News Official API
URL: https://github.com/HackerNews/API
Discussion: https://news.ycombinator.com/item?id=32540883
1 - Advice
2 - Career
3 - Hobby
4 - Money
5 - Opinion
6 - Product
7 - Project
8 - Reference
9 - Skill
0 - Tip
f - Finish

```

The only accepted input are the options shown above (`0` through `9` and `f`). Any other option is considered a "skip", and no category will be applied to that item. After inputting a category, an entry will be made in the appropriate category file. Exiting the sorting process -- either through inputting `f` or simply having processed all items -- will then save these categories back into the export files.

From here, you may repeat this process for all the other directories. As you read through all of your saved items, update their category file entries, then use `hn sort` to apply those changes back to the export files.

If you've saved any more items since the last update, simply run `hn update` then `hn update -shuffle` to merge those new items into your database. After merging, run `hn sort` to categorize all these new items.
