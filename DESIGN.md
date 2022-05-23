* if found `dir/key/end`, (confirm work with non-exist dir), count `cache_hit`, serve cache, append `ci_source` to `cache_served`
* if found `dir/key/start`, (confirm work with non-exit dir) skip cache, count `cache_miss`
* create `dir/key`, skip cache if failed, count `cache_miss`
* write `start` with timestamp, skip cache if failed or existant, count `cache_miss`
  > https://stackoverflow.com/questions/33223564/atomically-creating-a-file-if-it-doesnt-exist-in-python
* write cache and stdout
* if exit=0, append `end` with timestamp
* if exit!=0, force remove `dir/key`, if can't, force remove `dir/key/end`
* gauge `stdout_size`

cleanup
* iterate `dir/*`, read `end`, if older then x days
* force remove `end`
* force remove `dir/key`
