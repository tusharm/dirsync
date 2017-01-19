### dirsync

Quick and dirty way to sync a directory against another. Useful for incremental backups.

```
Usage: dirsync <src-dir> <dest-dir>
```

Things to fix:

+ circular dependencies due to symlinks
+ handle copying of symlinks
+ confgiurable parallel executions