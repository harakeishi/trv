<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

<div align="center">    
    <img src="logo.png" alt="logo">
</div>

`trv` is a remote viewer for [tbls](https://github.com/k1LoW/tbls).
This command is used to view `tbls` information stored in github, e.g. in a terminal.

# Table of Contents
- [Table of Contents](#table-of-contents)
- [DEMO](#demo)
- [Installation](#installation)
- [Usage](#usage)
    - [Setup](#setup)
        - [Example](#example)
- [Function](#function)
- [License](#license)

# DEMO
![image](https://user-images.githubusercontent.com/44335168/199011337-796428e4-88d7-4d40-983d-eeee6189fc45.gif)

# Installation
 
```bash
$ go install github.com/harakeishi/trv@latest
```

# Usage
## Setup
![setup](https://user-images.githubusercontent.com/44335168/199011411-8b9f3da0-df23-45d2-b9ca-a507ae007866.gif)

The first time it is started, a configuration file is created.
Fill in the information for documents generated by `tbls`.

 | | Required |Example|Description|
 |--|--|--|--|
 |Owner| * |`harakeishi`|Owner of the repository where the information is stored|
 |Repo| * |`trv`|Repository where information is stored|
 |Path| * |`sampledb`|Path of the document generated by `tbls`|
 |Token| * | |Github Token|
 |IsEnterprise|  || Is it GHES?|
 |BaseURL|| `https://git.hoge.com/api/v3/` | Base URL for GHES|
 |UploadURL|| `https://git.hoge.com/api/uploads/` | Upload URL for GHES|

Tokens for retrieving information in public repositories can be used without setting any specific permissions.

### Example
target: https://github.com/harakeishi/trv/tree/v0.0.11/sampledb

![image](https://user-images.githubusercontent.com/44335168/198053241-12f34946-1af5-4364-b53b-916eefc3e6a3.png)

# Function
- Easily switch between multiple data sources
- Filter by field name or table name
- View column details while displaying table information

# License
[MIT](LICENSE)
## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center"><a href="https://github.com/tommy6073"><img src="https://avatars.githubusercontent.com/u/3647470?v=4?s=100" width="100px;" alt="Takayuki Nagatomi"/><br /><sub><b>Takayuki Nagatomi</b></sub></a><br /><a href="https://github.com/harakeishi/trv/commits?author=tommy6073" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/dojineko"><img src="https://avatars.githubusercontent.com/u/1488898?v=4?s=100" width="100px;" alt="dojineko"/><br /><sub><b>dojineko</b></sub></a><br /><a href="https://github.com/harakeishi/trv/commits?author=dojineko" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
