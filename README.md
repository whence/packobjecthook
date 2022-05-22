A POC on https://git-scm.com/docs/git-config#Documentation/git-config.txt-uploadpackpackObjectsHook

Why packobjects can be cached
https://gitlab.com/gitlab-org/gitlab-git/-/commit/20b20a22f8f7c1420e259c97ef790cb93091f475

> You may want to insert a caching layer around
     pack-objects; it is the most CPU- and memory-intensive
     part of serving a fetch, and its output is a pure
     function[1] of its input, making it an ideal place to
     consolidate identical requests.

> [1] Pack-objects isn't _actually_ a pure function. Its
    output depends on the exact packing of the object
    database, and if multi-threading is used for delta
    compression, can even differ racily. But for the
    purposes of caching, that's OK; of the many possible
    outputs for a given input, it is sufficient only that we
    output one of them.
