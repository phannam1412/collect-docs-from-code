# collect-docs-from-code
Collect documentations from code and create a readme file

Motivation
==========

I usually write a lots of explainatory comments in source code, but after finishing a project, 

I should create a readme file for components, packages, modules that I have implemented. 

At this time, I need to choose:

- if I should copy all source code comments to the readme file and then delete these comments

- or if I should write a new readme and still keep all source code comments

- or not write a readme at all

I really want to keep all source code comments because it must easier for anyone to read and understand source code quickly.

So, my final solution is to write a tools to collect all source code comments and put everything into a readme file.

How it work
===========

The script will parse source code files utilizing a sophisticated pattern matching algorithm 
and create a dictionary for all items.

Then it will sort all items and write to the destination file. 

Reference
=========

- [Extended Backus–Naur form](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form)

- [Pattern matching implementation](https://github.com/phannam1412/go-pattern-matching)
