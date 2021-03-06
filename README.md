### GitBlog 

A simple web server written in Go that takes the contents of a Git repository, namely Markdown and HTML documents, then serves them as a static website. 

It is a similar concept as [Hugo](https://gohugo.io), [Jekyll](https://jekyllrb.com), and [Pelican](https://jekyllrb.com), but strives to be even simplier in its approach.

It will implement similar functionality as [GitHub pages](https://pages.github.com), where publishing new content is as simple as pushing to a master (or specified) branch of the repository containing the static content.


### Features

#### Dynamic Navigation 

Will automatically create a nav bar that can be inserted into a template with the tag `{{ nav }}` based on the directories in the Git repo. Note: Hidden directories prefixed with a "." will be skipped (for example `.git`). 